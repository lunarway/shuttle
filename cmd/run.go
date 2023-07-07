package cmd

import (
	stdcontext "context"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lunarway/shuttle/pkg/executors"
	"github.com/lunarway/shuttle/pkg/telemetry"
	"github.com/lunarway/shuttle/pkg/ui"
)

func newRun(uii *ui.UI, contextProvider contextProvider) *cobra.Command {
	var (
		flagTemplate string
		validateArgs bool
	)

	executorRegistry := executors.NewRegistry(executors.ShellExecutor, executors.TaskExecutor)

	runCmd := &cobra.Command{
		Use:          "run [command]",
		Short:        "Run a plan script",
		Long:         `Specify which plan script to run`,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			commandName := args[0]
			ctx := stdcontext.Background()
			ctx = telemetry.WithContextID(ctx)
			ctx = WithRunTelemetry(ctx, commandName, args)
			telemetry.Client.Trace(
				ctx,
				"run",
				telemetry.WithPhase("start"),
			)
			defer func(ctx stdcontext.Context) {
				telemetry.Client.Trace(
					ctx,
					"run",
					telemetry.WithPhase("end"),
				)
			}(ctx)

			context, err := contextProvider()
			if err != nil {
				telemetry.Client.TraceError(
					ctx,
					"run",
					err,
				)
				return err
			}
			telemetry.Client.Trace(
				ctx,
				"run",
				telemetry.WithPhase("after-plan-pull"),
				telemetry.WithText("plan", context.Config.Plan),
			)

			ctx, cancel := withSignal(ctx, uii)
			defer cancel()

			err = executorRegistry.Execute(ctx, context, commandName, args[1:], validateArgs)
			if err != nil {
				telemetry.Client.TraceError(
					ctx,
					"run",
					err,
				)
				return err
			}

			return nil
		},
	}

	runCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		scripts := cmd.Flags().Args()
		if len(scripts) == 0 {
			runCmd.Usage()
			return
		}
		context, err := contextProvider()
		checkError(uii, err)

		err = executors.Help(context.Scripts, scripts[0], cmd.OutOrStdout(), flagTemplate)
		checkError(uii, err)
	})
	runCmd.Flags().
		StringVar(&flagTemplate, "template", "", "Template string to use. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].")
	runCmd.Flags().
		BoolVar(&validateArgs, "validate", true, "Validate arguments against script definition in plan and exit with 1 on unknown or missing arguments")
	return runCmd
}

func WithRunTelemetry(
	ctx stdcontext.Context,
	commandName string,
	args []string,
) stdcontext.Context {
	ctx = stdcontext.WithValue(ctx, telemetry.TelemetryCommand, commandName)
	if len(args) != 0 {
		// TODO: Make sure we sanitize secrets, somehow
		ctx = stdcontext.WithValue(ctx, telemetry.TelemetryCommandArgs, strings.Join(args[1:], " "))
	}
	return ctx
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
