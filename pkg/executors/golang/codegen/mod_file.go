package codegen

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"golang.org/x/exp/slices"
)

type actionsModFile struct {
	info    fs.FileInfo
	content []string
	path    string

	writeFileFunc writeFileFunc
}

func (a *actionsModFile) containsModule(moduleName string) bool {
	return slices.ContainsFunc(a.content, func(s string) bool {
		return strings.Contains(s, moduleName)
	})
}

func (a *actionsModFile) replaceModulePath(rootDir string, module module) {
	relativeToActionsModulePath := path.Join(strings.Repeat("../", a.segmentsTo(rootDir)), module.path)

	foundReplace := false
	for i, line := range a.content {
		lineTrim := strings.TrimSpace(line)

		if strings.Contains(lineTrim, fmt.Sprintf("replace %s", module.name)) {
			a.content[i] = fmt.Sprintf("replace %s => %s", module.name, relativeToActionsModulePath)
			foundReplace = true
			break
		}

	}

	if !foundReplace {
		a.content = append(
			a.content,
			fmt.Sprintf("\nreplace %s => %s", module.name, relativeToActionsModulePath),
		)

	}
}

func (a *actionsModFile) segmentsTo(dirPath string) int {
	relativeActionsModFilePath := strings.TrimPrefix(
		strings.TrimPrefix(
			a.path,
			dirPath,
		),
		"/",
	)

	return strings.Count(relativeActionsModFilePath, "/")
}

func (a *actionsModFile) commit() error {
	return a.writeFileFunc(
		a.path,
		[]byte(strings.Join(a.content, "\n")),
		a.info.Mode(),
	)
}
