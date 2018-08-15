package config

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

// ShuttlePlanScript is a ShuttlePlan sub-element
type ShuttlePlanScript struct {
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
	PlanPath      string
	Dockerfile    string
	Configuration ShuttlePlanConfiguration
}

// Load loads a plan from project path and shuttle config
func (p *ShuttlePlan) Load(projectPath string, shuttleConfig ShuttleConfig) *ShuttlePlan {
	log.Printf(shuttleConfig.Plan)
	p.PlanPath = fetchPlan(shuttleConfig.Plan, projectPath)
	p.ProjectPath = projectPath
	p.Dockerfile = path.Join(p.PlanPath, "Dockerfile")
	p.Configuration.getPlanConf(p.PlanPath)

	return p
}

func (c *ShuttlePlanConfiguration) getPlanConf(planPath string) *ShuttlePlanConfiguration {
	var configPath = path.Join(planPath, "plan.yaml")
	log.Printf("configpath: %s", configPath)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func fetchPlan(plan string, projectPath string) string {
	switch {
	case isMatching("^git://", plan):
		panic("plan not valid: git is not supported yet")
	// 	_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
	// 		URL:      "https://github.com/src-d/go-git",
	// 		Progress: os.Stdout,
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return "/tmp/foo"
	case isMatching("^http://|^https://", plan):
		panic("plan not valid: http is not supported yet")
	case isFilePath(plan, true):
		return plan
	case isFilePath(plan, false):
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
