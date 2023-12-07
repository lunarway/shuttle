package codegen

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/lunarway/shuttle/pkg/ui"
)

type writeFileFunc = func(name string, contents []byte, permissions fs.FileMode) error

func PatchGoMod(ctx context.Context, rootDir string, shuttleLocalDir string) error {
	return patchGoMod(
		ctx,
		rootDir,
		shuttleLocalDir,
		ui,
		os.WriteFile,
	)
}

func patchGoMod(ctx context.Context, rootDir, shuttleLocalDir string, writeFileFunc writeFileFunc) error {
	packages, err := newChainedwPatchFinder(
		newWorkspaceFinder(rootDir),
		newGoModuleFinder(rootDir),
		newDefaultFinder(),
	).findPackages(ctx)
	if err != nil {
		return err
	}

	if err := newGoModPatcher(writeFileFunc).patchPackagesUsed(rootDir, shuttleLocalDir, packages); err != nil {
		return err
	}

	return nil

}

// packageFinder exists to find whatever patches are required for a given shuttle golang action to function
type packageFinder interface {
	// Find should return how many packages are required to function
	Find(ctx context.Context) (packages map[string]string, ok bool, err error)
}

type goModPatcher struct {
	writeFileFunc writeFileFunc
}

func newGoModPatcher(writeFileFunc writeFileFunc) *goModPatcher {
	return &goModPatcher{writeFileFunc: writeFileFunc}
}

func (p *goModPatcher) patchPackagesUsed(rootDir string, shuttleLocalDir string, packages map[string]string) error {
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
	err = p.writeFileFunc(actionsModFilePath, actionsFileWriter.Bytes(), actionsModFilePermissions.Mode())
	if err != nil {
		return err
	}

	return nil
}
