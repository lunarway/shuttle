package codegen

import (
	"context"
	"embed"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/lunarway/shuttle/pkg/executors/golang/discover"
	"github.com/lunarway/shuttle/pkg/executors/golang/parser"
)

//go:embed templates/mainFile.tmpl
var mainFileTmpl embed.FS

func GenerateMainFile(
	ctx context.Context,
	shuttlelocaldir string,
	actions *discover.ActionsDiscovered,
	functions []*parser.Function,
) error {
	tmpmainfile := path.Join(shuttlelocaldir, "tmp/main.go")

	file, err := os.Create(tmpmainfile)
	if err != nil {
		return err
	}

	tmpl := template.
		Must(
			template.
				New("mainFile.tmpl").
				Funcs(map[string]any{
					"lower": strings.ToLower,
				}).
				ParseFS(mainFileTmpl, "templates/mainFile.tmpl"),
		)

	err = tmpl.Execute(file, functions)

	return err
}
