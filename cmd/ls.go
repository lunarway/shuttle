package cmd

import (
	"os"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

const templ = `
{{- $max := .Max -}}
Available Scripts:
{{- range $key, $value := .Scripts}}
  {{rightPad $key $max }} {{upperFirst $value.Description}}
{{- end}}
`

type templData struct {
	Scripts map[string]config.ShuttlePlanScript
	Max     int
}

var lsCmd = &cobra.Command{
	Use:   "ls [command]",
	Short: "List possible commands",
	Run: func(cmd *cobra.Command, args []string) {
		context := getProjectContext()
		err := ui.Template(os.Stdout, "ls", templ, templData{
			Scripts: context.Scripts,
			Max:     calculateRightPadForKeys(context.Scripts),
		})
		context.UI.CheckIfError(err)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func calculateRightPadForKeys(m map[string]config.ShuttlePlanScript) int {
	max := 10
	for k := range m {
		if max < len(k) {
			max = len(k)
		}
	}
	return max + 2
}
