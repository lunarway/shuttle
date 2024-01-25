package extensions

import (
	"context"
	"fmt"
	"io"
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

	if accessToken := os.Getenv("SHUTTLE_EXTENSIONS_GITHUB_ACCESS_TOKEN"); accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	} else if accessToken := os.Getenv("GITHUB_ACCESS_TOKEN"); accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	extensionBinary, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer extensionBinary.Close()

	if _, err := io.Copy(extensionBinary, resp.Body); err != nil {
		return err
	}

	return nil
}
