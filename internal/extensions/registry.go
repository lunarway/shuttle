package extensions

import (
	"context"
	"fmt"
	"strings"

	"github.com/lunarway/shuttle/internal/global"
)

// Registry represents some kind of upstream registry where extension metadata lives, such as which ones should be downloaded, which versions they're on and how to download them
type Registry interface {
	Get(ctx context.Context) error
	Update(ctx context.Context) error
	Publish(ctx context.Context, extFile *shuttleExtensionsFile, version string) error
}

// NewRegistryFromCombined is a shim for concrete implementations of the registries, such as gitRegistry
func NewRegistryFromCombined(registry string, globalStore *global.GlobalStore) (Registry, error) {
	registryType, registryUrl, ok := strings.Cut(registry, "=")
	if !ok {
		return nil, fmt.Errorf("registry was not a valid url: %s", registry)
	}

	switch registryType {
	case "git":
		return newGitRegistry(registryUrl, globalStore), nil
	default:
		return nil, fmt.Errorf("registry type was not valid: %s", registryType)
	}
}

// NewRegistry is a shim for concrete implementations of the registries, such as gitRegistry
func NewRegistry(registryType string, registryUrl string, globalStore *global.GlobalStore) (Registry, error) {
	switch registryType {
	case "git":
		return newGitRegistry(registryUrl, globalStore), nil
	case "github":
		return newGitHubRegistry()
	default:
		return nil, fmt.Errorf("registry type was not valid: %s", registryType)
	}
}
