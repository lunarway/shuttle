package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuttleFileExists(t *testing.T) {
	t.Parallel()

	t.Run("full path, with file", func(t *testing.T) {
		actual := shuttleFileExists("/some/long/path", func(filePath string) bool {
			switch filePath {
			case "/some/long/path/shuttle.yaml":
				return true
			default:
				pathNotExpected(t, filePath)
				return false
			}
		})

		assert.True(t, actual)
	})

	t.Run("full path, no file", func(t *testing.T) {
		actual := shuttleFileExists("/some/long/path", func(filePath string) bool {
			switch filePath {
			case "/some/long/path/shuttle.yaml":
				return false
			default:
				pathNotExpected(t, filePath)
				return true
			}
		})

		assert.False(t, actual)
	})

	t.Run("current path, with file", func(t *testing.T) {
		actual := shuttleFileExists(".", func(filePath string) bool {
			switch filePath {
			case "shuttle.yaml":
				return true
			default:
				pathNotExpected(t, filePath)
				return false
			}
		})

		assert.True(t, actual)
	})

	t.Run("current path, no file", func(t *testing.T) {
		actual := shuttleFileExists(".", func(filePath string) bool {
			switch filePath {
			case "shuttle.yaml":
				return false
			default:
				pathNotExpected(t, filePath)
				return true
			}
		})

		assert.False(t, actual)
	})
}

func TestShuttleFileExistsRecursive(t *testing.T) {
	t.Parallel()

	t.Run("full path, file in given path", func(t *testing.T) {
		actual := shuttleFileExistsRecursive("/some/long/path", func(filePath string) bool {
			switch filePath {
			case "/some/long/path/shuttle.yaml":
				return true
			default:
				pathNotExpected(t, filePath)
				return false
			}
		})

		assert.True(t, actual)
	})

	t.Run("full path, file in sub directory", func(t *testing.T) {
		actual := shuttleFileExistsRecursive("/some/long/path", func(filePath string) bool {
			switch filePath {
			case "/some/long/path/shuttle.yaml":
				return false
			case "/some/long/shuttle.yaml":
				return true
			default:
				pathNotExpected(t, filePath)
				return false
			}
		})

		assert.True(t, actual)
	})

	t.Run("full path, file in root", func(t *testing.T) {
		actual := shuttleFileExistsRecursive("/some/long/path", func(filePath string) bool {
			switch filePath {
			case "/some/long/path/shuttle.yaml":
				return false
			case "/some/long/shuttle.yaml":
				return false
			case "/some/shuttle.yaml":
				return false
			case "/shuttle.yaml":
				return true
			default:
				pathNotExpected(t, filePath)
				return false
			}
		})

		assert.True(t, actual)
	})

	t.Run("full path, file not found", func(t *testing.T) {
		actual := shuttleFileExistsRecursive("/some/long/path", func(filePath string) bool {
			switch filePath {
			case "/some/long/path/shuttle.yaml":
				return false
			case "/some/long/shuttle.yaml":
				return false
			case "/some/shuttle.yaml":
				return false
			case "/shuttle.yaml":
				return false
			case "/":
				return false
			default:
				pathNotExpected(t, filePath)
				return true
			}
		})

		assert.False(t, actual)
	})

	t.Run("empty path, file false", func(t *testing.T) {
		actual := shuttleFileExistsRecursive("", func(filePath string) bool {
			switch filePath {
			case "shuttle.yaml":
				return false
			default:
				pathNotExpected(t, filePath)
				return false
			}
		})

		assert.False(t, actual)
	})

	t.Run("current dir, file found", func(t *testing.T) {
		actual := shuttleFileExistsRecursive(".", func(filePath string) bool {
			switch filePath {
			case "shuttle.yaml":
				return true
			default:
				pathNotExpected(t, filePath)
				return false
			}
		})

		assert.True(t, actual)
	})
}

func pathNotExpected(t *testing.T, filePath string) {
	t.Helper()

	assert.Fail(t, "path was not expected", "the path %s was not expected in matcher", filePath)
}
