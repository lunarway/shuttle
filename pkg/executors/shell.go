package executors

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/errors"

	go_cmd "github.com/go-cmd/cmd"
)

// Build builds the docker image from a shuttle plan
func executeShell(context ActionExecutionContext) error {
	//log.Printf("Exec: %s", context.Action.Shell)
	//cmdAndArgs := strings.Split(s.Shell, " ")
	//cmd := cmdAndArgs[0]
	//args := cmdAndArgs[1:]
	shuttlePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	cmdOptions := go_cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	execCmd := go_cmd.NewCmdOptions(cmdOptions, "sh", "-c", "cd '"+context.ScriptContext.Project.ProjectPath+"'; "+context.Action.Shell)

	//execCmd := exec.Command("sh", "-c", context.Action.Shell)
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

	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
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

	// Run and wait for Cmd to return, discard Status
	context.ScriptContext.Project.UI.Titleln("shell: %s", context.Action.Shell)
	status := <-execCmd.Start()
	<-doneChan

	if status.Exit > 0 {
		return errors.NewExitCode(4, "Failed executing script `%s`: shell script `%s`\nExit code: %v", context.ScriptContext.ScriptName, context.Action.Shell, status.Exit)
	}
	return nil
}

func init() {
	addExecutor(func(action config.ShuttleAction) bool {
		return action.Shell != ""
	}, executeShell)
}
