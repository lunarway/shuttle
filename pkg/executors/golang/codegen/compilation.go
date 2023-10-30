package codegen

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/lunarway/shuttle/pkg/ui"
)

func CompileBinary(ctx context.Context, ui *ui.UI, shuttlelocaldir string) (string, error) {
	cmd := exec.Command("go", "build")
	cmd.Env = os.Environ()
	// We need to set workspaces off, as we don't want users to have to add the golang modules to their go.work
	cmd.Env = append(cmd.Env, "GOWORK=off")

	cmd.Dir = path.Join(shuttlelocaldir, "tmp")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("compile-binary output: %s", string(output))
		return "", err
	}

	return path.Join(shuttlelocaldir, "tmp", "actions"), nil
}
