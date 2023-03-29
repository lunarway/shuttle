package executer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kjuulh/shuttletask/pkg/compile"
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

	for _, cmd := range localInquire {
		if cmd == cmdToExecute {
			return executeBinaryAction(ctx, &binaries.Local, args...)
		}
	}

	for _, cmd := range planInquire {
		if cmd == cmdToExecute {
			return executeBinaryAction(ctx, &binaries.Plan, args...)
		}
	}

	combinedOptions := make(map[string]struct{}, 0)
	for _, cmd := range localInquire {
		combinedOptions[cmd] = struct{}{}
	}
	for _, cmd := range planInquire {
		combinedOptions[cmd] = struct{}{}
	}

	return fmt.Errorf("no action available in commds, available options are: %s", combinedOptions)
}

func executeBinaryAction(ctx context.Context, binary *compile.Binary, args ...string) error {
	execmd := exec.Command(binary.Path, args...)

	workdir, err := os.Getwd()
	if err != nil {
		return err
	}
	execmd.Env = append(execmd.Env, fmt.Sprintf("TASK_CONTEXT_DIR=%s", workdir))

	output, err := execmd.CombinedOutput()
	log.Printf("%s\n", string(output))
	if err != nil {
		return err
	}

	return nil
}

func inquire(ctx context.Context, binary *compile.Binary) (actions []string, err error) {
	if binary == nil {
		return []string{}, nil
	}

	if binary.Path == "" {
		return []string{}, nil
	}

	cmd := exec.Command(binary.Path, "lsjson")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("inquire failed and could not get a list of commands: %v", err)
	}

	if err := json.Unmarshal(output, &actions); err != nil {
		return nil, fmt.Errorf("inquire failed with json unmarshal: %v", err)
	}

	return actions, nil
}
