package ui

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

var lsTemplateFuncs = template.FuncMap{
	"trim":       strings.TrimSpace,
	"upperFirst": upperFirst,
	"rightPad":   rightPad,
}

func Template(destination io.Writer, name, text string, data interface{}) error {
	t := template.New(name)
	t.Funcs(lsTemplateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(destination, data)
}

// rightPad adds padding to the right of a string.
func rightPad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

// upperFirst upper cases the first letter in the string
func upperFirst(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(s[0:1]), s[1:])
}
