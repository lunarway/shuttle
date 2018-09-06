package cmd

import (
	"bitbucket.org/LunarWay/shuttle/pkg/executors"
	"github.com/spf13/cobra"
)

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

var runCmd = &cobra.Command{
	Use:   "run [command]",
	Short: "Run a plan script",
	Long:  `Specify which plan script to run`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var commandName = args[0]
		context := getProjectContext()
		executors.Execute(context, commandName, args[1:])
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

func init() {
	rootCmd.AddCommand(runCmd)
}
