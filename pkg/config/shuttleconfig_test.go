package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuttleConfig_getConf(t *testing.T) {
	tt := []struct {
		name       string
		input      string
		strictMode bool
		err        error
		config     ShuttleConfig
		foundPath  string
	}{
		{
			name:  "empty path",
			input: "",
		},
		{
			name:  "unknown field",
			input: "testdata/unknown_field",
			err:   errors.New("exit code 2 - Failed to parse shuttle configuration: yaml: unmarshal errors:\n  line 1: field nothing not found in type config.ShuttleConfig\n\nMake sure your 'shuttle.yaml' is valid."),
		},
		{
			name:  "unknown file",
			input: "testdata/unknown_file",
			err:   errors.New("exit code 2 - Failed to load shuttle configuration: shuttle.yaml file not found\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available."),
		},
		{
			name:  "absolute path to unknown file",
			input: "/tmp/shuttle-test/unknown",
			err:   errors.New("exit code 2 - Failed to load shuttle configuration: shuttle.yaml file not found\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available."),
		},
		{
			name:      "valid",
			input:     "testdata/valid",
			err:       nil,
			foundPath: "testdata/valid",
			config: ShuttleConfig{
				Plan:    ".",
				PlanRaw: ".",
				Variables: map[string]interface{}{
					"squad": "nasa",
				},
				Scripts: map[string]ShuttlePlanScript{
					"shout": {
						Description: "Shout hello",
						Actions: []ShuttleAction{
							{
								Shell: `echo "HELLO WORLD"`,
							},
						},
					},
				},
			},
		},
		{
			name:      "subdir of shuttle.yaml file",
			input:     "testdata/valid/subdir",
			err:       nil,
			foundPath: "testdata/valid",
			config: ShuttleConfig{
				Plan:    ".",
				PlanRaw: ".",
				Variables: map[string]interface{}{
					"squad": "nasa",
				},
				Scripts: map[string]ShuttlePlanScript{
					"shout": {
						Description: "Shout hello",
						Actions: []ShuttleAction{
							{
								Shell: `echo "HELLO WORLD"`,
							},
						},
					},
				},
			},
		},
		{
			name:       "subdir of shuttle.yaml file in strict mode",
			input:      "testdata/valid/subdir",
			strictMode: true,
			err:        errors.New("exit code 2 - Failed to load shuttle configuration: shuttle.yaml file not found\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available."),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := &ShuttleConfig{}

			path, err := c.getConf(tc.input, tc.strictMode)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error(), "error not as expected")
			} else {
				assert.NoError(t, err, "unexpected error")
			}
			assert.Equal(t, tc.config, *c, "config not as expected")
			assert.Equal(t, tc.foundPath, path, "shuttle.yaml path not as expected")
		})
	}
}
