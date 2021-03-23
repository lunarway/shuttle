package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuttleConfig_getConf(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		err    error
		config ShuttleConfig
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
			err:   errors.New("exit code 2 - Failed to load shuttle configuration: open testdata/unknown_file/shuttle.yaml: no such file or directory\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available."),
		},
		{
			name:  "valid",
			input: "testdata/valid",
			err:   nil,
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
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := &ShuttleConfig{}
			var err error
			c, err = c.getConf(tc.input)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error(), "error not as expected")
			} else {
				assert.NoError(t, err, "unexpected error")
			}
			assert.Equal(t, tc.config, *c, "config not as expected")
		})
	}
}
