package codegen

import (
	"context"
	"log"
	"os/exec"
	"path"
)

func Format(ctx context.Context, shuttlelocaldir string) error {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = path.Join(shuttlelocaldir, "tmp")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s\nerror: %v", string(output), err)
		return err
	}

	return nil
}
