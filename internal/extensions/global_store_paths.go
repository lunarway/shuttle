package extensions

import (
	"errors"
	"os"
	"path"

	"github.com/lunarway/shuttle/internal/global"
)

func getRegistryPath(globalStore *global.GlobalStore) string {
	return path.Join(globalStore.Root(), "registry")
}

func getExtensionsPath(globalStore *global.GlobalStore) string {
	return path.Join(globalStore.Root(), "extensions")
}

func getExtensionsCachePath(globalStore *global.GlobalStore) string {
	return path.Join(getExtensionsPath(globalStore), "cache")
}

func ensureExists(dirPath string) error {
	return os.MkdirAll(dirPath, 0o666)
}

func exists(dirPath string) bool {
	_, err := os.Stat(dirPath)

	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}
