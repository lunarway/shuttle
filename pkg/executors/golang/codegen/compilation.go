package codegen

import (
	"context"
	"log"
	"os/exec"
	"path"
)

func CompileBinary(ctx context.Context, shuttlelocaldir string) (string, error) {
	cmd := exec.Command("go", "build")
	cmd.Dir = path.Join(shuttlelocaldir, "tmp")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s\n", string(output))
		return "", err
	}

	return path.Join(shuttlelocaldir, "tmp", "shuttletask"), nil
}
