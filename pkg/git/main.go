package git

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/lunarway/shuttle/pkg/ui"

	go_cmd "github.com/go-cmd/cmd"
)

type gitPlan struct {
	isGitPlan  bool
	protocol   string
	user       string
	host       string
	repository string
}

var gitRegex = regexp.MustCompile(`^git://((?P<user>[^@]+)@)?(?P<repository1>(?P<host>[^:]+):(?P<path>.*))$|^(?P<protocol>https)://(?P<repository2>.*\.git)$`)

func parseGitPlan(plan string) gitPlan {
	if !gitRegex.MatchString(plan) {
		return gitPlan{
			isGitPlan: false,
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

	return gitPlan{
		isGitPlan:  true,
		protocol:   protocol,
		user:       result["user"],
		host:       result["host"],
		repository: repository,
	}
}

// IsGitPlan returns true if specified plan is a git plan
func IsGitPlan(plan string) bool {
	parsedGitPlan := parseGitPlan(plan)
	return parsedGitPlan.isGitPlan
}

// GetGitPlan will pull git repository and return its path
func GetGitPlan(plan string, localShuttleDirectoryPath string, uii ui.UI) string {
	parsedGitPlan := parseGitPlan(plan)
	planPath := path.Join(localShuttleDirectoryPath, "plan")

	plansAlreadyValidated := strings.Split(os.Getenv("SHUTTLE_PLANS_ALREADY_VALIDATED"), string(os.PathListSeparator))
	for _, planAlreadyValidated := range plansAlreadyValidated {
		if planAlreadyValidated == planPath {
			uii.VerboseLn("Shuttle already validated plan. Skipping further plan validation")
			return planPath
		}
	}

	if fileAvailable(planPath) {
		status := getStatus(planPath)

		if status.mergeState {
			uii.ExitWithErrorCode(9, "Plan's cloned output is in merge state. Please resolve merge conflicts and the try again")
		} else if status.changes {
			uii.EmphasizeInfoLn("Found %v files locally changed in plan", len(status.files))
			uii.EmphasizeInfoLn("Skipping plan pull because of changes")
		} else {
			uii.InfoLn("Pulling latest plan changes")
			gitCmd("pull origin", planPath, uii)
		}
	} else {
		os.MkdirAll(localShuttleDirectoryPath, os.ModePerm)

		var cloneArg string
		if parsedGitPlan.protocol == "https" {
			cloneArg = "https://" + parsedGitPlan.repository
		} else if parsedGitPlan.protocol == "ssh" {
			cloneArg = parsedGitPlan.user + "@" + parsedGitPlan.repository
		} else {
			panic(fmt.Sprintf("Unknown protocol '%s'", parsedGitPlan.protocol))
		}

		uii.InfoLn("Cloning plan %s", cloneArg)
		gitCmd("clone "+cloneArg+" plan", localShuttleDirectoryPath, uii)
	}

	return planPath
}

func isMatching(r string, content string) bool {
	match, err := regexp.MatchString(r, content)
	if err != nil {
		panic(err)
	}
	return match
}

func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func fileAvailable(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func expandHome(path string) string {
	usr, err := user.Current()
	checkIfError(err)
	return strings.Replace(path, "~/", usr.HomeDir+"/", 1)
}

func gitCmd(command string, dir string, uii ui.UI) {
	cmdOptions := go_cmd.Options{
		Buffered:  true,
		Streaming: true,
	}
	execCmd := go_cmd.NewCmdOptions(cmdOptions, "sh", "-c", "cd '"+dir+"'; git "+command)
	execCmd.Env = os.Environ()
	go func() {
		for {
			select {
			case line := <-execCmd.Stdout:
				uii.VerboseLn("git> " + line)
			case line := <-execCmd.Stderr:
				uii.VerboseLn("git> " + line)
			}
		}
	}()
	status := <-execCmd.Start()
	for len(execCmd.Stdout) > 0 || len(execCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	if status.Exit > 0 {
		uii.ExitWithErrorCode(4, "Failed executing git command `%s` in `%s`. Got exit code: %v\n%s", command, dir, status.Exit, strings.Join(status.Stderr, "\n"))
	}
}
