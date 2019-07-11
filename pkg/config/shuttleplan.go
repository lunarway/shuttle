package config

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"fmt"
	"os"

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
func (p *ShuttlePlanConfiguration) Load(planPath string) *ShuttlePlanConfiguration {
	if planPath == "" {
		return p
	}
	var configPath = path.Join(planPath, "plan.yaml")
	//log.Printf("configpath: %s", configPath)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return p
}

// FetchPlan so it exists locally and return path to that plan
func FetchPlan(plan string, projectPath string, localShuttleDirectoryPath string, uii ui.UI, skipGitPlanPulling bool, planArgument string) string {
	if isPlanArgumentAPlan(planArgument) {
		uii.Infoln("Using overloaded plan %v", planArgument)
		return FetchPlan(getPlanFromPlanArgument(planArgument), projectPath, localShuttleDirectoryPath, uii, skipGitPlanPulling, "")
	}
	
	switch {
	case plan == "":
		uii.Verboseln("Using no plan")
		return ""
	case git.IsGitPlan(plan):
		uii.Verboseln("Using git plan at '%s'", plan)
		return git.GetGitPlan(plan, localShuttleDirectoryPath, uii, skipGitPlanPulling, planArgument)
	case isMatching("^http://|^https://", plan):
		panic(fmt.Sprintf("Plan '%v' is not valid: non-git http/https is not supported yet", plan))
	case isFilePath(plan, true):
		uii.Verboseln("Using local plan at '%s'", plan)
		return plan
	case isFilePath(plan, false):
		uii.Verboseln("Using local plan at '%s'", plan)
		return path.Join(projectPath, plan)
	}
	panic("Unknown plan path '" + plan + "'")
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

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
