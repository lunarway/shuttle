package executors

import (
	"errors"
	"io"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/ui"
)

const scriptHelpTemplate = `
{{- $max := .Max -}}
{{ .Description -}}
{{ if not .Description -}} No description available {{- end }}
{{ if .Args }}
Available arguments:
{{- range $i, $arg := .Args}}
  {{ rightPad (print $arg.Name " " $arg.Required) $max -}} {{- $arg.Description }}
{{- end}}
{{- end}}
`

type scriptHelpTemplateData struct {
	Description string
	Args        []scriptHelpTemplateArg
	Max         int
}

type scriptHelpTemplateArg struct {
	Name        string
	Required    string
	Description string
}

func Help(scripts map[string]config.ShuttlePlanScript, script string, output io.Writer) error {
	s, ok := scripts[script]
	if !ok {
		return errors.New("unrecognized script")
	}
	err := ui.Template(output, "runHelp", scriptHelpTemplate, scriptHelpTemplateData{
		Description: s.Description,
		Args:        templateArgs(s.Args),
		Max:         maxLength(s.Args),
	})
	if err != nil {
		return err
	}
	return nil
}

func maxLength(values []config.ShuttleScriptArgs) int {
	max := 10
	for _, value := range values {
		if max < len(value.Name) {
			max = len(value.Name)
		}
		if value.Required {
			max += len(required(true))
		}
	}
	return max + 2
}

func templateArgs(values []config.ShuttleScriptArgs) []scriptHelpTemplateArg {
	scriptArgs := make([]scriptHelpTemplateArg, len(values))
	for i := range values {
		scriptArgs[i] = scriptHelpTemplateArg{
			Name:        values[i].Name,
			Required:    required(values[i].Required),
			Description: values[i].Description,
		}
	}
	return scriptArgs
}

func required(b bool) string {
	if b {
		return "(required)"
	}
	return ""
}
