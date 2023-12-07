package codegen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type goModuleFinder struct{}

func newGoModuleFinder() *goModuleFinder {
	return &goModuleFinder{}
}

func (s *goModuleFinder) Find(ctx context.Context, rootDir string) (packages map[string]string, ok bool, err error) {
	contents, err := s.getGoModFile(rootDir)
	if err != nil {
		return nil, true, err
	}
	if contents == nil {
		return nil, false, nil
	}

	moduleName, modulePath, err := s.getModuleFromModFile(contents)
	if err != nil {
		return nil, true, fmt.Errorf("failed to parse go.mod in root of project: %w", err)
	}

	packages = make(map[string]string, 0)
	packages[moduleName] = modulePath

	return packages, true, nil
}

func (g *goModuleFinder) getGoModFile(rootDir string) (contents []string, err error) {
	goMod := path.Join(rootDir, "go.mod")
	if _, err := os.Stat(goMod); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	modFile, err := os.ReadFile(path.Join(rootDir, "go.mod"))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(modFile), "\n")

	if len(lines) == 0 {
		return nil, errors.New("go mod is empty")
	}

	return lines, nil
}

func (g *goModuleFinder) getModuleFromModFile(contents []string) (moduleName string, modulePath string, err error) {
	for _, line := range contents {
		modFileLine := strings.TrimSpace(line)
		if strings.HasPrefix(modFileLine, "module") {
			sections := strings.Split(modFileLine, " ")
			if len(sections) < 2 {
				return "", "", fmt.Errorf("invalid module line: %s", modFileLine)
			}

			moduleName := sections[1]

			return moduleName, ".", nil
		}
	}

	return "", "", errors.New("failed to find a valid go.mod file")
}
