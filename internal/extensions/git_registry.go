package extensions

import (
	"context"
	"fmt"

	"github.com/lunarway/shuttle/internal/global"
)

// gitRegistry represents a type of registry backed by a remote git registry, whether folder or url based. it is denoted by the variable git=github.com/lunarway/shuttle-extensions.git as an example
type gitRegistry struct {
	url         string
	globalStore *global.GlobalStore
}

func (*gitRegistry) Get(ctx context.Context) error {
	panic("unimplemented")
}

func (g *gitRegistry) Update(ctx context.Context) error {
	registry := getRegistryPath(g.globalStore)

	if exists(registry) {
		if err := g.fetchGitRepository(ctx); err != nil {
			return fmt.Errorf("failed to update registry: %w", err)
		}
	} else {
		if err := ensureExists(registry); err != nil {
			return fmt.Errorf("failed to create registry path: %w", err)
		}

		if err := g.cloneGitRepository(ctx); err != nil {
			return fmt.Errorf("failed to clone registry: %w", err)
		}
	}

	return nil
}

func newGitRegistry(url string, globalStore *global.GlobalStore) Registry {
	return &gitRegistry{
		url:         url,
		globalStore: globalStore,
	}
}

func (g *gitRegistry) fetchGitRepository(ctx context.Context) error {
	panic("unimplemented")
}

func (g *gitRegistry) cloneGitRepository(ctx context.Context) error {
	panic("unimplemented")
}
