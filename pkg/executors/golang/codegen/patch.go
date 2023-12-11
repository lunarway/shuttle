package codegen

import (
	"context"
	"io/fs"
	"os"
)

type writeFileFunc = func(name string, contents []byte, permissions fs.FileMode) error

type Patcher struct {
	patchFinder *chainedPackageFinder
	patcher     *goModPatcher
}

func NewPatcher() *Patcher {
	return &Patcher{
		patchFinder: newChainedPatchFinder(
			newWorkspaceFinder(),
			newGoModuleFinder(),
			newDefaultFinder(),
		),
		patcher: newGoModPatcher(os.WriteFile),
	}
}

func (p *Patcher) Patch(ctx context.Context, rootDir string, shuttleLocalDir string) error {
	packages, err := p.patchFinder.findPackages(ctx, rootDir)
	if err != nil {
		return err
	}

	if err := p.patcher.patch(rootDir, shuttleLocalDir, packages); err != nil {
		return err
	}

	return nil
}
