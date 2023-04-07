package executors

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-cmd/cmd"
	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors/golang/executer"
)

func TaskExecutor(action config.ShuttleAction) (Executor, bool) {

	return executeTask, action.Task != ""
}

func executeTask(ctx context.Context, context ActionExecutionContext) error {
	context.ScriptContext.Project.UI.Verboseln("Starting task command: %s", context.Action.Task)

	args := make([]string, 0)
	args = append(args, context.Action.Task)
	for name, value := range context.ScriptContext.Args {
		args = append(args, fmt.Sprintf("--%s", name))
		args = append(args, value)
	}

	err := executer.Run(ctx, &context.ScriptContext.Project, "shuttle.yaml", args...)
	if err != nil {
		return err
	}

	return nil
}

func setupTaskCommandEnvironmentVariables(execCmd *cmd.Cmd, context ActionExecutionContext) {
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
