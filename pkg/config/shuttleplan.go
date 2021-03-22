package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/git"
	"github.com/lunarway/shuttle/pkg/ui"
	"gopkg.in/yaml.v2"
)

// ShuttlePlanScript is a ShuttlePlan sub-element
type ShuttlePlanScript struct {
	Description string              `yaml:"description"`
	Actions     []ShuttleAction     `yaml:"actions"`
	Args        []ShuttleScriptArgs `yaml:"args"`
}

// ShuttleScriptArgs describes an arguments that a script accepts
type ShuttleScriptArgs struct {
	Name        string `yaml:"name"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
}

func (a ShuttleScriptArgs) String() string {
	var s strings.Builder
	s.WriteString(a.Name)
	if a.Required {
		s.WriteString(" (required)")
	}
	if len(a.Description) != 0 {
		fmt.Fprintf(&s, "  %s", a.Description)
	}
	return s.String()
}

// ShuttleAction describes an action done by a shuttle script
type ShuttleAction struct {
	Shell      string `yaml:"shell"`
	Dockerfile string `yaml:"dockerfile"`
}

// ShuttlePlanConfiguration is a ShuttlePlan sub-element
type ShuttlePlanConfiguration struct {
	Scripts map[string]ShuttlePlanScript `yaml:"scripts"`
}

// ShuttlePlan struct describes a plan
type ShuttlePlan struct {
	ProjectPath   string
	LocalPlanPath string
	Configuration ShuttlePlanConfiguration
}

// Load loads a plan from project path and shuttle config
func (p *ShuttlePlanConfiguration) Load(planPath string) (*ShuttlePlanConfiguration, error) {
	if planPath == "" {
		return p, nil
	}

	configPath := path.Join(planPath, "plan.yaml")

	file, err := os.Open(configPath)
	if err != nil {
		return p, errors.NewExitCode(2, "Failed to open plan configuration: %s\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available.", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)
	err = decoder.Decode(p)
	if err != nil {
		return p, errors.NewExitCode(1, "Failed to load plan configuration: %s\n\nThis is likely an issue with the referenced plan. Please, contact the plan maintainers.", err)
	}

	return p, nil
}

// FetchPlan so it exists locally and return path to that plan
func FetchPlan(plan string, projectPath string, localShuttleDirectoryPath string, uii ui.UI, skipGitPlanPulling bool, planArgument string) (string, error) {
	if isPlanArgumentAPlan(planArgument) {
		uii.Infoln("Using overloaded plan %v", planArgument)
		return FetchPlan(getPlanFromPlanArgument(planArgument), projectPath, localShuttleDirectoryPath, uii, skipGitPlanPulling, "")
	}

	switch {
	case plan == "":
		uii.Verboseln("Using no plan")
		return "", nil
	case git.IsGitPlan(plan):
		uii.Verboseln("Using git plan at '%s'", plan)
		return git.GetGitPlan(plan, localShuttleDirectoryPath, uii, skipGitPlanPulling, planArgument)
	case isMatching("^http://|^https://", plan):
		panic(fmt.Sprintf("Plan '%v' is not valid: non-git http/https is not supported yet", plan))
	case isFilePath(plan, true):
		uii.Verboseln("Using local plan at '%s'", plan)
		return plan, nil
	case isFilePath(plan, false):
		uii.Verboseln("Using local plan at '%s'", plan)
		return path.Join(projectPath, plan), nil
	default:
		return "", errors.NewExitCode(2, "Unknown plan path '%s'", plan)
	}
}

func isFilePath(path string, matchOnlyAbs bool) bool {
	return filepath.IsAbs(path) == matchOnlyAbs
}

func isMatching(r string, content string) bool {
	match, err := regexp.MatchString(r, content)
	if err != nil {
		panic(err)
	}
	return match
}
