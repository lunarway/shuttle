package executer

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/lunarway/shuttle/pkg/config"
	golangerrors "github.com/lunarway/shuttle/pkg/executors/golang/errors"
	"github.com/lunarway/shuttle/pkg/ui"
)

func List(
	ctx context.Context,
	ui *ui.UI,
	path string,
	c *config.ShuttleProjectContext,
) (*Actions, error) {
	if !isActionsEnabled() {
		ui.Verboseln("shuttle golang actions disabled")
		return NewActions(), nil
	}

	binaries, err := prepare(ctx, ui, path, c)
	if err != nil {
		if errors.Is(err, golangerrors.ErrGolangActionNoBuilder) {
			return NewActions(), nil
		}
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

	actions := NewActions().
		Merge(localInquire).
		Merge(planInquire)

	return actions, nil
}

const defaultGolangActions = false

func isActionsEnabled() bool {
	enabledRaw := os.Getenv("SHUTTLE_GOLANG_ACTIONS")
	if enabledRaw == "" {
		return defaultGolangActions
	}

	enabled, err := strconv.ParseBool(enabledRaw)
	if err != nil {
		return defaultGolangActions
	}

	return enabled

}
