package extensions

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	bearer, err := getGithubToken()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", bearer))
	req.Header.Add("Accept", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.RemoveAll(dest); err != nil {
		log.Printf("failed to remove extension before downloading new: %s, please try again", err.Error())
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
