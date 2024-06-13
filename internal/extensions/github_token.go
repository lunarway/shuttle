package extensions

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func getGithubToken() (string, error) {
	if accessToken := os.Getenv("SHUTTLE_EXTENSIONS_GITHUB_ACCESS_TOKEN"); accessToken != "" {
		return accessToken, nil
	} else if accessToken := os.Getenv("GITHUB_ACCESS_TOKEN"); accessToken != "" {
		return accessToken, nil
	} else {
		accessToken, err := getToken()
		if err != nil {
			return "", err
		}

		return accessToken, nil
	}
}

func getToken() (string, error) {
	tokenRaw, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", errors.New("github-cli (gh) is not installed")
		}

		return "", err
	}

	token := string(tokenRaw)

	if token != "" {
		return strings.TrimSpace(token), nil
	}

	return "", errors.New("no github token available (please sign in `gh auth login`)")
}
