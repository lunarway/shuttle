package shuttlefile

import (
	"bytes"
	"context"
	"fmt"
	"os"

	shuttleconfig "github.com/lunarway/shuttle/pkg/config"
	"gopkg.in/yaml.v3"
)

func ParseFile(ctx context.Context, path string) (*shuttleconfig.ShuttleConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := shuttleconfig.ShuttleConfig{}

	decoder := yaml.NewDecoder(bytes.NewReader(content))

	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal shuttle plan: %s", path)
	}

	switch config.PlanRaw {
	case false:
		// no plan
	default:
		config.Plan = config.PlanRaw.(string)
	}

	return &config, nil
}
