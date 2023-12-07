package codegen

import "context"

type defaultFinder struct{}

func newDefaultFinder() *defaultFinder {
	return &defaultFinder{}
}

func (s *defaultFinder) Find(ctx context.Context, _ string) (packages map[string]string, ok bool, err error) {
	// We return true, as this should be placed last in the chain
	return make(map[string]string, 0), true, nil
}
