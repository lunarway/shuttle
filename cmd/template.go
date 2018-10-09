package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	tmplFuncs "github.com/lunarway/shuttle/pkg/templates"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

type context struct {
	Vars interface{}
	Args map[string]string
}

var templateOutput string
var templateCmd = &cobra.Command{
	Use:   "template [template]",
	Short: "Execute a template",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uii = uii.SetContext(ui.LevelSilent)

		var templateName = args[0]
		projectContext := getProjectContext()

		namedArgs := map[string]string{}
		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			namedArgs[parts[0]] = parts[1]
		}

		templatePath := resolveFirstPath([]string{
			path.Join(projectContext.ProjectPath, "templates", templateName),
			path.Join(projectContext.ProjectPath, templateName),
			path.Join(projectContext.LocalPlanPath, "templates", templateName),
			path.Join(projectContext.LocalPlanPath, templateName),
		})
		if templatePath == "" {
			panic(fmt.Sprintf("Could not find a template named `%s`", templateName))
		}

		tmpl, err := template.New(templateName).Funcs(tmplFuncs.GetFuncMap()).ParseFiles(templatePath)

		if err != nil {
			panic(err)
		}

		context := context{
			Args: namedArgs,
			Vars: projectContext.Config.Variables,
		}

		if templateOutput == "" {
			err = tmpl.Execute(os.Stdout, context)
			if err != nil {
				panic(err)
			}
		} else {
			// TODO: This is probably not the right place to initialize the temp dir?
			os.MkdirAll(projectContext.TempDirectoryPath, os.ModePerm)

			file, err := os.Create(path.Join(projectContext.TempDirectoryPath, templateOutput))
			if err != nil {
				panic(err)
			}

			err = tmpl.Execute(file, context)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	templateCmd.Flags().StringVarP(&templateOutput, "output", "o", "", "Select filename to output file to in temporary directory")
	rootCmd.AddCommand(templateCmd)
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
