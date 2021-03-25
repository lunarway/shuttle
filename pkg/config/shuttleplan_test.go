package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuttlePlanConfiguration_Load(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		err    error
		config ShuttlePlanConfiguration
	}{
		{
			name:  "empty path",
			input: "",
		},
		{
			name:  "unknown field",
			input: "testdata/unknown_field",
			err:   errors.New("exit code 1 - Failed to load plan configuration from 'testdata/unknown_field/plan.yaml': yaml: unmarshal errors:\n  line 1: field unknown not found in type config.ShuttlePlanConfiguration\n\nThis is likely an issue with the referenced plan. Please, contact the plan maintainers."),
		},
		{
			name:  "unknown file",
			input: "testdata/unknown_file",
			err:   errors.New("exit code 2 - Failed to open plan configuration: open testdata/unknown_file/plan.yaml: no such file or directory\n\nMake sure you are in a project using shuttle and that a 'shuttle.yaml' file is available."),
		},
		{
			name:  "valid",
			input: "testdata/valid",
			err:   nil,
			config: ShuttlePlanConfiguration{
				Vars: map[string]interface{}{
					"shared": "var",
				},
				Scripts: map[string]ShuttlePlanScript{
					"hello": {
						Description: "Say hello",
						Actions: []ShuttleAction{
							{
								Shell: `echo "Hello world"`,
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := &ShuttlePlanConfiguration{}
			var err error
			c, err = c.Load(tc.input)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error(), "error not as expected")
			} else {
				assert.NoError(t, err, "unexpected error")
			}
			assert.Equal(t, tc.config, *c, "config not as expected")
		})
	}
}
