package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/spf13/cobra"
)

var (
	projectPath string
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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&projectPath, "project", "p", ".", "Project path")
}

func getProjectContext() config.ShuttleProjectContext {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var fullProjectPath = path.Join(dir, projectPath)
	var c config.ShuttleProjectContext
	c.Setup(fullProjectPath)
	return c
}
