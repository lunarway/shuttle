package codegen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsModule(t *testing.T) {
	t.Parallel()

	sut := actionsModFile{
		content: []string{
			"someModuleName",
			"someOtherModuleName",
			"  someSpacedModule",
			"package somePackaged",
			"module someModModule",
			"require someRequireModule",
		},
	}

	t.Run("contains a module name", func(t *testing.T) {
		actual := sut.containsModule("someModuleName")

		assert.True(t, actual)
	})

	t.Run("does not contain module", func(t *testing.T) {
		actual := sut.containsModule("someNonExistingModule")

		assert.False(t, actual)
	})

	t.Run("spaced matches", func(t *testing.T) {
		actual := sut.containsModule("someSpacedModule")

		assert.True(t, actual)
	})

	t.Run("packaged matches", func(t *testing.T) {
		actual := sut.containsModule("somePackaged")

		assert.True(t, actual)
	})

	t.Run("modded module matches", func(t *testing.T) {
		actual := sut.containsModule("someModModule")

		assert.True(t, actual)
	})

	t.Run("required module matches", func(t *testing.T) {
		actual := sut.containsModule("someRequireModule")

		assert.True(t, actual)
	})

	t.Run("case sensitive doesn't match", func(t *testing.T) {
		actual := sut.containsModule("SOMEMODULENAME")

		assert.False(t, actual)
	})
}

func TestReplaceModulePath(t *testing.T) {
	t.Parallel()

	createSut := func() actionsModFile {

		modFileContent := `module actions

require (
	root_workspace v0.0.0
	subpackage v0.0.0
	othersubpackage v0.0.0
)

go 1.21.4

replace othersubpackage => ../../../../othersubpackage`

		return actionsModFile{
			path:    "/some-path/some-other-path",
			content: strings.Split(modFileContent, "\n"),
		}

	}

	t.Run("module matches, not replaced already", func(t *testing.T) {
		sut := createSut()

		expected := `module actions

require (
	root_workspace v0.0.0
	subpackage v0.0.0
	othersubpackage v0.0.0
)

go 1.21.4

replace othersubpackage => ../../../../othersubpackage

replace subpackage => ../subpackage`

		sut.replaceModulePath("some-other-path/newpath", module{
			name: "subpackage",
			path: "subpackage",
		})

		assert.Equal(t, expected, strings.Join(sut.content, "\n"))
	})

	t.Run("module matches, not replaced already, deeper nesting", func(t *testing.T) {
		sut := createSut()

		expected := `module actions

require (
	root_workspace v0.0.0
	subpackage v0.0.0
	othersubpackage v0.0.0
)

go 1.21.4

replace othersubpackage => ../../../../othersubpackage

replace subpackage => subpackage`

		sut.replaceModulePath("/some-path", module{
			name: "subpackage",
			path: "subpackage",
		})

		assert.Equal(t, expected, strings.Join(sut.content, "\n"))
	})

	t.Run("module matches, already replaced already", func(t *testing.T) {
		sut := createSut()

		expected := `module actions

require (
	root_workspace v0.0.0
	subpackage v0.0.0
	othersubpackage v0.0.0
)

go 1.21.4

replace othersubpackage => ../othersubpackage`

		sut.replaceModulePath("some-other-path/newpath", module{
			name: "othersubpackage",
			path: "othersubpackage",
		})

		assert.Equal(t, expected, strings.Join(sut.content, "\n"))
	})

	t.Run("module matches, already replaced already deeper nesting", func(t *testing.T) {
		sut := createSut()

		expected := `module actions

require (
	root_workspace v0.0.0
	subpackage v0.0.0
	othersubpackage v0.0.0
)

go 1.21.4

replace othersubpackage => othersubpackage`

		sut.replaceModulePath("/some-path", module{
			name: "othersubpackage",
			path: "othersubpackage",
		})

		assert.Equal(t, expected, strings.Join(sut.content, "\n"))
	})
}

func TestSegmentsTo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		path     string
		rootDir  string
		expected int
	}{
		{
			name:     "empty",
			path:     "",
			rootDir:  "",
			expected: 0,
		},
		{
			name:     "current dir",
			path:     "/some-dir/",
			rootDir:  "/some-dir/",
			expected: 0,
		},
		{
			name:     "one level",
			path:     "/some-dir/some-other-dir/",
			rootDir:  "/some-dir",
			expected: 1,
		},
		{
			name:     "2 level",
			path:     "/some-dir/some-other-dir/some-third-dir/",
			rootDir:  "/some-dir",
			expected: 2,
		},
		{
			name:     "1 level",
			path:     "/some-dir/some-other-dir/some-third-dir/",
			rootDir:  "/some-dir/some-other-dir",
			expected: 1,
		},
		{
			name:     "without trailing",
			path:     "/some-dir/some-third-dir",
			rootDir:  "/some-dir",
			expected: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sut := actionsModFile{
				path: testCase.path,
			}

			actual := sut.segmentsTo(testCase.rootDir)

			assert.Equal(t, testCase.expected, actual)
		})
	}

}
