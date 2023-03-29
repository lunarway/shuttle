package executer

import (
	"context"
	"fmt"

	"github.com/kjuulh/shuttletask/pkg/compile"
	"github.com/kjuulh/shuttletask/pkg/discover"
)

func prepare(ctx context.Context, path string) (*compile.Binaries, error) {
	disc, err := discover.Discover(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to fiscover shuttletask: %v", err)
	}

	binaries, err := compile.Compile(ctx, disc)
	if err != nil {
		return nil, fmt.Errorf("failed to compile binaries: %v", err)
	}

	return binaries, nil
}
