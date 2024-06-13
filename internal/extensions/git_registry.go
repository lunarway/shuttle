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

// Publish isn't implemented yet for gitRegistry
func (*gitRegistry) Publish(ctx context.Context, extFile *shuttleExtensionsFile, version string) error {
	panic("unimplemented")
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

	return g.executeGit(ctx, "git", gitOptions{
		args: []string{
			"pull",
		},
		dir: registry,
	})
}

func (g *gitRegistry) cloneGitRepository(ctx context.Context) error {
	registry := getRegistryPath(g.globalStore)

	return g.executeGit(ctx, "git", gitOptions{
		args: []string{
			"clone",
			g.url,
			registry,
		},
	})

}

func (g *gitRegistry) registryClonedAlready() bool {
	registry := getRegistryPath(g.globalStore)

	return exists(path.Join(registry, ".git"))
}

type gitOptions struct {
	args []string
	dir  string
}

func (g *gitRegistry) executeGit(ctx context.Context, name string, gitOptions gitOptions) error {
	cmd := exec.CommandContext(ctx, name, gitOptions.args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if gitOptions.dir != "" {
		cmd.Dir = gitOptions.dir
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
