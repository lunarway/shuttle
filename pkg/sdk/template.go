package sdk

import (
	"fmt"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/lunarway/shuttle/pkg/templates"
	"github.com/pkg/errors"
)

type TemplateContext struct {
	Vars        interface{}
	Args        map[string]string
	PlanPath    string
	ProjectPath string
}

func resolveFirstPath(paths []string) string {
	for _, templatePath := range paths {
		if fileAvailable(templatePath) {
			return templatePath
		}
	}
	return ""
}

func fileAvailable(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ResolveTemplatePath(project ShuttleContext, templateName string) (string, error) {
	templatePath := resolveFirstPath([]string{
		path.Join(project.ProjectPath, "templates", templateName),
		path.Join(project.ProjectPath, templateName),
		path.Join(project.LocalPlanPath, "templates", templateName),
		path.Join(project.LocalPlanPath, templateName),
	})
	if templatePath == "" {
		return "", fmt.Errorf("template `%s` not found", templateName)
	}
	return templatePath, nil
}

func Generate(templatePath, templateName, outputFilepath string, context TemplateContext, leftDelim, rightDelim string) error {
	file, err := os.Create(outputFilepath)
	if err != nil {
		return errors.WithMessagef(err, "create template output file '%s'", outputFilepath)
	}

	err = renderTemplate(templatePath, templateName, file, context, leftDelim, rightDelim)
	if err != nil {
		return err
	}

	return nil
}

func renderTemplate(templatePath, templateName string, output io.Writer, context TemplateContext, leftDelim, rightDelim string) error {
	tmpl, err := template.New(templateName).
		Delims(leftDelim, rightDelim).
		Funcs(templates.GetFuncMap()).
		ParseFiles(templatePath)

	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(output, path.Base(templatePath), context)
	if err != nil {
		return err
	}
	return nil
}
