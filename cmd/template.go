package cmd

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	tmplFuncs "bitbucket.org/LunarWay/shuttle/pkg/templates"
	"github.com/spf13/cobra"
)

type context struct {
	Vars interface{}
	Args map[string]string
}

var templateCmd = &cobra.Command{
	Use:   "template [template]",
	Short: "Execute a template",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var templateName = args[0]
		projectContext := getProjectContext()

		namedArgs := map[string]string{}
		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			namedArgs[parts[0]] = parts[1]
		}
		templateDir := path.Join(projectContext.LocalPlanPath, "templates")

		templatePath := path.Join(templateDir, templateName)

		tmpl, err := template.New(templateName).Funcs(tmplFuncs.GetFuncMap()).ParseFiles(templatePath)

		if err != nil {
			panic(err)
		}
		// TODO: This is probably not the right place to initialize the temp dir?
		os.Mkdir(projectContext.TempDirectoryPath, 0755)

		outputYamlFilename := strings.Replace(templateName, "tmpl", "yaml", -1)

		fmt.Printf("Template generated: .shuttle/temp/%s\n", outputYamlFilename)
		file, err := os.Create(path.Join(projectContext.TempDirectoryPath, outputYamlFilename))

		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(file, context{
			Args: namedArgs,
			Vars: projectContext.Config.Variables,
		})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
