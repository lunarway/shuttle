package discover

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/kjuulh/shuttletask/pkg/shuttlefile"
)

var (
	InvalidShuttlePathFile = errors.New("shuttle path did not point ot a shuttle.yaml file")
)

const (
	shuttletaskdir  = "shuttletask"
	shuttlefilename = "shuttle.yaml"
)

type ShuttleTaskDiscovered struct {
	Files     []string
	DirPath   string
	ParentDir string
}

type Discovered struct {
	Local *ShuttleTaskDiscovered
	Plan  *ShuttleTaskDiscovered
}

// path: is a path to the shuttle.yaml file
// It will always look for the shuttletask directory relative to the shuttle.yaml file
//
// 1. Traverse shuttletaskdir
//
// 2. Traverse plan if exists (only 1 layer for now)
//
// 3. Collect file names
//
// 4. Return list of files to move to tmp dir
func Discover(ctx context.Context, shuttlepath string) (*Discovered, error) {
	if !strings.HasSuffix(shuttlepath, shuttlefilename) {
		return nil, InvalidShuttlePathFile
	}
	if _, err := os.Stat(shuttlepath); errors.Is(err, os.ErrNotExist) {
		return nil, InvalidShuttlePathFile
	}

	localdir := path.Dir(shuttlepath)
	localPlan, err := discoverPlan(localdir)
	if err != nil {
		return nil, err
	}

	config, err := shuttlefile.ParseFile(ctx, shuttlepath)
	if err != nil {
		return nil, err
	}

	discovered := Discovered{
		Local: localPlan,
	}

	if config.Plan != "" {
		planShuttleFile := path.Join(localdir, ".shuttle/plan")
		parentPlan, err := discoverPlan(planShuttleFile)
		if err != nil {
			return nil, err
		}

		discovered.Plan = parentPlan
	}

	return &discovered, nil
}

func discoverPlan(localdir string) (*ShuttleTaskDiscovered, error) {
	localshuttledirentries := make([]string, 0)

	shuttletaskpath := path.Join(localdir, shuttletaskdir)
	if fs, err := os.Stat(shuttletaskpath); err == nil {
		// list all local files
		if fs.IsDir() {
			entries, err := os.ReadDir(shuttletaskpath)
			if err != nil {
				return nil, err
			}

			for _, entry := range entries {
				// skip dirs
				if entry.IsDir() {
					continue
				}

				// skip non go files
				if !strings.HasSuffix(entry.Name(), ".go") {
					continue
				}

				// skip test files
				if strings.HasSuffix(entry.Name(), "test.go") {
					continue
				}

				localshuttledirentries = append(localshuttledirentries, entry.Name())
			}
		}

		return &ShuttleTaskDiscovered{
			DirPath:   shuttletaskpath,
			Files:     localshuttledirentries,
			ParentDir: localdir,
		}, nil

	}
	return nil, nil

}
