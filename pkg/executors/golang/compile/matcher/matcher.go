package matcher

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
	"github.com/lunarway/shuttle/pkg/ui"
	"golang.org/x/mod/sumdb/dirhash"
)

func BinaryMatches(
	ctx context.Context,
	ui *ui.UI,
	hash string,
	actions *discover.ActionsDiscovered,
) (string, bool, error) {
	shuttlebindir := path.Join(actions.ParentDir, ".shuttle/actions/binaries")

	if _, err := os.Stat(shuttlebindir); errors.Is(err, os.ErrNotExist) {
		ui.Verboseln("package doesn't exist continueing")
		return "", false, nil
	}

	entries, err := os.ReadDir(shuttlebindir)
	if err != nil {
		return "", false, err
	}

	if len(entries) == 0 {
		return "", false, err
	}

	// We only expect a single binary in the folder, so we just take the first entry if it exists
	binary := entries[0]

	expectedPath := fmt.Sprintf("actions-%s", hex.EncodeToString([]byte(hash)[:16]))
	actualName := binary.Name()
	if actualName == expectedPath {
		return path.Join(shuttlebindir, binary.Name()), true, nil
	} else {
		ui.Verboseln("binary does not match, rebuilding... (actual=%s, expected=%s)", actualName, expectedPath)
		return "", false, nil
	}
}

func GetHash(ctx context.Context, actions *discover.ActionsDiscovered) (string, error) {
	entries := make([]string, len(actions.Files))

	for i, task := range actions.Files {
		entries[i] = path.Join(actions.DirPath, task)
	}

	open := func(name string) (io.ReadCloser, error) {
		b, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}

		return io.NopCloser(bytes.NewReader(b)), nil
	}

	hash, err := dirhash.Hash1(entries, open)
	if err != nil {
		return "", err
	}

	return hash, nil
}
