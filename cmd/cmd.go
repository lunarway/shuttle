package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/ui"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/spf13/cobra"
)

var (
	version = "<dev-version>"
	commit  = "<unspecified-commit>"
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

	template=$'{{ range $name, $script := .Scripts }}{{ $name }}\n{{ end }}'
	if scripts_output=$(shuttle --skip-pull ls --template "$template" 2>/dev/null); then
		scripts=($(echo "${scripts_output}"))
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

func newRoot(uii *ui.UI) (*cobra.Command, contextProvider) {
	var (
		verboseFlag        bool
		projectPath        string
		clean              bool
		skipGitPlanPulling bool
		plan               string
	)

	rootCmd := &cobra.Command{
		Use:   "shuttle",
		Short: "A CLI for handling shared build and deploy tools between many projects no matter what technologies the project is using.",
		Long: fmt.Sprintf(`shuttle %s

A CLI for handling shared build and deploy tools between many
projects no matter what technologies the project is using.

Read more about shuttle at https://github.com/lunarway/shuttle`, version),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verboseFlag {
				uii.SetUserLevel(ui.LevelVerbose)
			}
			uii.Verboseln("Running shuttle")
			uii.Verboseln("- version: %s", version)
			uii.Verboseln("- commit: %s", commit)
			uii.Verboseln("- project-path: %s", projectPath)
		},
		BashCompletionFunction: rootCmdCompletion,
	}

	rootCmd.PersistentFlags().StringVarP(&projectPath, "project", "p", ".", "Project path")
	rootCmd.PersistentFlags().BoolVarP(&clean, "clean", "c", false, "Start from clean setup")
	rootCmd.PersistentFlags().BoolVar(&skipGitPlanPulling, "skip-pull", false, "Skip git plan pulling step")
	rootCmd.PersistentFlags().StringVar(&plan, "plan", "", `Overload the plan used.
Specifying a local path with either an absolute path (/some/plan) or a relative path (../some/plan) to another location
for the selected plan.
Select a version of a git plan by using #branch, #sha or #tag
If none of above is used, then the argument will expect a full plan spec.`)
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Print verbose output")

	ctxProvider := func() (config.ShuttleProjectContext, error) {
		return getProjectContext(rootCmd, uii, projectPath, clean, plan, skipGitPlanPulling)
	}

	return rootCmd, ctxProvider
}

func Execute(out, err io.Writer) {
	rootCmd, uii := initializedRoot(out, err)

	if err := rootCmd.Execute(); err != nil {
		checkError(uii, err)
	}
}

func initializedRoot(out, err io.Writer) (*cobra.Command, *ui.UI) {
	uii := ui.Create(out, err)

	rootCmd, ctxProvider := newRoot(uii)
	rootCmd.SetOut(out)
	rootCmd.SetErr(err)

	rootCmd.AddCommand(newDocumentation(uii, ctxProvider))
	rootCmd.AddCommand(newCompletion(uii))
	rootCmd.AddCommand(newGet(uii, ctxProvider))
	rootCmd.AddCommand(newGitPlan(uii, ctxProvider))
	rootCmd.AddCommand(newHas(uii, ctxProvider))
	rootCmd.AddCommand(newLs(uii, ctxProvider))
	rootCmd.AddCommand(newPlan(uii, ctxProvider))
	rootCmd.AddCommand(newPrepare(uii, ctxProvider))
	rootCmd.AddCommand(newRun(uii, ctxProvider))
	rootCmd.AddCommand(newTemplate(uii, ctxProvider))
	rootCmd.AddCommand(newVersion(uii))

	return rootCmd, uii
}

type contextProvider func() (config.ShuttleProjectContext, error)

func getProjectContext(rootCmd *cobra.Command, uii *ui.UI, projectPath string, clean bool, plan string, skipGitPlanPulling bool) (config.ShuttleProjectContext, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	projectFlagSet := rootCmd.Flags().Changed("project")

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
	_, err = c.Setup(fullProjectPath, uii, clean, skipGitPlanPulling, plan, projectFlagSet)
	if err != nil {
		return config.ShuttleProjectContext{}, err
	}
	return c, nil
}
