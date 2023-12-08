package codegen

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"golang.org/x/exp/slices"
)

type goModPatcher struct {
	writeFileFunc writeFileFunc
}

func newGoModPatcher(writeFileFunc writeFileFunc) *goModPatcher {
	return &goModPatcher{writeFileFunc: writeFileFunc}
}

func (p *goModPatcher) patch(rootDir string, shuttleLocalDir string, packages map[string]string) error {
	actionsModFile, err := p.readActionsMod(shuttleLocalDir)
	if err != nil {
		return err
	}

	for _, module := range modulesFromMap(packages) {
		if !actionsModFile.containsModule(module.name) {
			continue
		}

		actionsModFile.replaceModulePath(rootDir, module)
	}

	return actionsModFile.commit()
}

func (g *goModPatcher) readActionsMod(shuttleLocalDir string) (*actionsModFile, error) {
	path := path.Join(shuttleLocalDir, "tmp/go.mod")

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	return &actionsModFile{
		info:    info,
		content: lines,
		path:    path,

		writeFileFunc: g.writeFileFunc,
	}, nil
}

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
