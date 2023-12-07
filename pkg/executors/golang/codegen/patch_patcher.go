package codegen

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/exp/slices"
)

type module struct {
	name string
	path string
}

func modulesFromMap(packages map[string]string) []module {
	modules := make([]module, 0, len(packages))
	for moduleName, modulePath := range packages {
		modules = append(modules, module{
			name: moduleName,
			path: modulePath,
		})
	}
	slices.SortFunc(modules, func(a, b module) int {
		return strings.Compare(a.name, b.name)
	})

	return modules
}

type goModPatcher struct {
	writeFileFunc writeFileFunc
}

func newGoModPatcher(writeFileFunc writeFileFunc) *goModPatcher {
	return &goModPatcher{writeFileFunc: writeFileFunc}
}

func (p *goModPatcher) patch(rootDir string, shuttleLocalDir string, packages map[string]string) error {
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

	for _, module := range modulesFromMap(packages) {
		if !actionsModFileContainsModule(module.name) {
			continue
		}

		relativeToActionsModulePath := path.Join(strings.Repeat("../", segmentsToRoot), module.path)

		foundReplace := false
		for i, line := range actionsModFileLines {
			lineTrim := strings.TrimSpace(line)

			if strings.Contains(lineTrim, fmt.Sprintf("replace %s", module.name)) {
				actionsModFileLines[i] = fmt.Sprintf("replace %s => %s", module.name, relativeToActionsModulePath)
				foundReplace = true
				break
			}
		}

		if !foundReplace {
			actionsModFileLines = append(
				actionsModFileLines,
				fmt.Sprintf("\nreplace %s => %s", module.name, relativeToActionsModulePath),
			)

		}

		if err != nil {
			return err
		}
	}

	actionsFileWriter := bytes.NewBufferString(strings.Join(actionsModFileLines, "\n"))
	err = p.writeFileFunc(actionsModFilePath, actionsFileWriter.Bytes(), actionsModFilePermissions.Mode())
	if err != nil {
		return err
	}

	return nil
}
