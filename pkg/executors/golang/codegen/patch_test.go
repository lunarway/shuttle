package codegen

import (
	"context"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatchGoMod(t *testing.T) {
	t.Parallel()

	t.Run("finds root module adds to actions plan", func(t *testing.T) {
		sut := NewPatcher()
		sut.patcher = newGoModPatcher(func(name string, contents []byte, permissions fs.FileMode) error {
			assert.Equal(t, "testdata/patch/root_module/.shuttle/actions/tmp/go.mod", name)
			assert.Equal(t, `module actions

require (
	root_module
)

go 1.21.4


replace root_module => ../../..`, string(contents))

			return nil
		})

		err := sut.Patch(context.Background(), "testdata/patch/root_module/", "testdata/patch/root_module/.shuttle/actions")
		require.NoError(t, err)
	})

	t.Run("finds root module replaces existing", func(t *testing.T) {
		sut := NewPatcher()
		sut.patcher = newGoModPatcher(func(name string, contents []byte, permissions fs.FileMode) error {
			assert.Equal(t, "testdata/patch/replace_existing/.shuttle/actions/tmp/go.mod", name)
			assert.Equal(t, `module actions

require (
	replace_existing v0.0.0
)

go 1.21.4

replace replace_existing => ../../..
`, string(contents))

			return nil
		})

		err := sut.Patch(context.Background(), "testdata/patch/replace_existing/", "testdata/patch/replace_existing/.shuttle/actions")
		require.NoError(t, err)
	})

	t.Run("finds root workspace adds entries", func(t *testing.T) {
		sut := NewPatcher()
		sut.patcher = newGoModPatcher(func(name string, contents []byte, permissions fs.FileMode) error {
			assert.Equal(t, "testdata/patch/root_workspace/.shuttle/actions/tmp/go.mod", name)
			assert.Equal(t, `module actions

require (
	root_workspace v0.0.0
	subpackage v0.0.0
	othersubpackage v0.0.0
)

go 1.21.4


replace othersubpackage => ../../../other/subpackage

replace root_workspace => ../../..

replace subpackage => ../../../subpackage`, string(contents))

			return nil
		})

		err := sut.Patch(context.Background(), "testdata/patch/root_workspace/", "testdata/patch/root_workspace/.shuttle/actions")
		require.NoError(t, err)
	})
}
