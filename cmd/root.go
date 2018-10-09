package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/ui"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/spf13/cobra"
)

var (
	projectPath        string
	uii                ui.UI
	verboseFlag        bool
	clean              bool
	skipGitPlanPulling bool
	version            = "<dev-version>"
	commit             = "<unspecified-commit>"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shuttle",
	Short: "A CLI for handling shared build and deploy tools between many projects no matter what technologies the project is using.",
	Long: fmt.Sprintf(`shuttle %s

A CLI for handling shared build and deploy tools between many
projects no matter what technologies the project is using.

Read more about shuttle at https://github.com/lunarway/shuttle`, version),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		uii = ui.Create()
		if verboseFlag {
			uii = uii.SetUserLevel(ui.LevelVerbose)
		}

		uii.VerboseLn("Running shuttle")
		uii.VerboseLn("- version: %s", version)
		uii.VerboseLn("- commit: %s", commit)
		uii.VerboseLn("- project-path: %s", projectPath)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		var uii = ui.Create()
		uii.ExitWithErrorCode(1, "%s", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&projectPath, "project", "p", ".", "Project path")
	rootCmd.PersistentFlags().BoolVarP(&clean, "clean", "c", false, "Start from clean setup")
	rootCmd.PersistentFlags().BoolVar(&skipGitPlanPulling, "skip-pull", false, "Skip git plan pulling step")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Print verbose output")
}

func getProjectContext() config.ShuttleProjectContext {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var fullProjectPath = path.Join(dir, projectPath)
	var c config.ShuttleProjectContext
	c.Setup(fullProjectPath, uii, clean, skipGitPlanPulling)
	return c
}
