package codegen

import (
	"context"
	"os/exec"
	"path"

	"github.com/lunarway/shuttle/pkg/ui"
)

func ModTidy(ctx context.Context, ui *ui.UI, shuttlelocaldir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = path.Join(shuttlelocaldir, "tmp")

	output, err := cmd.CombinedOutput()
	if err != nil {
		ui.Errorln("mod tidy: %s, error: %v", string(output), err)
		return err
	}

	return nil
}
