package codegen

import (
	"context"
	"embed"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/kjuulh/shuttletask/pkg/discover"
	"github.com/kjuulh/shuttletask/pkg/parser"
)

var (
	//go:embed templates/mainFile.tmpl
	mainFileTmpl embed.FS
)

func GenerateMainFile(
	ctx context.Context,
	shuttlelocaldir string,
	shuttletask *discover.ShuttleTaskDiscovered,
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
