package extensions

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type shuttleExtensionsRegistry struct {
	GitHub *string `json:"github" yaml:"github"`
}

type shuttleExtensionProviderGitHubRelease struct {
	Owner string `json:"owner" yaml:"owner"`
	Repo  string `json:"repo" yaml:"repo"`
}

type shuttleExtensionsProvider struct {
	GitHubRelease *shuttleExtensionProviderGitHubRelease `json:"github-release" yaml:"github-release"`
}

type shuttleExtensionsFile struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`

	Provider shuttleExtensionsProvider `json:"provider" yaml:"provider"`
	Registry shuttleExtensionsRegistry `json:"registry" yaml:"registry"`
}

func getExtensionsFile(_ context.Context) (*shuttleExtensionsFile, error) {
	templateFileContent, err := os.ReadFile("shuttle.template.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to find shuttle.template.yaml: %w", err)
	}

	var templateFile shuttleExtensionsFile
	if err := yaml.Unmarshal(templateFileContent, &templateFile); err != nil {
		return nil, fmt.Errorf("failed to parse shuttle.template.yaml: %w", err)
	}

	return &templateFile, nil
}
