package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/ui"
	"gopkg.in/yaml.v2"
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
func (c *ShuttleProjectContext) Setup(projectPath string, uii ui.UI, clean bool, skipGitPlanPulling bool, planArgument string) *ShuttleProjectContext {
	c.Config.getConf(uii, projectPath)
	c.UI = uii
	c.ProjectPath = projectPath
	c.LocalShuttleDirectoryPath = path.Join(c.ProjectPath, ".shuttle")

	if clean {
		os.RemoveAll(c.LocalShuttleDirectoryPath)
		uii.Infoln("Cleaning %s", c.LocalShuttleDirectoryPath)
	}
	os.MkdirAll(c.LocalShuttleDirectoryPath, os.ModePerm)

	c.TempDirectoryPath = path.Join(c.LocalShuttleDirectoryPath, "temp")
	c.LocalPlanPath = FetchPlan(c.Config.Plan, projectPath, c.LocalShuttleDirectoryPath, uii, skipGitPlanPulling, planArgument)
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
func (c *ShuttleConfig) getConf(uii ui.UI, projectPath string) *ShuttleConfig {
	var configPath = path.Join(projectPath, "shuttle.yaml")

	//log.Printf("configpath: %s", configPath)

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		uii.ExitWithErrorCode(2, "Failed to load shuttle configuration: %s\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available.", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		uii.ExitWithErrorCode(2, "Failed to parse shuttle configuration: %s\n\nMake sure your 'shuttle.yaml' is valid.", err)
	}

	switch c.PlanRaw {
	case false:
		// no plan
	default:
		c.Plan = c.PlanRaw.(string)
	}

	return c
}
