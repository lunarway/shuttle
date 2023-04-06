package discover_test

import (
	"context"
	"os/exec"
	"testing"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
	"github.com/stretchr/testify/assert"
)

func TestDiscover(t *testing.T) {
	discovered, err := discover.Discover(context.Background(), "testdata/simple/shuttle.yaml", &config.ShuttleProjectContext{})
	assert.NoError(t, err)

	assert.Equal(t, discover.Discovered{
		Local: &discover.ShuttleTaskDiscovered{
			Files: []string{
				"build.go",
				"download.go",
			},
			DirPath:   "testdata/simple/shuttletask",
			ParentDir: "testdata/simple",
		},
	}, *discovered)
}

func TestDiscoverComplex(t *testing.T) {
	shuttleCmd := exec.Command("shuttle", "ls")
	shuttleCmd.Dir = "testdata/child/"

	output, err := shuttleCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("shuttle ls: %s", string(output))
	}

	discovered, err := discover.Discover(context.Background(), "testdata/child/shuttle.yaml", &config.ShuttleProjectContext{
		Config: config.ShuttleConfig{
			Plan: ".shuttle/plan",
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, discover.Discovered{
		Local: &discover.ShuttleTaskDiscovered{
			Files: []string{
				"build.go",
				"download.go",
			},
			DirPath:   "testdata/child/shuttletask",
			ParentDir: "testdata/child",
		},
		Plan: &discover.ShuttleTaskDiscovered{
			Files: []string{
				"build.go",
				"download.go",
			},
			DirPath:   "testdata/child/.shuttle/plan/shuttletask",
			ParentDir: "testdata/child/.shuttle/plan",
		},
	}, *discovered)
}
