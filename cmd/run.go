package cmd

import (
	stdcontext "context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors"
	"github.com/lunarway/shuttle/pkg/ui"
)

func newNoopRun() *cobra.Command {
	return &cobra.Command{
		Use:          "run [command]",
		Short:        "Run a plan script",
		Long:         `Specify which plan script to run`,
		SilenceUsage: true,
	}
}

func newNoContextRun() *cobra.Command {
	runCmd := newNoopRun()

	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("shuttle run is not available in this context. To use shuttle run you need to be in a project with a shuttle.yaml file")
	}

	return runCmd
}

func newRun(uii *ui.UI, contextProvider contextProvider) (*cobra.Command, error) {
	var (
		flagTemplate   string
		validateArgs   bool
		interactiveArg bool
	)
	shuttleInteractive := os.Getenv("SHUTTLE_INTERACTIVE")
	var shuttleInteractiveDefault bool
	if shuttleInteractive == "true" {
		shuttleInteractiveDefault = true
	}

	executorRegistry := executors.NewRegistry(executors.ShellExecutor, executors.TaskExecutor)

	runCmd := newNoopRun()

	context, err := contextProvider()
	if err != nil {
		return nil, err
	}

	// For each script construct a run command specific for said script
	for script, value := range context.Scripts {
		runCmd.AddCommand(
			newRunSubCommand(
				uii,
				context,
				script,
				value,
				executorRegistry,
				&interactiveArg,
				&validateArgs,
			),
		)
	}

	runCmd.PersistentFlags().
		StringVar(&flagTemplate, "template", "", "Template string to use. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].")
	runCmd.PersistentFlags().
		BoolVar(&validateArgs, "validate", true, "Validate arguments against script definition in plan and exit with 1 on unknown or missing arguments")
	runCmd.PersistentFlags().
		BoolVar(&interactiveArg, "interactive", shuttleInteractiveDefault, "sets whether to enable ui for getting missing values via. prompt instead of failing immediadly, default is set by [SHUTTLE_INTERACTIVE=true/false]")
	return runCmd, nil
}

func newRunSubCommand(
	uii *ui.UI,
	context config.ShuttleProjectContext,
	script string,
	value config.ShuttlePlanScript,
	executorRegistry *executors.Registry,
	interactiveArg *bool,
	validateArgs *bool,
) *cobra.Command {
	// Args are best suited as kebab-case on the command line
	argName := func(input string) string {
		return strcase.ToKebab(input)
	}

	parseKeyValuePair := func(arg string) (string, string, bool) {
		before, after, found := strings.Cut(arg, "=")
		if found {
			return before, after, true
		}

		return "", "", false
	}

	// Legacy key=value pairs into standard args that cobra can understand
	applyLegacyArgs := func(args []string, inputArgs map[string]*string) {
		for _, inputArg := range args {
			key, value, ok := parseKeyValuePair(inputArg)
			if ok {
				inputArgs[key] = &value
			}
		}
	}

	// In case interactive is turned on and arg is missing, we ask for missing values
	createPrompt := func(inputArgs map[string]*string, arg config.ShuttleScriptArgs) (string, error) {
		prompt := []*survey.Question{
			{
				Name: argName(arg.Name),
				Prompt: &survey.Input{
					Message: argName(arg.Name),
					Default: *inputArgs[arg.Name],
					Help:    arg.Description,
				},
			},
		}
		if arg.Required {
			prompt[0].Validate = survey.Required
		}
		var output string
		err := survey.Ask(prompt, &output)
		if err != nil {
			return "", err
		}

		return output, nil
	}

	// Decide whether to fall back on prompt or give a hard error
	validateInputArgs := func(value config.ShuttlePlanScript, inputArgs map[string]*string) error {
		for _, arg := range value.Args {
			if !arg.Required {
				continue
			}

			arg := arg

			if *inputArgs[arg.Name] == "" && *interactiveArg {
				output, err := createPrompt(inputArgs, arg)
				if err != nil {
					return err
				}
				if output != "" {
					inputArgs[arg.Name] = &output
				}

			} else if *inputArgs[arg.Name] == "" && arg.Required && *validateArgs {
				return fmt.Errorf("required flag(s) \"%s\" not set", argName(arg.Name))
			}
		}

		return nil
	}

	// Produce a stable list of arguments
	sort.Slice(value.Args, func(i, j int) bool {
		return value.Args[i].Name < value.Args[j].Name
	})

	// Initialize a collection of variables for use in args set pr command
	inputArgs := make(map[string]*string, 0)
	for _, arg := range value.Args {
		arg := arg
		inputArgs[arg.Name] = new(string)
	}

	cmd := &cobra.Command{
		Use:          script,
		Short:        value.Description,
		Long:         value.Description,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if *interactiveArg {
				uii.Verboseln("Running using interactive mode!")
			}

			ctx := cmd.Context()
			ctx, _, traceError, traceEnd := trace(ctx, script, args)
			defer traceEnd()

			applyLegacyArgs(args, inputArgs)
			if err := validateInputArgs(value, inputArgs); err != nil {
				return err
			}

			ctx, cancel := withSignal(ctx, uii)
			defer cancel()
			actualArgs := make(map[string]string, len(inputArgs))
			for k, v := range inputArgs {
				actualArgs[k] = *v
			}

			err := executorRegistry.Execute(ctx, context, script, actualArgs, *validateArgs)
			if err != nil {
				traceError(err)
				return err
			}

			return nil
		},
	}

	if !*validateArgs {
		cmd.Args = cobra.ArbitraryArgs
	}

	for _, arg := range value.Args {
		arg := arg
		cmd.Flags().StringVar(inputArgs[arg.Name], argName(arg.Name), "", arg.Description)
	}

	return cmd
}

// withSignal returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called,
// if the parent context's Done channel is closed, if a SIGINT signal is
// catched, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func withSignal(parent stdcontext.Context, uii *ui.UI) (stdcontext.Context, func()) {
	parent, cancel := stdcontext.WithCancel(parent)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case s := <-c:
			uii.Infoln("Received %v signal...", s)
			cancel()
		case <-parent.Done():
		}
	}()

	return parent, func() {
		signal.Stop(c)
		cancel()
	}
}
