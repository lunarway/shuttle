package executer

import (
	"context"
	"fmt"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors/golang/compile"
	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
)

func prepare(ctx context.Context, path string, c *config.ShuttleProjectContext) (*compile.Binaries, error) {
	disc, err := discover.Discover(ctx, path, c)
	if err != nil {
		return nil, fmt.Errorf("failed to fiscover shuttletask: %v", err)
	}

	binaries, err := compile.Compile(ctx, disc)
	if err != nil {
		return nil, fmt.Errorf("failed to compile binaries: %v", err)
	}

	return binaries, nil
}
