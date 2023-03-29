package compile_test

import (
	"context"
	"testing"

	"github.com/kjuulh/shuttletask/pkg/compile"
	"github.com/kjuulh/shuttletask/pkg/discover"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	ctx := context.Background()
	discovered, err := discover.Discover(ctx, "testdata/simple/shuttle.yaml")
	assert.NoError(t, err)

	path, err := compile.Compile(ctx, discovered)
	assert.NoError(t, err)

	assert.Contains(t, path.Local.Path, "testdata/simple/.shuttle/shuttletask/binaries/shuttletask-")
}
