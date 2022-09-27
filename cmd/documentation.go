package cmd

import (
	"github.com/lunarway/shuttle/cmd/utility"
	"github.com/lunarway/shuttle/pkg/browser"
	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func newDocumentation(uii *ui.UI, contextProvider utility.ContextProvider) *cobra.Command {
	documentationCommand := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := contextProvider()
			if err != nil {
				return err
			}

			url, err := context.DocumentationURL()
			if err != nil {
				return err
			}

			uii.Infoln("Documentation available at: %s", url)

			browseCmd, err := browser.Command(url, cmd.ErrOrStderr())
			if err != nil {
				return err
			}

			err = browseCmd.Run()
			if err != nil {
				return errors.NewExitCode(1, "Failed to open document reference: %v", err)
			}

			return nil
		},
	}

	return documentationCommand
}
