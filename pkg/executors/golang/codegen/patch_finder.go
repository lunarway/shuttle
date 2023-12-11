package codegen

import (
	"context"
	"errors"
)

// packageFinder exists to find whatever patches are required for a given shuttle golang action to function
type packageFinder interface {
	// Find should return how many packages are required to function
	Find(ctx context.Context, rootDir string) (packages map[string]string, ok bool, err error)
}

type chainedPackageFinder struct {
	finders []packageFinder
}

func newChainedPatchFinder(finders ...packageFinder) *chainedPackageFinder {
	return &chainedPackageFinder{
		finders: finders,
	}
}

// FindPackages is setup as a chain of responsibility, which means that from most significant it will attempt to find packages
// to be used. However, each finder needs to return how many packages it needs to function, as returning ok means that the finder has exclusive access to the packages
func (p *chainedPackageFinder) findPackages(ctx context.Context, rootDir string) (packages map[string]string, err error) {
	for _, finder := range p.finders {
		packages, ok, err := finder.Find(ctx, rootDir)
		if err != nil {
			return nil, err
		}
		if ok {
			return packages, nil
		}
	}

	return nil, errors.New("failed to find a valid patcher")
}
