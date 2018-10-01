package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/output"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/spf13/cobra"
)

var (
	projectPath string
	verbose     bool
	clean       bool
	version     = "<dev-version>"
	commit      = "<unspecified-commit>"
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
		if verbose {
			fmt.Println("Running shuttle")
			fmt.Println(fmt.Sprintf("- version: %s", version))
			fmt.Println(fmt.Sprintf("- commit: %s", commit))
			fmt.Println(fmt.Sprintf("- project-path: %s", projectPath))
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.ExitWithErrorCode(1, fmt.Sprintf("%s", err))
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&projectPath, "project", "p", ".", "Project path")
	rootCmd.PersistentFlags().BoolVarP(&clean, "clean", "c", false, "Start from clean setup")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output")
}

func getProjectContext() config.ShuttleProjectContext {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var fullProjectPath = path.Join(dir, projectPath)
	var c config.ShuttleProjectContext
	c.Setup(fullProjectPath, verbose, clean)
	return c
}
