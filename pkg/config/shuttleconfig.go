package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	shuttleerrors "github.com/lunarway/shuttle/pkg/errors"
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
	UI                        *ui.UI
}

// Setup the ShuttleProjectContext for a specific path
func (c *ShuttleProjectContext) Setup(
	projectPath string,
	uii *ui.UI,
	clean bool,
	skipGitPlanPulling bool,
	planArgument string,
	strictConfigLookup bool,
) (*ShuttleProjectContext, error) {
	projectPath, err := c.Config.getConf(projectPath, strictConfigLookup)
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
	c.LocalPlanPath, err = FetchPlan(
		c.Config.Plan,
		projectPath,
		c.LocalShuttleDirectoryPath,
		uii,
		skipGitPlanPulling,
		planArgument,
	)
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
func (c *ShuttleConfig) getConf(projectPath string, strictConfigLookup bool) (string, error) {
	if projectPath == "" {
		return projectPath, nil
	}

	file, err := locateShuttleConfigurationFile(projectPath, strictConfigLookup)
	if err != nil {
		return "", shuttleerrors.NewExitCode(
			2,
			"Failed to load shuttle configuration: %s\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available.",
			err,
		)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)
	err = decoder.Decode(c)
	if err != nil {
		return "", shuttleerrors.NewExitCode(
			2,
			"Failed to parse shuttle configuration: %s\n\nMake sure your 'shuttle.yaml' is valid.",
			err,
		)
	}

	switch c.PlanRaw {
	case false:
		// no plan
	default:
		c.Plan = c.PlanRaw.(string)
	}
	// return the path where the shuttle.yaml file was found
	return path.Dir(file.Name()), nil
}

var errShuttleFileNotFound = errors.New("shuttle.yaml file not found")

func locateShuttleConfigurationFile(startPath string, strictConfigLookup bool) (*os.File, error) {
	var err error
	for {
		configPath := path.Join(startPath, "shuttle.yaml")

		var file *os.File
		file, err = os.Open(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				if startPath == "" || startPath == "/" {
					err = errShuttleFileNotFound
					break
				}

				if strictConfigLookup {
					err = errShuttleFileNotFound
					break
				}

				startPath = removeLastDirectory(startPath)
				continue
			}
			break
		}

		return file, nil
	}

	return nil, err
}

func removeLastDirectory(projectPath string) string {
	parts := strings.Split(projectPath, "/")

	newProjectPath := strings.Join(parts[0:len(parts)-1], "/")

	// when handling the root path / the split and join will produce an empty
	// string so to keep the absolute path set the root path directly.
	if path.IsAbs(projectPath) && newProjectPath == "" {
		return "/"
	}

	return newProjectPath
}
