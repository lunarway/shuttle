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

const rootCmdCompletion = `
__shuttle_run_script_args() {
	local cur prev args_output args
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"

	template=$'{{ range $i, $arg := .Args }}{{ $arg.Name }}\n{{ end }}'
	if args_output=$(shuttle --skip-pull run "$1" --help --template "$template" 2>/dev/null); then
		args=($(echo "${args_output}"))
		COMPREPLY=( $( compgen -W "${args[*]}" -- "$cur" ) )
		compopt -o nospace
	fi
}

# find available scripts to run
__shuttle_run_scripts() {
	local cur prev scripts currentScript
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	currentScript="${COMP_WORDS[2]}"

	if [[ ! "${prev}" = "run" ]]; then
		__shuttle_run_script_args $currentScript
		return 0
	fi

	if scripts_output=$(shuttle --skip-pull ls 2>/dev/null); then
		scripts=($(echo "${scripts_output}" | tail +3 | awk '{print $1}'))
		COMPREPLY=( $( compgen -W "${scripts[*]}" -- "$cur" ) )
	fi
	return 0
}

# called when the build in completion fails to match
__shuttle_custom_func() {
  case ${last_command} in
      shuttle_run)
          __shuttle_run_scripts
          return
          ;;
      *)
          ;;
  esac
}
`

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
	BashCompletionFunction: rootCmdCompletion,
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
