package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/ui"
	yaml "gopkg.in/yaml.v2"
)

// DynamicYaml are any yaml document
type DynamicYaml = map[string]interface{}

// ShuttleConfig describes the actual config for each project
type ShuttleConfig struct {
	Plan      string                       `yaml:"not_plan"`
	PlanRaw   interface{}                  `yaml:"plan"`
	Variables DynamicYaml                  `yaml:"vars"`
	Scripts   map[string]ShuttlePlanScript `yaml:"scripts"`
}

// ShuttleProjectContext describes the context of the project using shuttle
type ShuttleProjectContext struct {
	ProjectPath               string
	LocalShuttleDirectoryPath string
	TempDirectoryPath         string
	Config                    ShuttleConfig
	LocalPlanPath             string
	Plan                      ShuttlePlanConfiguration
	Scripts                   map[string]ShuttlePlanScript
	UI                        ui.UI
}

// Setup the ShuttleProjectContext for a specific path
func (c *ShuttleProjectContext) Setup(projectPath string, uii ui.UI, clean bool, skipGitPlanPulling bool) *ShuttleProjectContext {
	c.Config.getConf(projectPath)
	c.UI = uii
	c.ProjectPath = projectPath
	c.LocalShuttleDirectoryPath = path.Join(c.ProjectPath, ".shuttle")

	if clean {
		os.RemoveAll(c.LocalShuttleDirectoryPath)
		uii.InfoLn("Cleaning %s", c.LocalShuttleDirectoryPath)
	}
	os.MkdirAll(c.LocalShuttleDirectoryPath, os.ModePerm)

	c.TempDirectoryPath = path.Join(c.LocalShuttleDirectoryPath, "temp")
	c.LocalPlanPath = FetchPlan(c.Config.Plan, projectPath, c.LocalShuttleDirectoryPath, uii, skipGitPlanPulling)
	c.Plan.Load(c.LocalPlanPath)
	c.Scripts = make(map[string]ShuttlePlanScript)
	for scriptName, script := range c.Plan.Scripts {
		c.Scripts[scriptName] = script
	}
	for scriptName, script := range c.Config.Scripts {
		c.Scripts[scriptName] = script
	}
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

	switch c.PlanRaw {
	case false:
		// no plan
	default:
		c.Plan = c.PlanRaw.(string)
	}

	return c
}
