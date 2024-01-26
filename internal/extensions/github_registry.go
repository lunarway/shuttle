package extensions

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type gitHubRegistry struct {
	client *githubClient
}

func (g *gitHubRegistry) Publish(ctx context.Context, extFile *shuttleExtensionsFile, version string) error {
	release, err := g.client.GetRelease(ctx, extFile, version)
	if err != nil {
		return err
	}

	sha, err := g.client.GetFile(ctx, extFile)
	if err != nil {
		log.Printf("failed to find file: %s", err.Error())
		// Ignore file as it probably means that the file wasn't there
	}

	if err := g.client.UpsertFile(ctx, extFile, release, version, sha); err != nil {
		return err
	}

	log.Println("done")
	return nil
}

// Get isn't implemented yet for GitHubRegistry
func (*gitHubRegistry) Get(ctx context.Context) error {
	panic("unimplemented")
}

// Update isn't implemented yet for GitHubRegistry
func (*gitHubRegistry) Update(ctx context.Context) error {
	panic("unimplemented")
}

func newGitHubRegistry() (Registry, error) {
	client, err := newGitHubClient()
	if err != nil {
		return nil, err
	}

	return &gitHubRegistry{
		client: client,
	}, nil
}

type githubClient struct {
	accessToken string
	httpClient  *http.Client
}

func newGitHubClient() (*githubClient, error) {
	var token string
	if accessToken := os.Getenv("SHUTTLE_EXTENSIONS_GITHUB_ACCESS_TOKEN"); accessToken != "" {
		token = accessToken
	} else if accessToken := os.Getenv("GITHUB_ACCESS_TOKEN"); accessToken != "" {
		token = accessToken
	}

	if token == "" {
		return nil, errors.New("GITHUB_ACCESS_TOKEN was not set")
	}

	return &githubClient{
		accessToken: token,
		httpClient:  http.DefaultClient,
	}, nil
}

func (gc *githubClient) GetFile(ctx context.Context, shuttleExtensionsFile *shuttleExtensionsFile) (string, error) {
	owner, repo, ok := strings.Cut(*shuttleExtensionsFile.Registry.GitHub, "/")
	if !ok {
		return "", fmt.Errorf("failed to find owner and repo in registry: %s", *shuttleExtensionsFile.Registry.GitHub)
	}

	extensionsFile, err := githubClientDo[any, githubFileShaResp](
		ctx,
		gc,
		http.MethodGet,
		fmt.Sprintf(
			"/repos/%s/%s/contents/%s",
			owner,
			repo,
			getRemoteRegistryExtensionPathFile(shuttleExtensionsFile.Name),
		),
		nil,
	)
	if err != nil {
		return "", err
	}

	return extensionsFile.Sha, nil
}

func (gc *githubClient) UpsertFile(ctx context.Context, shuttleExtensionsFile *shuttleExtensionsFile, releaseInformation *githubReleaseInformation, version string, sha string) error {
	registryExtensionsReq := registryExtension{
		Name:         shuttleExtensionsFile.Name,
		Description:  shuttleExtensionsFile.Description,
		Version:      version,
		DownloadUrls: make([]registryExtensionDownloadLink, 0),
	}

	for _, releaseAsset := range releaseInformation.Assets {
		arch, os, err := releaseAsset.ParseDownloadLink(shuttleExtensionsFile.Name)
		if err != nil {
			log.Printf("file did not match an actual binary: %s, %s", releaseAsset.DownloadUrl, err.Error())
			continue
		}

		downloadLink := registryExtensionDownloadLink{
			Architecture: arch,
			Os:           os,
			Url:          releaseAsset.DownloadUrl,
			Provider:     "github-release",
		}

		registryExtensionsReq.DownloadUrls = append(registryExtensionsReq.DownloadUrls, downloadLink)
	}

	upsertRequest, err := newGitHubUpsertRequest(
		shuttleExtensionsFile.Name,
		version,
		registryExtensionsReq,
		sha,
	)
	if err != nil {
		return err
	}

	owner, repo, ok := strings.Cut(*shuttleExtensionsFile.Registry.GitHub, "/")
	if !ok {
		return fmt.Errorf("failed to find owner and repo in registry: %s", *shuttleExtensionsFile.Registry.GitHub)
	}

	_, err = githubClientDo[githubUpsertFileRequest, any](
		ctx,
		gc,
		http.MethodPut,
		fmt.Sprintf(
			"/repos/%s/%s/contents/%s",
			owner,
			repo,
			getRemoteRegistryExtensionPathFile(shuttleExtensionsFile.Name),
		),
		upsertRequest,
	)
	if err != nil {
		return err
	}

	return nil
}

func (gc *githubClient) GetRelease(ctx context.Context, shuttleExtensionsFile *shuttleExtensionsFile, version string) (*githubReleaseInformation, error) {
	release, err := githubClientDo[any, githubReleaseInformation](
		ctx,
		gc,
		http.MethodGet,
		fmt.Sprintf(
			"/repos/%s/%s/releases/tags/%s",
			shuttleExtensionsFile.Provider.GitHubRelease.Owner,
			shuttleExtensionsFile.Provider.GitHubRelease.Repo,
			version,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if len(release.Assets) == 0 {
		return nil, errors.New("found no releases for github release")
	}

	return release, nil
}

type githubReleaseAsset struct {
	DownloadUrl string `json:"browser_download_url"`
}

func (gra *githubReleaseAsset) ParseDownloadLink(name string) (arch string, os string, err error) {
	components := strings.Split(gra.DownloadUrl, "/")
	if len(components) < 3 {
		return "", "", errors.New("failed to find a proper github download link")
	}

	file := components[len(components)-1]

	rest, ok := strings.CutPrefix(file, name)
	if !ok {
		return "", "", errors.New("file link did not contain extension name")
	}
	rest = strings.TrimPrefix(rest, "-")

	os, arch, ok = strings.Cut(rest, "-")
	if !ok {
		return "", "", errors.New("file did not match os-arch")
	}

	return arch, os, nil
}

type githubReleaseInformation struct {
	Assets []githubReleaseAsset `json:"assets"`
}

type githubFileShaResp struct {
	Sha string `json:"sha"`
}

type githubCommitter struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type githubUpsertFileRequest struct {
	Message   string          `json:"message"`
	Committer githubCommitter `json:"committer"`
	Content   string          `json:"content"`
	Sha       *string         `json:"sha"`
}

func newGitHubUpsertRequest(name string, version string, registryExtensionsReq registryExtension, sha string) (*githubUpsertFileRequest, error) {
	committerName := os.Getenv("GITHUB_COMMITTER_NAME")
	if committerName == "" {
		return nil, errors.New("GITHUB_COMMITTER_NAME was not found")
	}

	committerEmail := os.Getenv("GITHUB_COMMITTER_EMAIL")
	if committerEmail == "" {
		return nil, errors.New("GITHUB_COMMITTER_EMAIL was not found")
	}

	content, err := json.MarshalIndent(registryExtensionsReq, "", "  ")
	if err != nil {
		return nil, err
	}

	contentB64 := base64.StdEncoding.EncodeToString(content)

	req := &githubUpsertFileRequest{
		Message: fmt.Sprintf("chore(extensions): updating %s to %s", name, version),
		Committer: githubCommitter{
			Name:  committerName,
			Email: committerEmail,
		},
		Content: contentB64,
	}

	if sha != "" {
		req.Sha = &sha
	}

	return req, nil
}

func githubClientDo[TReq any, TResp any](ctx context.Context, githubClient *githubClient, method string, path string, reqBody *TReq) (*TResp, error) {
	var bodyReader io.Reader
	if reqBody != nil {
		contents, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(contents)
	}

	url := fmt.Sprintf("https://api.github.com/%s", strings.TrimPrefix(path, "/"))

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		bodyReader,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", githubClient.accessToken))

	resp, err := githubClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 399 {
		var message githubMessage
		if err := json.Unmarshal(respContent, &message); err != nil {
			return nil, fmt.Errorf("failed to unmarshal resp: %w", err)
		}

		return nil, fmt.Errorf("failed github request with: %s", message.Message)
	}

	var returnObject TResp
	if err := json.Unmarshal(respContent, &returnObject); err != nil {
		return nil, err
	}

	return &returnObject, nil
}

type githubMessage struct {
	Message string `json:"message"`
}
