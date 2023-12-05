package codegen

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

type writeFileFunc = func(name string, contents []byte, permissions fs.FileMode) error

func PatchGoMod(rootDir string, shuttleLocalDir string) error {
	return patchGoMod(
		rootDir,
		shuttleLocalDir,
		os.WriteFile,
	)
}

func patchGoMod(rootDir, shuttleLocalDir string, writeFileFunc writeFileFunc) error {
	packages := make(map[string]string, 0)

	if rootWorkspaceExists(rootDir) {
		modules, err := GetWorkspaceModules(rootDir)
		if err != nil {
			return fmt.Errorf("failed to parse go.mod in root of project: %w", err)
		}

		for _, module := range modules {
			moduleName, modulePath, err := GetWorkspaceModule(rootDir, module)
			if err != nil {
				return err
			}
			packages[moduleName] = modulePath
		}

	} else if rootModExists(rootDir) {
		moduleName, modulePath, err := GetRootModule(rootDir)
		if err != nil {
			return fmt.Errorf("failed to parse go.mod in root of project: %w", err)
		}

		packages[moduleName] = modulePath
	}

	if err := patchPackagesUsed(rootDir, shuttleLocalDir, packages, writeFileFunc); err != nil {
		return err
	}

	return nil

}

func patchPackagesUsed(rootDir string, shuttleLocalDir string, packages map[string]string, writeFileFunc writeFileFunc) error {
	actionsModFilePath := path.Join(shuttleLocalDir, "tmp/go.mod")
	relativeActionsModFilePath := strings.TrimPrefix(path.Join(strings.TrimPrefix(shuttleLocalDir, rootDir), "tmp/go.mod"), "/")

	segmentsToRoot := strings.Count(relativeActionsModFilePath, "/")
	actionsModFileContents, err := os.ReadFile(actionsModFilePath)
	if err != nil {
		return err
	}
	actionsModFilePermissions, err := os.Stat(actionsModFilePath)
	if err != nil {
		return err
	}

	actionsModFile := string(actionsModFileContents)
	actionsModFileLines := strings.Split(actionsModFile, "\n")

	actionsModFileContainsModule := func(moduleName string) bool {
		return strings.Contains(actionsModFile, moduleName)
	}

	for moduleName, modulePath := range packages {
		if !actionsModFileContainsModule(moduleName) {
			continue
		}

		relativeToActionsModulePath := path.Join(strings.Repeat("../", segmentsToRoot), modulePath)

		foundReplace := false
		for i, line := range actionsModFileLines {
			lineTrim := strings.TrimSpace(line)

			if strings.Contains(lineTrim, fmt.Sprintf("replace %s", moduleName)) {
				actionsModFileLines[i] = fmt.Sprintf("replace %s => %s", moduleName, relativeToActionsModulePath)
				foundReplace = true
				break
			}
		}

		if !foundReplace {
			actionsModFileLines = append(
				actionsModFileLines,
				fmt.Sprintf("\nreplace %s => %s", moduleName, relativeToActionsModulePath),
			)

		}

		if err != nil {
			return err
		}
	}

	actionsFileWriter := bytes.NewBufferString(strings.Join(actionsModFileLines, "\n"))
	err = writeFileFunc(actionsModFilePath, actionsFileWriter.Bytes(), actionsModFilePermissions.Mode())
	if err != nil {
		return err
	}

	return nil
}

func GetRootModule(rootDir string) (moduleName string, modulePath string, err error) {
	modFile, err := os.ReadFile(path.Join(rootDir, "go.mod"))
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

func GetWorkspaceModule(rootDir string, absoluteModulePath string) (moduleName string, modulePath string, err error) {
	modFile, err := os.ReadFile(path.Join(rootDir, absoluteModulePath, "go.mod"))
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
			modulePath = strings.TrimPrefix(absoluteModulePath, rootDir)

			return moduleName, modulePath, nil
		}
	}

	return "", "", errors.New("failed to find a valid go.mod file")
}

func GetWorkspaceModules(rootDir string) (modules []string, err error) {
	workFile, err := os.ReadFile(path.Join(rootDir, "go.work"))
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

func rootModExists(rootDir string) bool {
	goMod := path.Join(rootDir, "go.mod")
	if _, err := os.Stat(goMod); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func rootWorkspaceExists(rootDir string) bool {
	goWork := path.Join(rootDir, "go.work")
	if _, err := os.Stat(goWork); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
