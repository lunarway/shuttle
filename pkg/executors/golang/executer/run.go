package executer

import (
	"context"

	"github.com/lunarway/shuttle/pkg/config"
)

func Run(ctx context.Context, c *config.ShuttleProjectContext, path string, args ...string) error {
	binaries, err := prepare(ctx, path, c)
	if err != nil {
		return err
	}

	if err := executeAction(ctx, binaries, args...); err != nil {
		return err
	}

	return nil
}
