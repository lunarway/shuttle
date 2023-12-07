package codegen

import (
	"context"
	"errors"
)

type chainedPackageFinder struct {
	finders []packageFinder
}

func newChainedwPatchFinder(finders ...packageFinder) *chainedPackageFinder {
	return &chainedPackageFinder{
		finders: finders,
	}
}

// FindPackages is setup as a chain of responsibility, which means that from most significant it will attempt to find packages
// to be used. However, each finder needs to return how many packages it needs to function, as returning ok means that the finder has exclusive access to the packages
func (p *chainedPackageFinder) findPackages(ctx context.Context) (packages map[string]string, err error) {
	for _, finder := range p.finders {
		packages, ok, err := finder.Find(ctx)
		if err != nil {
			return nil, err
		}
		if ok {
			return packages, nil
		}
	}

	return nil, errors.New("failed to find a valid patcher")
}
