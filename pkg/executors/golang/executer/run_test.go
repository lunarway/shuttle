package executer_test

import (
	"context"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors/golang/executer"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestRunVersion(t *testing.T) {
	updateShuttle(t, "testdata/child")
	ctx := context.Background()

	c := &config.ShuttleProjectContext{Config: config.ShuttleConfig{Plan: "something"}}

	ui := ui.Create(os.Stdout, os.Stderr)

	err := executer.Run(ctx, ui, c, "testdata/child/shuttle.yaml", "version")
	assert.NoError(t, err)

	err = executer.Run(ctx, ui, c, "testdata/child/shuttle.yaml", "build")
	assert.NoError(t, err)

	err = executer.Run(
		ctx,
		ui,
		c,
		"testdata/child/shuttle.yaml",
		"build",
		"--some-unexisting-arg",
		"something",
	)
	assert.Error(t, err)
}

func updateShuttle(t *testing.T, path string) {
	err := os.RemoveAll("testdata/child/.shuttle/")
	assert.NoError(t, err)

	shuttleCmd := exec.Command("shuttle", "ls")
	shuttleCmd.Dir = path
	if output, err := shuttleCmd.CombinedOutput(); err != nil {
		log.Printf("%s\n", string(output))
		assert.Error(t, err)
	}
}
