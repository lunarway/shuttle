package config

import (
	"io/ioutil"
	"log"
	"os/user"
	"path"
	"path/filepath"
	"regexp"

	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	go_git_ssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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
	LocalPlanPath string
	Configuration ShuttlePlanConfiguration
}

// Load loads a plan from project path and shuttle config
func (p *ShuttlePlanConfiguration) Load(planPath string) *ShuttlePlanConfiguration {
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
func FetchPlan(plan string, projectPath string, localShuttleDirectoryPath string) string {
	switch {
	case isMatching("^git://", plan):
		// We need the user to find the homedir.
		usr, err := user.Current()
		CheckIfError(err)
		planPath := path.Join(localShuttleDirectoryPath, "plan")
		if _, err := os.Stat(planPath); err == nil {
			repo, err := git.PlainOpen(planPath)
			CheckIfError(err)
			worktree, err := repo.Worktree()
			CheckIfError(err)
			err = worktree.Pull(&git.PullOptions{
				RemoteName: "origin",
				Auth:       getSshKeyAuth(usr.HomeDir + "/.ssh/bitbucket_key"),
			})
			if err != nil && err.Error() != "already up-to-date" {
				CheckIfError(err)
			}
		} else {
			_, err := git.PlainClone(planPath, false, &git.CloneOptions{
				URL:      strings.Replace(plan, "git://", "", 1),
				Progress: os.Stdout,
				Auth:     getSshKeyAuth(usr.HomeDir + "/.ssh/bitbucket_key"),
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

func getSshKeyAuth(privateSshKeyFile string) transport.AuthMethod {
	var auth transport.AuthMethod
	sshKey, _ := ioutil.ReadFile(privateSshKeyFile)
	signer, _ := ssh.ParsePrivateKey([]byte(sshKey))
	auth = &go_git_ssh.PublicKeys{User: "git", Signer: signer}
	return auth
}
