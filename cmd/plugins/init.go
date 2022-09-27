package plugins

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/lunarway/shuttle/cmd/utility"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

//go:embed templates/plugin.yaml
var pluginFile embed.FS

//go:embed templates/schema.json
var schemaFile embed.FS

func newInitCmd(uii *ui.UI, contextProvider utility.ContextProvider) *cobra.Command {
	var pluginName string
	var pluginPath string
	cmd := &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.ParseFlags(args); err != nil {
				return err
			}

			_, err := contextProvider()
			if err != nil {
				return err
			}

			if _, err := os.Stat(pluginPath); errors.Is(err, os.ErrNotExist) {
				err = os.MkdirAll(pluginPath, 0777)
				if err != nil {
					return fmt.Errorf("could not create plugin dir: %w", err)
				}
			}

			pluginTmpl, err := template.ParseFS(pluginFile, "templates/plugin.yaml")
			if err != nil {
				return err
			}
			pluginF, err := os.Create(path.Join(pluginPath, "plugin.yaml"))
			if err != nil {
				return err
			}

			err = pluginTmpl.Execute(pluginF, struct{ Name string }{Name: pluginName})
			if err != nil {
				return err
			}

			schemaTmpl, err := template.ParseFS(schemaFile, "templates/schema.json")
			if err != nil {
				return err
			}
			schemaF, err := os.Create(path.Join(pluginPath, "schema.json"))
			if err != nil {
				return err
			}

			err = schemaTmpl.Execute(schemaF, struct{ Name string }{Name: pluginName})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&pluginName, "name", "n", "", "the name of the plugin to template")
	cmd.MarkPersistentFlagRequired("name")

	cmd.PersistentFlags().StringVar(&pluginPath, "path", ".", "the place of the plugin to template")

	return cmd
}
