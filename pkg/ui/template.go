package ui

import (
	"fmt"
	"io"
	"text/template"

	"github.com/lunarway/shuttle/pkg/templates"
)

func Template(destination io.Writer, name, text string, data interface{}) error {
	t := template.New(name)
	t.Funcs(templates.GetFuncMap())
	t, err := t.Parse(text)
	if err != nil {
		return fmt.Errorf("invalid template: %v", err)
	}
	return t.Execute(destination, data)
}
