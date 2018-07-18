package main

import (
	"log"
	"os"
	"path"

	"github.com/emilingerslev/shuttle/pkg/config"
	"github.com/emilingerslev/shuttle/pkg/plan"
	"github.com/spf13/cobra"
)

func main() {

	var ProjectPath string

	// 	var cmdBuild = &cobra.Command{
	// 		Use:   "build [path of project]",
	// 		Short: "Build the docker image of a project",
	// 		Long: `Build the docker image of a project.
	// This will use the internal Dockerfile of the shuttle.`,
	// 		Args: cobra.MinimumNArgs(0),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			plan := getPlan(ProjectPath)
	// 			//docker.Build(plan)
	// 		},
	// 	}

	var cmdRun = &cobra.Command{
		Use:   "run [command]",
		Short: "Run a plan command",
		Long:  `Specify which plan command to run`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p := getPlan(ProjectPath)
			var commandName = args[0]
			log.Println(plan.Execute(p, commandName))
		},
	}

	// 	var cmdEcho = &cobra.Command{
	// 		Use:   "echo [string to echo]",
	// 		Short: "Echo anything to the screen",
	// 		Long: `echo is for echoing anything back.
	// Echo works a lot like print, except it has a child command.`,
	// 		Args: cobra.MinimumNArgs(1),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			fmt.Println("Print: " + strings.Join(args, " "))
	// 		},
	// 	}

	// 	var cmdTimes = &cobra.Command{
	// 		Use:   "times [# times] [string to echo]",
	// 		Short: "Echo anything to the screen more times",
	// 		Long: `echo things multiple times back to the user by providing
	// a count and a string.`,
	// 		Args: cobra.MinimumNArgs(1),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			for i := 0; i < echoTimes; i++ {
	// 				fmt.Println("Echo: " + strings.Join(args, " "))
	// 			}
	// 		},
	// 	}

	// 	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{Use: "shuttle"}
	rootCmd.PersistentFlags().StringVarP(&ProjectPath, "project", "p", ".", "Project path")
	// rootCmd.AddCommand(cmdPrint, cmdEcho)
	// cmdEcho.AddCommand(cmdTimes)
	// rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdRun)
	rootCmd.Execute()
}

func getPlan(projectPath string) config.ShuttlePlan {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var fullProjectPath = path.Join(dir, projectPath)
	var c config.ShuttleConfig
	c.GetConf(fullProjectPath)
	var plan config.ShuttlePlan
	plan.Load(fullProjectPath, c)

	return plan
}
