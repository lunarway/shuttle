package config

import (
	"fmt"
	"os"
	"path"

	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/ui"
	"gopkg.in/yaml.v2"
)

// DynamicYaml are any yaml document
type DynamicYaml = map[string]interface{}

// ShuttleConfig describes the actual config for each project
type ShuttleConfig struct {
	Plan      string                       `yaml:"-"`
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
func (c *ShuttleProjectContext) Setup(projectPath string, uii ui.UI, clean bool, skipGitPlanPulling bool, planArgument string) (*ShuttleProjectContext, error) {
	_, err := c.Config.getConf(projectPath)
	if err != nil {
		return nil, err
	}
	c.UI = uii
	c.ProjectPath = projectPath
	c.LocalShuttleDirectoryPath = path.Join(c.ProjectPath, ".shuttle")

	if clean {
		uii.Infoln("Cleaning %s", c.LocalShuttleDirectoryPath)
		err := os.RemoveAll(c.LocalShuttleDirectoryPath)
		if err != nil {
			return nil, fmt.Errorf("remove '%s': %w", c.LocalShuttleDirectoryPath, err)
		}
	}
	err = os.MkdirAll(c.LocalShuttleDirectoryPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("create '%s' directory: %w", c.LocalShuttleDirectoryPath, err)
	}

	c.TempDirectoryPath = path.Join(c.LocalShuttleDirectoryPath, "temp")
	c.LocalPlanPath, err = FetchPlan(c.Config.Plan, projectPath, c.LocalShuttleDirectoryPath, uii, skipGitPlanPulling, planArgument)
	if err != nil {
		return nil, err
	}
	_, err = c.Plan.Load(c.LocalPlanPath)
	if err != nil {
		return nil, err
	}
	c.Scripts = make(map[string]ShuttlePlanScript)
	for scriptName, script := range c.Plan.Scripts {
		c.Scripts[scriptName] = script
	}
	for scriptName, script := range c.Config.Scripts {
		c.Scripts[scriptName] = script
	}
	return c, nil
}

// getConf loads the ShuttleConfig from yaml file in the project path
func (c *ShuttleConfig) getConf(projectPath string) (*ShuttleConfig, error) {
	if projectPath == "" {
		return c, nil
	}

	configPath := path.Join(projectPath, "shuttle.yaml")

	file, err := os.Open(configPath)
	if err != nil {
		return c, errors.NewExitCode(2, "Failed to load shuttle configuration: %s\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available.", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)
	err = decoder.Decode(c)
	if err != nil {
		return c, errors.NewExitCode(2, "Failed to parse shuttle configuration: %s\n\nMake sure your 'shuttle.yaml' is valid.", err)
	}

	switch c.PlanRaw {
	case false:
		// no plan
	default:
		c.Plan = c.PlanRaw.(string)
	}

	return c, nil
}
