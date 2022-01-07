package cmd

import (
	"github.com/lunarway/shuttle/pkg/browser"
	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/spf13/cobra"
)

var documentationCommand = &cobra.Command{
	Use:     "documentation",
	Aliases: []string{"docs"},
	Short:   "Open documentation for the configured shuttle plan",
	Long: `Open documentation for the configured shuttle plan.
By default shuttle will try to open the plan's documentation in a web browser.

If no docs are explicitly configured in the plan, the plan it self is opened.
Usually this will target a hosted git repository, eg. GitHub README.

The application to open the documentation is inferred from the operating system
and respects the BROWSER environment variable.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		context, err := getProjectContext()
		checkError(err)

		url, err := context.DocumentationURL()
		checkError(err)
		uii.Infoln("Documentation available at: %s", url)

		browseCmd, err := browser.Command(url, cmd.ErrOrStderr())
		checkError(err)

		err = browseCmd.Run()
		if err != nil {
			checkError(errors.NewExitCode(1, "Failed to open document reference: %v", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(documentationCommand)
}
