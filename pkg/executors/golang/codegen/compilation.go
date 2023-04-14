package codegen

import (
	"context"
	"os/exec"
	"path"

	"github.com/lunarway/shuttle/pkg/ui"
)

func CompileBinary(ctx context.Context, ui *ui.UI, shuttlelocaldir string) (string, error) {
	cmd := exec.Command("go", "build")
	cmd.Dir = path.Join(shuttlelocaldir, "tmp")

	output, err := cmd.CombinedOutput()
	if err != nil {
		ui.Verboseln("compile-binary output: %s", string(output))
		return "", err
	}

	return path.Join(shuttlelocaldir, "tmp", "actions"), nil
}
