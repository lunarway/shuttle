package git

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/lunarway/shuttle/pkg/output"

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
func GetGitPlan(plan string, localShuttleDirectoryPath string, verbose bool) string {
	parsedGitPlan := parseGitPlan(plan)
	planPath := path.Join(localShuttleDirectoryPath, "plan")

	if fileAvailable(planPath) {
		output.Verbose(verbose, "Pulling latest git changes")
		gitCmd("pull origin", planPath, verbose)
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

		output.Verbose(verbose, "Cloning repository %s", cloneArg)
		gitCmd("clone "+cloneArg+" plan", localShuttleDirectoryPath, verbose)
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

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
	// never reached
	panic(true)
	return nil, nil
}

func gitCmd(command string, dir string, printOutput bool) {
	cmdOptions := go_cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	execCmd := go_cmd.NewCmdOptions(cmdOptions, "sh", "-c", "cd '"+dir+"'; git "+command)
	execCmd.Env = os.Environ()
	go func() {
		for {
			select {
			case line := <-execCmd.Stdout:
				if printOutput {
					fmt.Println("git> " + line)
				}
			case line := <-execCmd.Stderr:
				if printOutput {
					fmt.Fprintln(os.Stderr, "git> "+line)
				}
			}
		}
	}()
	status := <-execCmd.Start()
	for len(execCmd.Stdout) > 0 || len(execCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	if status.Exit > 0 {
		output.ExitWithErrorCode(4, fmt.Sprintf("Failed executing git command `%s` in `%s`\nExit code: %v", command, dir, status.Exit))
	}
}
