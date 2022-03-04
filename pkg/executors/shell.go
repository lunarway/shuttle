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
	"syscall"
	"time"

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

	// Set process group ID so the cmd and all its children become a new process
	// group. This allows Stop to SIGTERM the cmd's process group without killing
	// this process (i.e. this code here). Note that this does not support
	// windows. this is copied from go-cmd:
	// https://github.com/go-cmd/cmd/blob/9e40bcc1acc0e559af2c2e99ae5674ae25cde397/cmd_darwin.go#L15
	execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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
	defer func() {
		wg.Wait()
		// Do we actually see a race here????
		time.Sleep(1 * time.Second)
	}()

	go stdScanner(&wg, stderr, context.ScriptContext.Project.UI.Errorln)
	go stdScanner(&wg, stdout, context.ScriptContext.Project.UI.Infoln)

	err = waitOrStop(ctx, execCmd, 5*time.Second)
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

// waitOrStop waits for the already-started command cmd by calling its Wait
// method.
//
// If cmd does not return before ctx is done, waitOrStop sends it the given
// interrupt signal. If killDelay is positive, waitOrStop waits that additional
// period for Wait to return before sending os.Kill.
//
// This function is inspired by
// https://github.com/golang/go/blob/cacac8bdc5c93e7bc71df71981fdf32dded017bf/src/cmd/go/script_test.go#L1091
// and adapted to a non-test context along with handling group termination.
func waitOrStop(ctx context.Context, cmd *exec.Cmd, killDelay time.Duration) error {
	interruptChan := make(chan error)
	go func() {
		select {
		case interruptChan <- nil:
			return
		case <-ctx.Done():
		}

		// Signal the process group (-pid), not just the process, so that the
		// process and all its children are signaled. Else, child procs can keep
		// running and keep the stdout/stderr fd open and cause cmd.Wait to hang.
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		if err == nil {
			err = ctx.Err() // Report ctx.Err() as the reason we interrupted.
		} else if err.Error() == "os: process already finished" {
			interruptChan <- nil
			return
		}

		if killDelay > 0 {
			timer := time.NewTimer(killDelay)
			select {
			// Report ctx.Err() as the reason we interrupted the process...
			case interruptChan <- ctx.Err():
				timer.Stop()
				return
			// ...but after killDelay has elapsed, fall back to a stronger signal.
			case <-timer.C:
			}

			// Wait still hasn't returned.
			// Kill the process harder to make sure that it exits.
			//
			// Ignore any error: if cmd.Process has already terminated, we still
			// want to send ctx.Err() (or the error from the Interrupt call)
			// to properly attribute the signal that may have terminated it.
			_ = cmd.Process.Kill()
		}

		interruptChan <- err
	}()

	waitErr := cmd.Wait()

	interruptErr := <-interruptChan
	if interruptErr != nil {
		return interruptErr
	}

	return waitErr
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
