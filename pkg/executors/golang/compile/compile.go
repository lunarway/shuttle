package compile

import (
	"context"
	"path"

	"github.com/lunarway/shuttle/pkg/executors/golang/codegen"
	"github.com/lunarway/shuttle/pkg/executors/golang/compile/matcher"
	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
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

	if err = codegen.Format(ctx, shuttlelocaldir); err != nil {
		return "", err
	}

	if err = codegen.ModTidy(ctx, shuttlelocaldir); err != nil {
		return "", err
	}
	binarypath, err := codegen.CompileBinary(ctx, ui, shuttlelocaldir)
	if err != nil {
		return "", err
	}

	finalBinaryPath := shuttlefolder.CalculateBinaryPath(shuttlelocaldir, hash)
	shuttlefolder.Move(binarypath, finalBinaryPath)

	return finalBinaryPath, nil
}
