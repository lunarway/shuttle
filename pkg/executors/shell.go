package executors

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lunarway/shuttle/pkg/config"

	go_cmd "github.com/go-cmd/cmd"
)

// Build builds the docker image from a shuttle plan
func executeShell(context ActionExecutionContext) {
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

	go func() {
		for {
			select {
			case line := <-execCmd.Stdout:
				context.ScriptContext.Project.UI.InfoLn(line)
			case line := <-execCmd.Stderr:
				context.ScriptContext.Project.UI.ErrorLn(line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	context.ScriptContext.Project.UI.TitleLn("shell: %s", context.Action.Shell)
	status := <-execCmd.Start()

	// Cmd has finished but wait for goroutine to print all lines
	for len(execCmd.Stdout) > 0 || len(execCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	if status.Exit > 0 {
		context.ScriptContext.Project.UI.ExitWithErrorCode(4, "Failed executing shell script `%s`\nExit code: %v", context.Action.Shell, status.Exit)
	}
}

func init() {
	addExecutor(func(action config.ShuttleAction) bool {
		return action.Shell != ""
	}, executeShell)
}
