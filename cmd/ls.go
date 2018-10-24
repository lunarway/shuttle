package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/lunarway/shuttle/pkg/config"
)

const templ = `
{{- $max := calcrpad .Scripts }}
Available Scripts:
{{- range $key, $value := .Scripts}}
  {{rpad $key $max }} {{upperfirst $value.Description}}
{{- end}}
`
var lsCmd = &cobra.Command{
	Use:   "ls [command]",
	Short: "List possible commands",
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		context := getProjectContext()
		tmpl(os.Stdout, templ, context)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsTemplateFuncs = template.FuncMap{
	"trim": strings.TrimSpace,
	"calcrpad": calcrpad,
	"upperfirst": upperfirst,
	//"trimRightSpace":          trimRightSpace,
	//"trimTrailingWhitespaces": trimRightSpace,
	//"appendIfNotPresent":      appendIfNotPresent,
	"rpad": rpad,
	//"gt":                      Gt,
	//"eq":                      Eq,
}

func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(lsTemplateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

// upperfirst upper cases the first letter in the string
func upperfirst(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(s[0:1]), s[1:])
}

// calcrpad calculates the right padding to use for the scripts
func calcrpad(m map[string]config.ShuttlePlanScript) int {
	max := 10
	for k := range m {
		if max < len(k) {
			max = len(k)
		}
	}
	return max + 2
}