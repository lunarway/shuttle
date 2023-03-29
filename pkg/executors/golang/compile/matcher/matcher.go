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

	"github.com/kjuulh/shuttletask/pkg/discover"
	"golang.org/x/mod/sumdb/dirhash"
)

func BinaryMatches(
	ctx context.Context,
	hash string,
	shuttletask *discover.ShuttleTaskDiscovered,
) (string, bool, error) {
	shuttlebindir := path.Join(shuttletask.ParentDir, ".shuttle/shuttletask/binaries")

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

	if binary.Name() == fmt.Sprintf("shuttletask-%s", hex.EncodeToString([]byte(hash)[:16])) {
		return path.Join(shuttlebindir, binary.Name()), true, nil
	} else {
		log.Printf("DEBUG: binary does not match, rebuilding...")
		return "", false, nil
	}
}

func GetHash(ctx context.Context, shuttletask *discover.ShuttleTaskDiscovered) (string, error) {
	entries := make([]string, len(shuttletask.Files))

	for i, task := range shuttletask.Files {
		entries[i] = path.Join(shuttletask.DirPath, task)
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
