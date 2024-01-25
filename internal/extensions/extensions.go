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
	return nil, nil
}

// Install will ensure that all known extensions are installed and ready for use
func (e *ExtensionsManager) Install(ctx context.Context) error {
	return nil
}

// Update will fetch the latest extensions from a registry and install them afterwards so that they're ready for use
func (e *ExtensionsManager) Update(ctx context.Context, registry string) error {
	reg, err := NewRegistry(registry, e.globalStore)
	if err != nil {
		return fmt.Errorf("failed to update extensions: %w", err)
	}

	if err := reg.Update(ctx); err != nil {
		return err
	}

	// 3. Initiate install

	return nil
}
