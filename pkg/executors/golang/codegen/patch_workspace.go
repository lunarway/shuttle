package codegen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type workspaceFinder struct {
	rootDir string
}

func newWorkspaceFinder(rootDir string) *workspaceFinder {
	return &workspaceFinder{
		rootDir: rootDir,
	}
}

func (w *workspaceFinder) rootWorkspaceExists() bool {
	goWork := path.Join(w.rootDir, "go.work")
	if _, err := os.Stat(goWork); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (s *workspaceFinder) Find(ctx context.Context) (packages map[string]string, ok bool, err error) {
	if !s.rootWorkspaceExists() {
		return nil, false, nil
	}

	modules, err := s.getWorkspaceModules()
	if err != nil {
		return nil, true, err
	}

	packages = make(map[string]string, 0)
	for _, module := range modules {
		moduleName, modulePath, err := s.getWorkspaceModule(module)
		if err != nil {
			return nil, true, err
		}
		packages[moduleName] = modulePath
	}

	return packages, true, nil
}

func (w *workspaceFinder) getWorkspaceModules() (modules []string, err error) {
	workFile, err := os.ReadFile(path.Join(w.rootDir, "go.work"))
	if err != nil {
		return nil, err
	}

	workFileContent := string(workFile)
	lines := strings.Split(workFileContent, "\n")
	if len(lines) == 0 {
		return nil, errors.New("go work is empty")
	}

	modules = make([]string, 0)
	for _, line := range lines {
		modFileLine := strings.Trim(strings.TrimSpace(line), "\t")
		if strings.HasPrefix(modFileLine, ".") && modFileLine != "./actions" {
			modules = append(
				modules,
				strings.TrimPrefix(
					strings.TrimPrefix(modFileLine, "."),
					"/",
				),
			)
		}
	}

	return modules, nil
}

func (w *workspaceFinder) getWorkspaceModule(absoluteModulePath string) (moduleName string, modulePath string, err error) {
	modFile, err := os.ReadFile(path.Join(w.rootDir, absoluteModulePath, "go.mod"))
	if err != nil {
		return "", "", fmt.Errorf("failed to find go.mod at: %s: %w", absoluteModulePath, err)
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
			modulePath = strings.TrimPrefix(absoluteModulePath, w.rootDir)

			return moduleName, modulePath, nil
		}
	}

	return "", "", errors.New("failed to find a valid go.mod file")
}
