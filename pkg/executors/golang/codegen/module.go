package codegen

import (
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
