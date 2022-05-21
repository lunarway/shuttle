package executors

import (
	"context"

	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-cmd/cmd"
	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/errors"
)

func ShellExecutor(action config.ShuttleAction) (Executor, bool) {
	return executeShell, action.Shell != ""
}

// Build builds the docker image from a shuttle plan
func executeShell(ctx context.Context, context ActionExecutionContext) error {
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
		// support large outputs from scripts
		LineBufferSize: 512e3,
	}

	cmdArgs := []string{"-c", fmt.Sprintf("cd '%s'; %s", context.ScriptContext.Project.ProjectPath, context.Action.Shell)}
	execCmd := cmd.NewCmdOptions(cmdOptions, "sh", cmdArgs...)
	context.ScriptContext.Project.UI.Verboseln("Starting shell command: %s %s", execCmd.Name, strings.Join(cmdArgs, " "))

	setupCommandEnvironmentVariables(execCmd, context)

	outputReadCompleted := make(chan struct{})

	go func() {
		defer close(outputReadCompleted)

		for execCmd.Stdout != nil || execCmd.Stderr != nil {
			select {
			case line, open := <-execCmd.Stdout:
				if !open {
					execCmd.Stdout = nil
					continue
				}
				context.ScriptContext.Project.UI.Infoln("%s", line)
			case line, open := <-execCmd.Stderr:
				if !open {
					execCmd.Stderr = nil
					continue
				}
				context.ScriptContext.Project.UI.Errorln("%s", line)
			}
		}
	}()

	// stop cmd if context is cancelled
	go func() {
		select {
		case <-ctx.Done():
			err := execCmd.Stop()
			if err != nil {
				context.ScriptContext.Project.UI.Errorln("Failed to stop script '%s': %v", context.Action.Shell, err)
			}
		case <-outputReadCompleted:
		}
	}()

	select {
	case status := <-execCmd.Start():
		<-outputReadCompleted
		if status.Exit > 0 {
			return errors.NewExitCode(4, "Failed executing script `%s`: shell script `%s`\nExit code: %v", context.ScriptContext.ScriptName, context.Action.Shell, status.Exit)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func setupCommandEnvironmentVariables(execCmd *cmd.Cmd, context ActionExecutionContext) {
	shuttlePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	execCmd.Env = os.Environ()
	for name, value := range context.ScriptContext.Args {
		execCmd.Env = append(execCmd.Env, fmt.Sprintf("%s=%s", name, value))
	}
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("plan=%s", context.ScriptContext.Project.LocalPlanPath))
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("tmp=%s", context.ScriptContext.Project.TempDirectoryPath))
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("project=%s", context.ScriptContext.Project.ProjectPath))
	// TODO: Add project path as a shuttle specific ENV
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("PATH=%s", shuttlePath+string(os.PathListSeparator)+os.Getenv("PATH")))
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("SHUTTLE_PLANS_ALREADY_VALIDATED=%s", context.ScriptContext.Project.LocalPlanPath))
}
