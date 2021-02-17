package sdk

import (
	"fmt"
	"github.com/lunarway/shuttle/pkg/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

type ShuttleContext struct {
	Variables                 config.DynamicYaml `yaml:"vars"` //temporarily include a dynamic representation of the variables here so the go-based plans can use this for templating so we're backwards compatible with the existing templates (for the time being)
	ProjectPath               string
	LocalPlanPath             string
	LocalShuttleDirectoryPath string
	TempDirectoryPath         string
}

func LoadShuttleContext(projectPath, localPlanPath string) (ShuttleContext, error) {

	yamlFile, err := LoadShuttleYaml(projectPath)
	if err != nil {
		return ShuttleContext{}, err
	}
	var result ShuttleContext
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return ShuttleContext{}, fmt.Errorf("Failed to parse shuttle configuration. \n\nMake sure your 'shuttle.yaml' is valid. %w", err)
	}

	result.ProjectPath = projectPath
	result.LocalShuttleDirectoryPath = path.Join(result.ProjectPath, ".shuttle")
	result.TempDirectoryPath = path.Join(result.LocalShuttleDirectoryPath, "temp")
	result.LocalPlanPath = localPlanPath

	return result, nil
}

func LoadShuttleYaml(projectPath string) ([]byte, error) {
	file, err := ioutil.ReadFile(path.Join(projectPath, "shuttle.yaml"))
	if err != nil {
		return nil, fmt.Errorf ("Failed to load shuttle configuration. \n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available. %w", err)
	}
	return file, nil
}

