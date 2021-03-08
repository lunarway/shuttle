package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/template"

	tmplFuncs "github.com/lunarway/shuttle/pkg/templates"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type context struct {
	Vars        interface{}
	Args        map[string]string
	PlanPath    string
	ProjectPath string
}

var templateOutput, leftDelimArg, rightDelimArg, delimsArg string
var ignoreProjectOverrides bool
var templateCmd = &cobra.Command{
	Use:   "template [template]",
	Short: "Execute a template",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var templateName = args[0]
		projectContext := getProjectContext()

		namedArgs := map[string]string{}
		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			namedArgs[parts[0]] = parts[1]
		}

		planPaths := []string{
			path.Join(projectContext.LocalPlanPath, "templates", templateName),
			path.Join(projectContext.LocalPlanPath, templateName),
		}

		projectPaths := []string{
			path.Join(projectContext.ProjectPath, "templates", templateName),
			path.Join(projectContext.ProjectPath, templateName),
		}

		var paths []string
		if ignoreProjectOverrides {
			paths = planPaths
		} else {
			paths = append(projectPaths, planPaths...)
		}

		templatePath := resolveFirstPath(paths)
		if templatePath == "" {
			return fmt.Errorf("template `%s` not found", templateName)
		}

		leftDelim, rightDelim, err := parseDelims(leftDelimArg, rightDelimArg, delimsArg)
		if err != nil {
			return err
		}

		tmpl, err := template.New(templateName).Delims(leftDelim, rightDelim).Funcs(tmplFuncs.GetFuncMap()).ParseFiles(templatePath)
		if err != nil {
			uii.Errorln("Parse template file failed\nFile: %s", templatePath)
			return err
		}

		context := context{
			Args:        namedArgs,
			Vars:        projectContext.Config.Variables,
			PlanPath:    projectContext.LocalPlanPath,
			ProjectPath: projectContext.ProjectPath,
		}
		var output io.Writer
		if templateOutput == "" {
			output = os.Stdout
		} else {
			// TODO: This is probably not the right place to initialize the temp dir?
			os.MkdirAll(projectContext.TempDirectoryPath, os.ModePerm)
			templateOutputPath := path.Join(projectContext.TempDirectoryPath, templateOutput)
			file, err := os.Create(templateOutputPath)
			if err != nil {
				return errors.WithMessagef(err, "create template output file '%s'", templateOutputPath)
			}
			output = file
		}

		err = tmpl.ExecuteTemplate(output, path.Base(templatePath), context)
		if err != nil {
			uii.Errorln("Failed to execute template\nPlan: %s\nProject: %s", context.PlanPath, context.ProjectPath)
			return err
		}
		return nil
	},
}

func init() {
	templateCmd.Flags().StringVarP(&templateOutput, "output", "o", "", "Select filename to output file to in temporary directory")
	templateCmd.Flags().StringVarP(&delimsArg, "delims", "", "", "Select delims for templating. Split by ','. If ',' is in the delims, then use --left-delim and --right-delim instead")
	templateCmd.Flags().StringVarP(&leftDelimArg, "left-delim", "", "", "Select delims for templating. Defaults to '{{'")
	templateCmd.Flags().StringVarP(&rightDelimArg, "right-delim", "", "", "Select delims for templating. Defaults to '}}'")
	templateCmd.Flags().BoolVarP(&ignoreProjectOverrides, "ignore-project-overrides", "", false, "Set flag to ignore template files located in the project folder")
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

func parseDelims(leftDelimArg, rightDelimArg, delimsArg string) (string, string, error) {
	if (leftDelimArg != "" && rightDelimArg == "") || (leftDelimArg == "" && rightDelimArg != "") {
		return "", "", fmt.Errorf("--left-delim and --right-delim should always be used together")
	}
	if delimsArg != "" && (leftDelimArg != "" || rightDelimArg != "") {
		return "", "", fmt.Errorf("either use --left-delim and --right-delim together or use --delims")
	}
	if delimsArg != "" {
		parts := strings.Split(delimsArg, ",")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("--delims should have exactly 2 values split by ',' but value was '%s'", delimsArg)
		}
		return parts[0], parts[1], nil
	}
	if leftDelimArg != "" && rightDelimArg != "" {
		return leftDelimArg, rightDelimArg, nil
	}
	return "{{", "}}", nil
}
