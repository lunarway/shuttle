package config

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"

	"fmt"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
)

// ShuttlePlanScript is a ShuttlePlan sub-element
type ShuttlePlanScript struct {
	Description string              `yaml:"description"`
	Actions     []ShuttleAction     `yaml:"actions"`
	Args        []ShuttleScriptArgs `yaml:"args"`
}

type ShuttleScriptArgs struct {
	Name     string `yaml:"name"`
	Required bool   `yaml:"required"`
}

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
	ShuttleConfig ShuttleConfig
	PlanPath      string
	Dockerfile    string
	Configuration ShuttlePlanConfiguration
}

// Load loads a plan from project path and shuttle config
func (p *ShuttlePlan) Load(projectPath string, shuttleConfig ShuttleConfig) *ShuttlePlan {
	p.PlanPath = fetchPlan(shuttleConfig.Plan, projectPath)
	p.ProjectPath = projectPath
	p.ShuttleConfig = shuttleConfig
	p.Dockerfile = path.Join(p.PlanPath, "Dockerfile")
	p.Configuration.getPlanConf(p.PlanPath)
	return p
}

func (c *ShuttlePlanConfiguration) getPlanConf(planPath string) *ShuttlePlanConfiguration {
	var configPath = path.Join(planPath, "plan.yaml")
	//log.Printf("configpath: %s", configPath)
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
		//panic("plan not valid: git is not supported yet")
		planPath := path.Join(projectPath, ".shuttle/plan")

		if _, err := os.Stat(planPath); err == nil {
			repo, err := git.PlainOpen(planPath)
			CheckIfError(err)
			worktree, err := repo.Worktree()
			CheckIfError(err)
			err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil && err.Error() != "already up-to-date" {
				CheckIfError(err)
			}
		} else {
			_, err := git.PlainClone(planPath, false, &git.CloneOptions{
				URL:      strings.Replace(plan, "git://", "", 1),
				Progress: os.Stdout,
			})
			if err != nil {
				panic(err)
			}
		}
		return planPath
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

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
