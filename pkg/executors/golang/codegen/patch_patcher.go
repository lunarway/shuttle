package codegen

import (
	"os"
	"path"
	"strings"
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
