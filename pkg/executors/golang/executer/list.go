package executer

import (
	"context"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/ui"
)

type TaskArg struct{}

func List(
	ctx context.Context,
	ui *ui.UI,
	path string,
	c *config.ShuttleProjectContext,
) (map[string]TaskArg, error) {
	binaries, err := prepare(ctx, ui, path, c)
	if err != nil {
		return nil, err
	}

	localInquire, err := inquire(ctx, &binaries.Local)
	if err != nil {
		return nil, err
	}
	planInquire, err := inquire(ctx, &binaries.Plan)
	if err != nil {
		return nil, err
	}

	combinedOptions := make(map[string]TaskArg, 0)
	for _, cmd := range localInquire {
		combinedOptions[cmd] = struct{}{}
	}
	for _, cmd := range planInquire {
		combinedOptions[cmd] = struct{}{}
	}

	return combinedOptions, nil
}
