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
	uii                ui.UI = ui.Create()
	verboseFlag        bool
	clean              bool
	skipGitPlanPulling bool
	plan               string
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
		if verboseFlag {
			uii = uii.SetUserLevel(ui.LevelVerbose)
		}

		uii.Verboseln("Running shuttle")
		uii.Verboseln("- version: %s", version)
		uii.Verboseln("- commit: %s", commit)
		uii.Verboseln("- project-path: %s", projectPath)
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
	rootCmd.PersistentFlags().StringVar(&plan, "plan", "", `Overload the plan used.
Specifying a local path with either an absolute path (/some/plan) or a relative path (../some/plan) to another location
for the selected plan.
Select a version of a git plan by using #branch, #sha or #tag
If none of above is used, then the argument will expect a full plan spec.`)
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Print verbose output")
}

func getProjectContext() config.ShuttleProjectContext {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var fullProjectPath string
	if path.IsAbs(projectPath) {
		fullProjectPath = projectPath
	} else {
		fullProjectPath = path.Join(dir, projectPath)
	}

	if plan == "" {
		env := os.Getenv("SHUTTLE_PLAN_OVERLOAD")
		if env != "" {
			plan = env
		}
	}

	var c config.ShuttleProjectContext
	c.Setup(fullProjectPath, uii, clean, skipGitPlanPulling, plan)
	return c
}
