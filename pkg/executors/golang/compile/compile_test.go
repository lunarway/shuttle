package compile_test

import (
	"context"
	"testing"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors/golang/compile"
	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	ctx := context.Background()
	discovered, err := discover.Discover(ctx, "testdata/simple/shuttle.yaml", &config.ShuttleProjectContext{})
	assert.NoError(t, err)

	path, err := compile.Compile(ctx, discovered)
	assert.NoError(t, err)

	assert.Contains(t, path.Local.Path, "testdata/simple/.shuttle/actions/binaries/actions-")
}
