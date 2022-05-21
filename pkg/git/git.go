package git

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	go_cmd "github.com/go-cmd/cmd"
	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/ui"
)

type Plan struct {
	IsGitPlan  bool
	Protocol   string
	User       string
	Host       string
	Repository string
	Head       string
}

var gitRegex = regexp.MustCompile(`^((git://((?P<user>[^@]+)@)?(?P<repository1>(?P<host>[^:]+):(?P<path>[^#]*)))|((?P<protocol>https)://(?P<repository2>.*\.git)))(#(?P<head>.*))?$`)

func ParsePlan(plan string) Plan {
	if !gitRegex.MatchString(plan) {
		return Plan{
			IsGitPlan: false,
		}
	}

	match := gitRegex.FindStringSubmatch(plan)
	result := make(map[string]string)
	for i, name := range gitRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	protocol := result["protocol"]
	if protocol == "" {
		protocol = "ssh"
	}
	repository := result["repository1"]
	if repository == "" {
		repository = result["repository2"]
	}
	head := result["head"]
	if head == "" {
		head = "master"
	}

	return Plan{
		IsGitPlan:  true,
		Protocol:   protocol,
		User:       result["user"],
		Host:       result["host"],
		Repository: repository,
		Head:       head,
	}
}

// IsPlan returns true if specified plan is a git plan
func IsPlan(plan string) bool {
	parsedGitPlan := ParsePlan(plan)
	return parsedGitPlan.IsGitPlan
}

// GetGitPlan will pull git repository and return its path
func GetGitPlan(plan string, localShuttleDirectoryPath string, uii *ui.UI, skipGitPlanPulling bool, planArgument string) (string, error) {
	parsedGitPlan := ParsePlan(plan)

	if strings.HasPrefix(planArgument, "#") {
		parsedGitPlan.Head = planArgument[1:]
		uii.EmphasizeInfoln("Overload git plan branch/tag/sha with %v", parsedGitPlan.Head)
	}

	planPath := path.Join(localShuttleDirectoryPath, "plan")

	plansAlreadyValidated := strings.Split(os.Getenv("SHUTTLE_PLANS_ALREADY_VALIDATED"), string(os.PathListSeparator))
	for _, planAlreadyValidated := range plansAlreadyValidated {
		if planAlreadyValidated == planPath {
			uii.Verboseln("Shuttle already validated plan. Skipping further plan validation")
			return planPath, nil
		}
	}

	if fileAvailable(planPath) {
		status := getStatus(planPath)

		if status.mergeState {
			return "", errors.NewExitCode(9, "Plan's cloned output is in merge state. Please resolve merge conflicts and the try again")
		} else if status.changes {
			uii.EmphasizeInfoln("Found %v files locally changed in plan", len(status.files))
			uii.EmphasizeInfoln("Skipping plan pull because of changes")
		} else {
			if skipGitPlanPulling {
				uii.Verboseln("Skipping git plan pulling")
				return planPath, nil
			}
			err := gitCmd("fetch origin", planPath, uii)
			if err != nil {
				return "", err
			}
			err = gitCmd(fmt.Sprintf("checkout %s", parsedGitPlan.Head), planPath, uii)
			if err != nil {
				return "", err
			}
			status := getStatus(planPath)
			if !status.isDetached {
				uii.Infoln("Pulling latest plan changes on %v", parsedGitPlan.Head)
				err = gitCmd(fmt.Sprintf("pull origin %v", parsedGitPlan.Head), planPath, uii)
				if err != nil {
					return "", err
				}
				status = getStatus(planPath)
			} else {
				uii.EmphasizeInfoln("Skipping plan pull because its running on detached head")
			}
			uii.Verboseln("Using %s - branch %s - commit %s", plan, status.branch, status.commit)
		}
		return planPath, nil
	} else {
		err := os.MkdirAll(localShuttleDirectoryPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("create '%s' directory: %w", localShuttleDirectoryPath, err)
		}

		var cloneArg string
		if parsedGitPlan.Protocol == "https" {
			cloneArg = "https://" + parsedGitPlan.Repository
		} else if parsedGitPlan.Protocol == "ssh" {
			cloneArg = parsedGitPlan.User + "@" + parsedGitPlan.Repository
		} else {
			panic(fmt.Sprintf("Unknown protocol '%s'", parsedGitPlan.Protocol))
		}

		uii.Infoln("Cloning plan %s", cloneArg)
		err = gitCmd(fmt.Sprintf("clone %v --branch %v plan", cloneArg, parsedGitPlan.Head), localShuttleDirectoryPath, uii)
		if err != nil {
			return "", err
		}
	}

	return planPath, nil
}

func RunGitPlanCommand(command string, plan string, uii *ui.UI) {

	cmdOptions := go_cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	execCmd := go_cmd.NewCmdOptions(cmdOptions, "sh", "-c", "cd '"+plan+"'; git "+command)
	execCmd.Env = os.Environ()

	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)

		for execCmd.Stdout != nil || execCmd.Stderr != nil {
			select {
			case line, open := <-execCmd.Stdout:
				if !open {
					execCmd.Stdout = nil
					continue
				}
				uii.Infoln("%s", line)
			case line, open := <-execCmd.Stderr:
				if !open {
					execCmd.Stderr = nil
					continue
				}
				uii.Infoln("%s", line)
			}
		}
	}()

	status := <-execCmd.Start()
	<-doneChan

	if status.Exit > 0 {
		os.Exit(status.Exit)
	}
}

func fileAvailable(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func gitCmd(command string, dir string, uii *ui.UI) error {
	cmdOptions := go_cmd.Options{
		Buffered:  true,
		Streaming: true,
	}
	execCmd := go_cmd.NewCmdOptions(cmdOptions, "sh", "-c", "cd '"+dir+"'; git "+command)
	execCmd.Env = os.Environ()
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		for execCmd.Stdout != nil || execCmd.Stderr != nil {
			select {
			case line, open := <-execCmd.Stdout:
				if !open {
					execCmd.Stdout = nil
					continue
				}
				uii.Verboseln("git> %s", line)
			case line, open := <-execCmd.Stderr:
				if !open {
					execCmd.Stderr = nil
					continue
				}
				uii.Verboseln("git> %s", line)
			}
		}
	}()
	status := <-execCmd.Start()
	<-doneChan

	if status.Exit > 0 {
		return errors.NewExitCode(4, "Failed executing git command `%s` in `%s`. Got exit code: %v\n%s", command, dir, status.Exit, strings.Join(status.Stderr, "\n"))
	}
	return nil
}
