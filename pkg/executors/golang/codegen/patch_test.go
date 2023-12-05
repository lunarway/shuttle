package codegen

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatchGoMod(t *testing.T) {
	t.Parallel()

	t.Run("finds root module adds to actions plan", func(t *testing.T) {
		err := patchGoMod("testdata/patch/root_module/", ".shuttle/actions", func(name string, contents []byte, permissions fs.FileMode) error {
			assert.Equal(t, "testdata/patch/root_module/.shuttle/actions/tmp/go.mod", name)
			assert.Equal(t, `module actions

require (
	root_module
)

go 1.21.4

replace root_module => ../../..
`, string(contents))

			return nil
		})
		require.NoError(t, err)
	})
}
