package executer

import (
	"context"
	"log"

	"github.com/lunarway/shuttle/pkg/config"
)

func Run(ctx context.Context, c *config.ShuttleProjectContext, path string, args ...string) error {
	binaries, err := prepare(ctx, path, c)
	if err != nil {
		log.Printf("failed to run command: %v", err)
		return err
	}

	if err := executeAction(ctx, binaries, args...); err != nil {
		return err
	}

	return nil
}
