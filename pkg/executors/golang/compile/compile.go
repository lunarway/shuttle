package compile

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"dagger.io/dagger"
	"github.com/lunarway/shuttle/pkg/executors/golang/codegen"
	"github.com/lunarway/shuttle/pkg/executors/golang/compile/matcher"
	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
	golangerrors "github.com/lunarway/shuttle/pkg/executors/golang/errors"
	"github.com/lunarway/shuttle/pkg/executors/golang/parser"
	"github.com/lunarway/shuttle/pkg/executors/golang/shuttlefolder"
	"github.com/lunarway/shuttle/pkg/ui"
	"golang.org/x/sync/errgroup"
)

const (
	alwaysBuild = false
)

type Binary struct {
	Path string
}

type Binaries struct {
	Local Binary
	Plan  Binary
}

// discovered: Discovered actions projects
//
// 1. Check hash for each dir
//
// 2. Compile for each discovered dir
//
// 2.1. Copy to tmp dir
//
// 2.2. Generate main file
//
// 3. Move binary to .shuttle/actions/binary-<hash>
func Compile(ctx context.Context, ui *ui.UI, discovered *discover.Discovered) (*Binaries, error) {
	egrp, ctx := errgroup.WithContext(ctx)
	binaries := &Binaries{}
	if discovered.Local != nil {
		egrp.Go(func() error {
			ui.Verboseln("compiling golang actions binary for: %s", discovered.Local.DirPath)

			path, err := compile(ctx, ui, discovered.Local)
			if err != nil {
				return err
			}

			binaries.Local = Binary{Path: path}
			return nil
		})
	}
	if discovered.Plan != nil {
		egrp.Go(func() error {
			ui.Verboseln("compiling golang actions binary for: %s", discovered.Plan.DirPath)

			path, err := compile(ctx, ui, discovered.Plan)
			if err != nil {
				return err
			}

			binaries.Plan = Binary{Path: path}
			return nil
		})
	}

	if err := egrp.Wait(); err != nil {
		return nil, err
	}

	return binaries, nil
}

func compile(ctx context.Context, ui *ui.UI, actions *discover.ActionsDiscovered) (string, error) {
	hash, err := matcher.GetHash(ctx, actions)
	if err != nil {
		return "", err
	}

	binaryPath, ok, err := matcher.BinaryMatches(ctx, ui, hash, actions)
	if err != nil {
		return "", err
	}

	if ok && !alwaysBuild {
		ui.Verboseln("file already matches continueing")
		// The binary is the same so we short circuit
		return binaryPath, nil
	}

	shuttlelocaldir := path.Join(actions.ParentDir, ".shuttle/actions")

	if err = shuttlefolder.GenerateTmpDir(ctx, shuttlelocaldir); err != nil {
		return "", err
	}
	if err = shuttlefolder.CopyFiles(ctx, shuttlelocaldir, actions); err != nil {
		return "", err
	}

	contents, err := parser.GenerateAst(ctx, shuttlelocaldir, actions)
	if err != nil {
		return "", err
	}

	if err = codegen.GenerateMainFile(ctx, shuttlelocaldir, actions, contents); err != nil {
		return "", err
	}

	var binarypath string

	if err := codegen.NewPatcher().Patch(ctx, actions.ParentDir, shuttlelocaldir); err != nil {
		return "", fmt.Errorf("failed to patch generated go.mod: %w", err)
	}

	if goInstalled() {
		if err = codegen.ModTidy(ctx, ui, shuttlelocaldir); err != nil {
			return "", fmt.Errorf("go mod tidy failed: %w", err)
		}

		if err = codegen.Format(ctx, ui, shuttlelocaldir); err != nil {
			return "", fmt.Errorf("go fmt failed: %w", err)
		}

		binarypath, err = codegen.CompileBinary(ctx, ui, shuttlelocaldir)
		if err != nil {
			return "", fmt.Errorf("go build failed: %w", err)
		}
	} else if goDaggerFallback() {
		binarypath, err = compileWithDagger(ctx, ui, shuttlelocaldir)
		if err != nil {
			return "", fmt.Errorf("failed to compile with dagger: %w", err)
		}
	} else {
		return "", golangerrors.ErrGolangActionNoBuilder
	}

	finalBinaryPath := shuttlefolder.CalculateBinaryPath(shuttlelocaldir, hash)
	if err := shuttlefolder.Move(binarypath, finalBinaryPath); err != nil {
		return "", fmt.Errorf("failed to remove actions binary to final destination: %w", err)
	}

	return finalBinaryPath, nil
}

func compileWithDagger(ctx context.Context, ui *ui.UI, shuttlelocaldir string) (string, error) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return "", fmt.Errorf("failed to start dagger: %w", err)
	}

	src := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{
			".git/",
			".node_modules/",
			"target/",
		},
		Include: []string{},
	})

	log.Printf("shuttlelocaldir: %s", shuttlelocaldir)

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("dagger failed to get your current dir: %w", err)
	}

	nakedShuttleDir := strings.TrimPrefix(strings.TrimPrefix(shuttlelocaldir, dir), "/")
	log.Printf("nakedShuttleDir: %s", nakedShuttleDir)

	shuttleBinary := client.Container().
		From(getGolangImage()).
		WithWorkdir("/app").
		WithDirectory(".", src).
		WithWorkdir(path.Join(nakedShuttleDir, "tmp")).
		WithExec([]string{
			"go", "mod", "tidy",
		}).
		WithExec([]string{
			"go", "fmt", "./...",
		}).
		WithExec([]string{
			"go",
			"build",
			// TODO: add cross compilation
		})

	_, err = shuttleBinary.Sync(ctx)
	if err != nil {
		return "", fmt.Errorf("dagger failed to build binary, see shuttle ls -v to see error output: %w", err)
	}

	shuttleActionsDirectory := shuttleBinary.File("actions")
	exported, err := shuttleActionsDirectory.Export(ctx, path.Join(shuttlelocaldir, "tmp", "actions"))
	if err != nil {
		return "", fmt.Errorf("could not export dagger shuttle actions binary, err: %w", err)
	}
	if !exported {
		return "", fmt.Errorf("failed to export binary")
	}

	return path.Join(shuttlelocaldir, "tmp", "actions"), nil
}

func goInstalled() bool {
	gopath, err := exec.LookPath("go")
	if err != nil {
		return false
	}

	if gopath == "" {
		return false
	}

	return true
}

func getGolangImage() string {
	const (
		// renovate: datasource=docker depName=golang
		golangImageVersion = "1.22.5-alpine"
	)

	golangImage := fmt.Sprintf("golang:%s", golangImageVersion)
	golangImageOverride := os.Getenv("SHUTTLE_GOLANG_ACTIONS_IMAGE")
	if golangImageOverride != "" {
		return golangImageOverride
	}

	return golangImage
}

func goDaggerFallback() bool {
	daggerFallback := os.Getenv("SHUTTLE_GOLANG_ACTIONS_DAGGER_FALLBACK")

	daggerFallbackEnabled, err := strconv.ParseBool(daggerFallback)
	if err != nil {
		return false
	}

	return daggerFallbackEnabled
}
