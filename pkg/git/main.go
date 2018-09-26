package git

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"strings"
)

type gitPlan struct {
	isGitPlan  bool
	protocol   string
	user       string
	host       string
	repository string
}

var gitRegex = regexp.MustCompile(`^git://((?P<user>[^@]+)@)?(?P<repository>(?P<host>[^:]+):(?P<path>.*))$|^(?P<protocol>https)://(?P<repository>.*\.git)$`)

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

	return gitPlan{
		isGitPlan:  true,
		protocol:   protocol,
		user:       result["user"],
		host:       result["host"],
		repository: result["repository"],
	}
}

// IsGitPlan returns true if specified plan is a git plan
func IsGitPlan(plan string) bool {
	parsedGitPlan := parseGitPlan(plan)
	return parsedGitPlan.isGitPlan
}

// GetGitPlan will pull git repository and return its path
func GetGitPlan(plan string, localShuttleDirectoryPath string) string {
	// We need the user to find the homedir.

	parsedGitPlan := parseGitPlan(plan)
	planPath := path.Join(localShuttleDirectoryPath, "plan")

	if fileAvailable(planPath) {

		execCmd := exec.Command("git", "pull", "origin")
		execCmd.Env = append(os.Environ())
		execCmd.Dir = planPath

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := execCmd.StdoutPipe()
		stderrIn, _ := execCmd.StderrPipe()
		startErr := execCmd.Start()
		checkIfError(startErr)

		go func() {
			stdout, errStdout = copyAndCapture(ioutil.Discard, stdoutIn)
		}()

		go func() {
			stderr, errStderr = copyAndCapture(ioutil.Discard, stderrIn)
		}()

		err := execCmd.Wait()
		checkIfError(err)

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

		execCmd := exec.Command("git", "clone", cloneArg, "plan")
		execCmd.Env = append(os.Environ())
		execCmd.Dir = localShuttleDirectoryPath

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := execCmd.StdoutPipe()
		stderrIn, _ := execCmd.StderrPipe()
		startErr := execCmd.Start()
		checkIfError(startErr)

		go func() {
			stdout, errStdout = copyAndCapture(ioutil.Discard, stdoutIn)
		}()

		go func() {
			stderr, errStderr = copyAndCapture(ioutil.Discard, stderrIn)
		}()

		err := execCmd.Wait()

		if err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("Could not clone %s\ngit output:%v\n%v", plan, string(stdout), string(stderr)))
			os.Exit(3)
		}

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
