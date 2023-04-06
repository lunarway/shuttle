package executer

import (
	"context"
	"fmt"

	"github.com/lunarway/shuttle/pkg/config"
)

func List(ctx context.Context, path string, c *config.ShuttleProjectContext) error {
	binaries, err := prepare(ctx, path, c)
	if err != nil {
		return err
	}

	localInquire, err := inquire(ctx, &binaries.Local)
	if err != nil {
		return err
	}
	planInquire, err := inquire(ctx, &binaries.Plan)
	if err != nil {
		return err
	}

	combinedOptions := make(map[string]struct{}, 0)
	for _, cmd := range localInquire {
		combinedOptions[cmd] = struct{}{}
	}
	for _, cmd := range planInquire {
		combinedOptions[cmd] = struct{}{}
	}

	println("Args: ")
	for k := range combinedOptions {
		fmt.Printf("\t%s\n", k)
	}

	return nil
}
