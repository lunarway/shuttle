package extensions

import "context"

type ExtensionsManager struct {
	registry string
}

func NewExtensionsManager(registry string) *ExtensionsManager {
	return &ExtensionsManager{
		registry: registry,
	}
}

func (e *ExtensionsManager) Init(ctx context.Context) error {
	return nil
}

func (e *ExtensionsManager) GetAll(ctx context.Context) ([]Extension, error) {
	return nil, nil
}

func (e *ExtensionsManager) Install(ctx context.Context) error {
	return nil
}

func (e *ExtensionsManager) Update(ctx context.Context) error {
	return nil
}

type Extension struct{}
