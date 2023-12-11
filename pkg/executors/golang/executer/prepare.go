package executer

import (
	"context"
	"fmt"
	"log"
	"time"

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
	ui.Verboseln("preparing shuttle golang actions")
	start := time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	disc, err := discover.Discover(ctx, path, c)
	if err != nil {
		return nil, fmt.Errorf("failed to discover actions: %v", err)
	}

	binaries, err := compile.Compile(ctx, ui, disc)
	if err != nil {
		return nil, fmt.Errorf("failed to compile binaries: %v", err)
	}

	elapsed := time.Since(start)
	ui.Verboseln("preparing shuttle golang actions took: %d ms", elapsed.Milliseconds())

	return binaries, nil
}
