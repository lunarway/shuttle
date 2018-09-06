package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls [command]",
	Short: "List possible commands",
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		loadedPlan := getPlan()
		tmpl(os.Stdout, `Available Scripts:{{range $key, $value := .Scripts}}
  {{rpad $key 10 }} {{$value.Description}}{{end}}
`, loadedPlan.Configuration)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsTemplateFuncs = template.FuncMap{
	"trim": strings.TrimSpace,
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
