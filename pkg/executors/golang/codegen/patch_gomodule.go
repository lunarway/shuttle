package codegen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type goModuleFinder struct {
	rootDir string
}

func newGoModuleFinder(rootDir string) *goModuleFinder {
	return &goModuleFinder{
		rootDir: rootDir,
	}
}

func (s *goModuleFinder) Find(ctx context.Context) (packages map[string]string, ok bool, err error) {
	if !s.rootModExists() {
		return nil, false, nil
	}

	moduleName, modulePath, err := s.getRootModule()
	if err != nil {
		return nil, true, fmt.Errorf("failed to parse go.mod in root of project: %w", err)
	}

	packages = make(map[string]string, 0)
	packages[moduleName] = modulePath

	return packages, true, nil
}

func (g *goModuleFinder) getRootModule() (moduleName string, modulePath string, err error) {
	modFile, err := os.ReadFile(path.Join(g.rootDir, "go.mod"))
	if err != nil {
		return "", "", err
	}

	modFileContent := string(modFile)
	lines := strings.Split(modFileContent, "\n")
	if len(lines) == 0 {
		return "", "", errors.New("go mod is empty")
	}

	for _, line := range lines {
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

func (g *goModuleFinder) rootModExists() bool {
	goMod := path.Join(g.rootDir, "go.mod")
	if _, err := os.Stat(goMod); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
