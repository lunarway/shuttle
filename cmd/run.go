package cmd

import (
	stdcontext "context"
	"os"
	"os/signal"

	"github.com/lunarway/shuttle/pkg/executors"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [command]",
	Short: "Run a plan script",
	Long:  `Specify which plan script to run`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var commandName = args[0]
		context, err := getProjectContext()
		checkError(err)
		ctx, cancel := withSignal(stdcontext.Background())
		defer cancel()
		err = executors.Execute(ctx, context, commandName, args[1:], validateArgs)
		checkError(err)
	},
}

var (
	flagTemplate string
	validateArgs bool
)

func init() {
	runCmd.SetHelpFunc(func(f *cobra.Command, args []string) {
		scripts := f.Flags().Args()
		if len(scripts) == 0 {
			runCmd.Usage()
			return
		}
		context, err := getProjectContext()
		checkError(err)

		err = executors.Help(context.Scripts, scripts[0], os.Stdout, flagTemplate)
		checkError(err)
	})
	runCmd.Flags().StringVar(&flagTemplate, "template", "", "Template string to use. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].")
	runCmd.Flags().BoolVar(&validateArgs, "validate", true, "Validate arguments against script definition in plan and exit with 1 on unknown or missing arguments")
	rootCmd.AddCommand(runCmd)
}

// withSignal returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called,
// if the parent context's Done channel is closed, if a SIGINT signal is
// catched, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func withSignal(parent stdcontext.Context) (stdcontext.Context, func()) {
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
