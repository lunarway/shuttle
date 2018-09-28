package executors

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lunarway/shuttle/pkg/output"

	"github.com/lunarway/shuttle/pkg/config"
)

// Build builds the docker image from a shuttle plan
func executeShell(context ActionExecutionContext) {
	//log.Printf("Exec: %s", context.Action.Shell)
	//cmdAndArgs := strings.Split(s.Shell, " ")
	//cmd := cmdAndArgs[0]
	//args := cmdAndArgs[1:]
	shuttlePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	execCmd := exec.Command("sh", "-c", context.Action.Shell)
	execCmd.Env = os.Environ()
	for name, value := range context.ScriptContext.Args {
		execCmd.Env = append(execCmd.Env, fmt.Sprintf("%s=%s", name, value))
	}
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("plan=%s", context.ScriptContext.Project.LocalPlanPath))
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("tmp=%s", context.ScriptContext.Project.TempDirectoryPath))
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("project=%s", context.ScriptContext.Project.ProjectPath))
	// TODO: Add project path as a shuttle specific ENV
	execCmd.Env = append(execCmd.Env, fmt.Sprintf("PATH=%s", shuttlePath+string(os.PathListSeparator)+os.Getenv("PATH")))
	execCmd.Dir = context.ScriptContext.Project.ProjectPath

	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := execCmd.StdoutPipe()
	stderrIn, _ := execCmd.StderrPipe()

	execCmd.Start()

	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	}()

	go func() {
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	}()

	err := execCmd.Wait()
	if err != nil {
		output.ExitWithErrorCode(4, fmt.Sprintf("Failed executing shell script `%s`\nError: %s", context.Action.Shell, err))
	}
	if errStdout != nil {
		output.ExitWithErrorCode(4, fmt.Sprintf("Failed to capture Stdout for shell script `%s`˜\nError: %s", context.Action.Shell, errStdout))
	}
	if errStderr != nil {
		output.ExitWithErrorCode(4, fmt.Sprintf("Failed to capture Stderr for shell script `%s`˜\nError: %s", context.Action.Shell, errStderr))
	}
	//outStr, errStr := string(stdout), string(stderr)
	//fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	//return context.ScriptContext.ScriptName + "> Executed shell - " + context.Action.Shell
}

func init() {
	addExecutor(func(action config.ShuttleAction) bool {
		return action.Shell != ""
	}, executeShell)
}
