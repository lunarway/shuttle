package extensions

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

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
	if g.registryClonedAlready() {
		if err := g.fetchGitRepository(ctx); err != nil {
			return fmt.Errorf("failed to update registry: %w", err)
		}
	} else {
		registry := getRegistryPath(g.globalStore)

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
	registry := getRegistryPath(g.globalStore)

	cmd := exec.CommandContext(ctx, "git", "pull")

	cmd.Dir = registry
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func (g *gitRegistry) cloneGitRepository(ctx context.Context) error {
	registry := getRegistryPath(g.globalStore)

	cmd := exec.CommandContext(ctx, "git", "clone", g.url, registry)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func (g *gitRegistry) registryClonedAlready() bool {
	registry := getRegistryPath(g.globalStore)

	return exists(path.Join(registry, ".git"))
}
