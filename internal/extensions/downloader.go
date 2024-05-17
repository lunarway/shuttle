package extensions

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Downloader interface {
	Download(ctx context.Context, dest string) error
}

func NewDownloader(downloadLink *registryExtensionDownloadLink) (Downloader, error) {
	switch downloadLink.Provider {
	case "github-release":
		return newGitHubReleaseDownloader(downloadLink), nil
	default:
		return nil, fmt.Errorf("invalid provider type: %s", downloadLink.Provider)
	}
}

type gitHubReleaseDownloader struct {
	link *registryExtensionDownloadLink
}

func newGitHubReleaseDownloader(downloadLink *registryExtensionDownloadLink) Downloader {
	return &gitHubReleaseDownloader{
		link: downloadLink,
	}
}

func (d *gitHubReleaseDownloader) Download(ctx context.Context, dest string) error {
	client := http.DefaultClient
	client.Timeout = time.Second * 60

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.link.Url, nil)
	if err != nil {
		return err
	}

	var bearer string
	if accessToken := os.Getenv("SHUTTLE_EXTENSIONS_GITHUB_ACCESS_TOKEN"); accessToken != "" {
		bearer = accessToken
	} else if accessToken := os.Getenv("GITHUB_ACCESS_TOKEN"); accessToken != "" {
		bearer = accessToken
	} else if accessToken, ok := getToken(); ok {
		bearer = accessToken
	}

	if bearer == "" {
		return errors.New("failed to find a valid authorization token for github. Please make sure you're logged into github-cli (gh), or have followed the setup documentation")
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", bearer))
	req.Header.Add("Accept", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.RemoveAll(dest); err != nil {
		log.Printf("failed to remove extension before downloading new: %s", err.Error())
	}

	extensionBinary, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer extensionBinary.Close()
	extensionBinary.Chmod(0o755)

	if _, err := io.Copy(extensionBinary, resp.Body); err != nil {
		return err
	}

	return nil
}

func getToken() (string, bool) {
	tokenRaw, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return "", false
	}

	token := string(tokenRaw)

	if token != "" {
		return strings.TrimSpace(token), false
	}

	return "", false
}
