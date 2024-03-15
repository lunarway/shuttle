package extensions

import (
	"context"
	"fmt"

	"github.com/lunarway/shuttle/internal/global"
)

// ExtensionsManager is the entry into installing, updating and using extensions. It is the orchestrator of all the parts that consist of extensions
type ExtensionsManager struct {
	globalStore *global.GlobalStore
}

func NewExtensionsManager(globalStore *global.GlobalStore) *ExtensionsManager {
	return &ExtensionsManager{
		globalStore: globalStore,
	}
}

// Init will initialize a repository with a sample extension package
func (e *ExtensionsManager) Init(ctx context.Context) error {
	return nil
}

// GetAll will return all known and installed extensions
func (e *ExtensionsManager) GetAll(ctx context.Context) ([]Extension, error) {
	registry := getRegistryPath(e.globalStore)

	index := newRegistryIndex(registry)

	registryExtensions, err := index.getExtensions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to install extensions, could not get extensions from index: %w", err)
	}

	extensions := make([]Extension, 0)
	for _, registryExtension := range registryExtensions {
		registryExtension := registryExtension

		extension, err := newExtensionFromRegistry(e.globalStore, &registryExtension)
		if err != nil {
			return nil, err
		}

		if extension != nil {
			extensions = append(extensions, *extension)
		}
	}

	return extensions, nil
}

// Install will ensure that all known extensions are installed and ready for use
func (e *ExtensionsManager) Install(ctx context.Context) error {
	registry := getRegistryPath(e.globalStore)
	index := newRegistryIndex(registry)
	extensions, err := index.getExtensions(ctx)
	if err != nil {
		return fmt.Errorf("failed to install extensions, could not get extensions from index: %w", err)
	}

	for _, registryExtension := range extensions {
		extension, err := newExtensionFromRegistry(e.globalStore, &registryExtension)
		if err != nil {
			return err
		}

		if err := extension.Ensure(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Update will fetch the latest extensions from a registry and install them afterwards so that they're ready for use
func (e *ExtensionsManager) Update(ctx context.Context, registry string) error {
	reg, err := NewRegistryFromCombined(registry, e.globalStore)
	if err != nil {
		return fmt.Errorf("failed to update extensions: %w", err)
	}

	if err := reg.Update(ctx); err != nil {
		return err
	}

	if err := e.Install(ctx); err != nil {
		return err
	}

	return nil
}

func (e *ExtensionsManager) Publish(ctx context.Context, version string) error {
	extensionsFile, err := getExtensionsFile(ctx)
	if err != nil {
		return err
	}

	registry, err := NewRegistry("github", "", e.globalStore)
	if err != nil {
		return err
	}

	if err := registry.Publish(ctx, extensionsFile, version); err != nil {
		return err
	}

	return nil
}
