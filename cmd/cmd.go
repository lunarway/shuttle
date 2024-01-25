package cmd

import (
	stdcontext "context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors/golang/executer"
	"github.com/lunarway/shuttle/pkg/telemetry"
	"github.com/lunarway/shuttle/pkg/ui"
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

func newRoot(uii *ui.UI) (*cobra.Command, contextProvider, repositoryContext) {
	telemetry.Setup()

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
	rootCmd.PersistentFlags().
		BoolVar(&skipGitPlanPulling, "skip-pull", false, "Skip git plan pulling step")
	rootCmd.PersistentFlags().StringVar(&plan, "plan", "", `Overload the plan used.
Specifying a local path with either an absolute path (/some/plan) or a relative path (../some/plan) to another location
for the selected plan.
Select a version of a git plan by using #branch, #sha or #tag
If none of above is used, then the argument will expect a full plan spec.`)
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Print verbose output")

	ctxProvider := func() (config.ShuttleProjectContext, error) {
		return getProjectContext(rootCmd, uii, projectPath, clean, plan, skipGitPlanPulling)
	}

	repositoryCtxProvider := func() bool {
		return getRepositoryContext(projectPath)
	}

	return rootCmd, ctxProvider, repositoryCtxProvider
}

func Execute(stdout, stderr io.Writer) {
	rootCmd, uii, err := initializedRoot(stdout, stderr)
	if err != nil {
		telemetry.TraceError(stdcontext.Background(), "init", err)
		if uii != nil {
			checkError(uii, err)
		}
		fmt.Printf("failed to initialize with error: %s", err)
		return
	}

	if err := rootCmd.Execute(); err != nil {
		telemetry.TraceError(
			stdcontext.Background(),
			"execute",
			err,
		)

		checkError(uii, err)
	}
}

func initializedRootFromArgs(stdout, stderr io.Writer, args []string) (*cobra.Command, *ui.UI, error) {
	uii := ui.Create(stdout, stderr)

	rootCmd, ctxProvider, isInRepoContext := newRoot(uii)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	// Parses falgs early such that we can build PersistentFlags on rootCmd used
	// for building various subcommands in both run and ls. This is required otherwise
	// Run and LS will not get closured variables from contextProvider
	rootCmd.ParseFlags(args)

	rootCmd.AddCommand(newExtCmd())

	if isInRepoContext() {
		runCmd, err := newRun(uii, ctxProvider)
		if err != nil {
			return nil, nil, err
		}
		rootCmd.AddCommand(
			newDocumentation(uii, ctxProvider),
			newCompletion(uii),
			newGet(uii, ctxProvider),
			newGitPlan(uii, ctxProvider),
			newHas(uii, ctxProvider),
			newLs(uii, ctxProvider),
			newPlan(uii, ctxProvider),
			runCmd,
			newPrepare(uii, ctxProvider),
			newTemplate(uii, ctxProvider),
			newVersion(uii),
			newConfig(uii, ctxProvider),
			newTelemetry(uii),
		)

		return rootCmd, uii, nil
	} else {
		rootCmd.AddCommand(
			newNoContextRun(),
			newCompletion(uii),
			newVersion(uii),
			newTelemetry(uii),
			newHas(uii, ctxProvider),
			newConfig(uii, ctxProvider),
		)

		return rootCmd, uii, nil
	}
}

func initializedRoot(out, err io.Writer) (*cobra.Command, *ui.UI, error) {
	return initializedRootFromArgs(out, err, os.Args[1:])
}

type contextProvider func() (config.ShuttleProjectContext, error)
type repositoryContext func() bool

func getProjectContext(
	rootCmd *cobra.Command,
	uii *ui.UI,
	projectPath string,
	clean bool,
	plan string,
	skipGitPlanPulling bool,
) (config.ShuttleProjectContext, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
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
	projectContext, err := c.Setup(
		fullProjectPath,
		uii,
		clean,
		skipGitPlanPulling,
		plan,
		projectFlagSet,
	)
	if err != nil {
		return config.ShuttleProjectContext{}, err
	}

	ctx := stdcontext.Background()
	taskActions, err := executer.List(
		ctx,
		uii,
		fmt.Sprintf("%s/shuttle.yaml", projectContext.ProjectPath),
		&c,
	)
	if err != nil {
		return config.ShuttleProjectContext{}, err
	}

	for name, action := range taskActions.Actions {
		args := make([]config.ShuttleScriptArgs, 0)

		for _, taskArg := range action.Args {
			args = append(args, config.ShuttleScriptArgs{
				Name:     taskArg.Name,
				Required: true,
			})
		}

		c.Scripts[name] = config.ShuttlePlanScript{
			Description: name,
			Actions: []config.ShuttleAction{
				{
					Task: name,
				},
			},
			Args: args,
		}
	}

	return c, nil
}

// getRepositoryContext makes sure that we're in a repository context, this is useful to add extra commands, which are only useful when in a repository with a shuttle file
func getRepositoryContext(projectPath string) bool {
	if projectPath != "" && projectPath != "." {
		return shuttleFileExists(projectPath, fileExists)
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		fullProjectPath := path.Join(dir, projectPath)
		exists := shuttleFileExistsRecursive(fullProjectPath, fileExists)
		return exists
	}
}

type fileExistsFunc func(filePath string) bool

// shuttleFileExistsRecursive tries to find a shuttle file by going towards the root the path, it will check each folder towards the root.
func shuttleFileExistsRecursive(projectPath string, existsFunc fileExistsFunc) bool {
	if strings.Contains(projectPath, "/") {
		exists := shuttleFileExists(projectPath, existsFunc)
		if exists {
			return true
		}

		parentProjectDir := path.Dir(projectPath)
		if parentProjectDir == projectPath {
			return false
		}

		return shuttleFileExistsRecursive(parentProjectDir, existsFunc)
	}

	return shuttleFileExists(projectPath, existsFunc)

}

// shuttleFileExists will check the given directory and return if a shuttle.yaml file is found
func shuttleFileExists(projectPath string, existsFunc fileExistsFunc) bool {
	shuttleFile := path.Join(projectPath, "shuttle.yaml")
	return existsFunc(shuttleFile)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return true
}
