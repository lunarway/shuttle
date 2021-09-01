package executors

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/lunarway/shuttle/pkg/config"
	shuttleerrors "github.com/lunarway/shuttle/pkg/errors"
)

// Build builds the docker image from a shuttle plan
func executeShell(ctx context.Context, context ActionExecutionContext) error {
	execCmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("cd '%s'; %s", context.ScriptContext.Project.ProjectPath, context.Action.Shell))
	context.ScriptContext.Project.UI.Verboseln("Starting shell command: %+v", execCmd.String())

	setupCommandEnvironmentVariables(execCmd, context)

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("get stdout pipe: %w", err)
	}

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("get stderr pipe: %w", err)
	}

	err = execCmd.Start()
	if err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// unusual use of sync.WaitGroup here. The os/exec.Cmd#Wait method will wait
	// for the stdout and stderr streams to be read but to be sure that the
	// execution does not move on before our scanner Go routines have returned we
	// add this wait.
	defer wg.Wait()

	go stdScanner(&wg, stderr, context.ScriptContext.Project.UI.Errorln)
	go stdScanner(&wg, stdout, context.ScriptContext.Project.UI.Infoln)

	err = execCmd.Wait()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode := exitError.ExitCode()
			if exitCode > 0 {
				return shuttleerrors.NewExitCode(4, "Failed executing script `%s`: shell script `%s`\nExit code: %v", context.ScriptContext.ScriptName, context.Action.Shell, exitCode)
			}
		}
		return err
	}

	return nil
}

func setupCommandEnvironmentVariables(execCmd *exec.Cmd, context ActionExecutionContext) {
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

func stdScanner(wg *sync.WaitGroup, reader io.ReadCloser, output func(format string, args ...interface{})) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)

	// increase max line capacity of the buffer to allow for reading very long
	// lines. This sets the maximum to 512k characters per line.
	const maxCapacity = 512 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()
		output("%s", line)
	}
}

func init() {
	addExecutor(func(action config.ShuttleAction) bool {
		return action.Shell != ""
	}, executeShell)
}
