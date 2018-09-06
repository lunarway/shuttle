package config

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// DynamicYaml are any yaml document
type DynamicYaml = map[string]interface{}

// ShuttleConfig describes the actual config for each project
type ShuttleConfig struct {
	Plan      string      `yaml:"plan"`
	Variables DynamicYaml `yaml:"vars"`
}

// ShuttleProjectContext describes the context of the project using shuttle
type ShuttleProjectContext struct {
	ProjectPath               string
	LocalShuttleDirectoryPath string
	TempDirectoryPath         string
	Config                    ShuttleConfig
	LocalPlanPath             string
	Plan                      ShuttlePlanConfiguration
}

// Setup the ShuttleProjectContext for a specific path
func (c *ShuttleProjectContext) Setup(projectPath string) *ShuttleProjectContext {
	c.Config.getConf(projectPath)
	c.ProjectPath = projectPath
	c.LocalShuttleDirectoryPath = path.Join(c.ProjectPath, ".shuttle")
	c.TempDirectoryPath = path.Join(c.LocalShuttleDirectoryPath, "temp")
	c.LocalPlanPath = FetchPlan(c.Config.Plan, projectPath, c.LocalShuttleDirectoryPath)
	c.Plan.Load(c.LocalPlanPath)
	return c
}

// getConf loads the ShuttleConfig from yaml file in the project path
func (c *ShuttleConfig) getConf(projectPath string) *ShuttleConfig {
	var configPath = path.Join(projectPath, "shuttle.yaml")

	//log.Printf("configpath: %s", configPath)

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("yamlFile.Get err   #%v ", err))
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal: %v", err))
	}

	return c
}
