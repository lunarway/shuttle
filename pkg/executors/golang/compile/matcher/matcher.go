package matcher

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
	"golang.org/x/mod/sumdb/dirhash"
)

func BinaryMatches(
	ctx context.Context,
	hash string,
	actions *discover.ActionsDiscovered,
) (string, bool, error) {
	shuttlebindir := path.Join(actions.ParentDir, ".shuttle/actions/binaries")

	if _, err := os.Stat(shuttlebindir); errors.Is(err, os.ErrNotExist) {
		log.Println("DEBUG: package doesn't exist continueing")
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

	if binary.Name() == fmt.Sprintf("actions-%s", hex.EncodeToString([]byte(hash)[:16])) {
		return path.Join(shuttlebindir, binary.Name()), true, nil
	} else {
		log.Printf("DEBUG: binary does not match, rebuilding...")
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
