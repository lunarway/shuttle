package codegen

import (
	"context"
	"os/exec"
	"path"

	"github.com/lunarway/shuttle/pkg/ui"
)

func Format(ctx context.Context, ui *ui.UI, shuttlelocaldir string) error {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = path.Join(shuttlelocaldir, "tmp")

	output, err := cmd.CombinedOutput()
	if err != nil {
		ui.Errorln("go fmt: %s, error: %v", string(output), err)
		return err
	}

	return nil
}
