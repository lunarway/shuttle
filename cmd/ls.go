package cmd

import (
	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

const lsDefaultTempl = `
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

func newLs(uii *ui.UI, contextProvider contextProvider) *cobra.Command {
	var lsFlagTemplate string

	lsCmd := &cobra.Command{
		Use:          "ls [command]",
		Short:        "List possible commands",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := contextProvider()
			if err != nil {
				return err
			}

			var templ string
			if lsFlagTemplate != "" {
				templ = lsFlagTemplate
			} else {
				templ = lsDefaultTempl
			}
			err = ui.Template(cmd.OutOrStdout(), "ls", templ, templData{
				Scripts: context.Scripts,
				Max:     calculateRightPadForKeys(context.Scripts),
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	lsCmd.Flags().StringVar(&lsFlagTemplate, "template", "", "Template string to use. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].")

	return lsCmd
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
