package compile_test

import (
	"context"
	"os"
	"testing"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors/golang/compile"
	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	ctx := context.Background()
	discovered, err := discover.Discover(
		ctx,
		"testdata/simple/shuttle.yaml",
		&config.ShuttleProjectContext{},
	)
	assert.NoError(t, err)

	uiout := ui.Create(os.Stdout, os.Stderr)

	path, err := compile.Compile(ctx, uiout, discovered)
	assert.NoError(t, err)

	assert.Contains(t, path.Local.Path, "testdata/simple/.shuttle/actions/binaries/actions-")
}
