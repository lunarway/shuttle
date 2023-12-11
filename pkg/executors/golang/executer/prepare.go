package executer

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors/golang/compile"
	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
	"github.com/lunarway/shuttle/pkg/ui"
)

func prepare(
	ctx context.Context,
	ui *ui.UI,
	path string,
	c *config.ShuttleProjectContext,
) (*compile.Binaries, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	disc, err := discover.Discover(ctx, path, c)
	if err != nil {
		return nil, fmt.Errorf("failed to discover actions: %v", err)
	}

	binaries, err := NewActionProvider(newGoBinary()).GetBinaries(ctx, disc)
	if err != nil {
		return nil, fmt.Errorf("failed to compile binaries: %v", err)
	}

	return binaries, nil
}

type actionProvider struct {
	providers []ActionProvider
}

func NewActionProvider(providers ...ActionProvider) *actionProvider {
	return &actionProvider{
		providers: providers,
	}
}

func (a *actionProvider) GetBinaries(ctx context.Context, discover *discover.Discovered) (*compile.Binaries, error) {
	binaries := compile.Binaries{}

	for _, actionProvider := range a.providers {
		b, err := actionProvider.GetBinaries(ctx, discover)
		if err != nil {
			return nil, err
		}

		if b.Local != nil {
			binaries.Local = b.Local
		}
		if b.Plan != nil {
			binaries.Plan = b.Plan
		}

		if binaries.Local != nil && binaries.Plan != nil {
			return &binaries, nil
		}
	}

	return &binaries, nil
}

type ActionProvider interface {
	GetBinaries(ctx context.Context, discover *discover.Discovered) (*compile.Binaries, error)
}

type goBinary struct {
	ui *ui.UI
}

func newGoBinary() *goBinary {
	return &goBinary{}
}

var _ ActionProvider = &goBinary{}

func (g *goBinary) GetBinaries(ctx context.Context, discover *discover.Discovered) (*compile.Binaries, error) {
	binaries, err := compile.Compile(ctx, g.ui, discover)

	return binaries, err
}
