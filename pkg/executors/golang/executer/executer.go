package executer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/lunarway/shuttle/pkg/executors/golang/compile"
	"github.com/lunarway/shuttle/pkg/telemetry"
)

// Executes an action based on which plan is used
// Get a list of actions for each binary if they exist
// Take child if available otherwise pick plan, else error
func executeAction(ctx context.Context, binaries *compile.Binaries, args ...string) error {
	localInquire, err := inquire(ctx, &binaries.Local)
	if err != nil {
		return err
	}
	planInquire, err := inquire(ctx, &binaries.Plan)
	if err != nil {
		return err
	}

	cmdToExecute := args[0]

	ran, err := localInquire.Execute(cmdToExecute, func() error {
		return executeBinaryAction(ctx, &binaries.Local, args...)
	})
	if err != nil {
		return err
	}
	if ran {
		return nil
	}

	ran, err = planInquire.Execute(cmdToExecute, func() error {
		return executeBinaryAction(ctx, &binaries.Plan, args...)
	})
	if err != nil {
		return err
	}
	if ran {
		return nil
	}

	return fmt.Errorf("no action available in commands, available options are available through shuttle run -h")
}

func executeBinaryAction(ctx context.Context, binary *compile.Binary, args ...string) error {
	execmd := exec.Command(binary.Path, args...)
	execmd.Stdout = os.Stdout
	execmd.Stderr = os.Stderr

	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	execmd.Env = os.Environ()
	execmd.Env = append(execmd.Env, fmt.Sprintf("TASK_CONTEXT_DIR=%s", workdir))
	execmd.Env = append(execmd.Env, "SHUTTLE_INTERACTIVE=default")
	execmd.Env = append(
		execmd.Env,
		fmt.Sprintf("%s=%s",
			"SHUTTLE_CONTEXT_ID",
			telemetry.ContextIDFrom(ctx),
		),
	)

	err = execmd.Run()

	os.Stdout.Sync()
	os.Stderr.Sync()

	if err != nil {
		return err
	}

	return nil
}

func inquire(ctx context.Context, binary *compile.Binary) (actions *Actions, err error) {
	if binary == nil {
		return nil, nil
	}

	if binary.Path == "" {
		return nil, nil
	}

	cmd := exec.Command(binary.Path, "lsjson")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("shuttle_actions: err: %s", string(exitErr.Stderr))
		}

		fmt.Printf("shuttle_actions: %s", string(output))
		return nil, fmt.Errorf("inquire failed and could not get a list of commands: %v", err)

	}

	actions = &Actions{}
	if err := json.Unmarshal(output, actions); err != nil {
		return nil, fmt.Errorf("inquire failed with json unmarshal: %v", err)
	}

	return actions, nil
}
