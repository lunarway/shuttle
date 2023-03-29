package shuttlefolder

import (
	"context"
	"os"
	"path"

	"github.com/kjuulh/shuttletask/pkg/discover"
	cp "github.com/otiai10/copy"
)

func CopyFiles(
	ctx context.Context,
	shuttlelocaldir string,
	shuttletask *discover.ShuttleTaskDiscovered,
) error {
	tmpdir := path.Join(shuttlelocaldir, "tmp")

	return cp.Copy(shuttletask.DirPath, tmpdir)
}

func Move(src, dest string) error {
	return os.Rename(src, dest)
}

func GenerateTmpDir(ctx context.Context, shuttlelocaldir string) error {
	if err := os.MkdirAll(shuttlelocaldir, 0755); err != nil {
		return err
	}

	binarydir := path.Join(shuttlelocaldir, "binaries")
	if err := os.RemoveAll(binarydir); err != nil {
		return nil
	}
	if err := os.MkdirAll(binarydir, 0755); err != nil {
		return err
	}

	tmpdir := path.Join(shuttlelocaldir, "tmp")
	if err := os.RemoveAll(tmpdir); err != nil {
		return nil
	}
	if err := os.MkdirAll(tmpdir, 0755); err != nil {
		return err
	}

	return nil
}
