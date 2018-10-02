package executors

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lunarway/shuttle/pkg/output"

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

	go func() {
		for {
			select {
			case line := <-execCmd.Stdout:
				fmt.Println(line)
			case line := <-execCmd.Stderr:
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	status := <-execCmd.Start()

	// Cmd has finished but wait for goroutine to print all lines
	for len(execCmd.Stdout) > 0 || len(execCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	if status.Exit > 0 {
		output.ExitWithErrorCode(4, fmt.Sprintf("Failed executing shell script `%s`\nExit code: %v", context.Action.Shell, status.Exit))
	}

	// go func() {
	// 	stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	// }()

	// go func() {
	// 	stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	// }()

	// err := execCmd.Wait()
	//outStr, errStr := string(stdout), string(stderr)
	//fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	//return context.ScriptContext.ScriptName + "> Executed shell - " + context.Action.Shell
}

func init() {
	addExecutor(func(action config.ShuttleAction) bool {
		return action.Shell != ""
	}, executeShell)
}
