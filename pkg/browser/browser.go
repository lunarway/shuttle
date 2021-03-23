package browser

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cli/safeexec"
	"github.com/google/shlex"
)

// This package is copied from github.com/cli/cli

// Command returns an exec.Cmd instance respecting runtime.GOOS and $BROWSER
// environment variable.
func Command(url string) (*exec.Cmd, error) {
	launcher := os.Getenv("BROWSER")
	if launcher != "" {
		return fromBrowserEnv(launcher, url)
	}
	return forOS(runtime.GOOS, url), nil
}

func forOS(goos, url string) *exec.Cmd {
	exe := "open"
	var args []string
	switch goos {
	case "darwin":
		args = append(args, url)
	case "windows":
		exe, _ = lookPath("cmd")
		r := strings.NewReplacer("&", "^&")
		args = append(args, "/c", "start", r.Replace(url))
	default:
		exe = linuxExe()
		args = append(args, url)
	}

	cmd := exec.Command(exe, args...)
	cmd.Stderr = os.Stderr
	return cmd
}

// fromBrowserEnv parses the BROWSER string based on shell splitting rules.
func fromBrowserEnv(launcher, url string) (*exec.Cmd, error) {
	args, err := shlex.Split(launcher)
	if err != nil {
		return nil, err
	}

	exe, err := lookPath(args[0])
	if err != nil {
		return nil, err
	}

	args = append(args, url)
	cmd := exec.Command(exe, args[1:]...)
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func linuxExe() string {
	exe := "xdg-open"

	_, err := lookPath(exe)
	if err != nil {
		_, err := lookPath("wslview")
		if err == nil {
			exe = "wslview"
		}
	}

	return exe
}

var lookPath = safeexec.LookPath
