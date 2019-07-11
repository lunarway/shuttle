package executors

import (
	"testing"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestValidateUnknownArgs(t *testing.T) {
	tt := []struct {
		name       string
		scriptArgs []config.ShuttleScriptArgs
		inputArgs  map[string]string
		output     []string
	}{
		{
			name:       "no input or script",
			scriptArgs: nil,
			inputArgs:  nil,
			output:     nil,
		},
		{
			name:       "single input without script",
			scriptArgs: nil,
			inputArgs: map[string]string{
				"foo": "1",
			},
			output: []string{
				"'foo' unknown",
			},
		},
		{
			name:       "multiple input without script",
			scriptArgs: nil,
			inputArgs: map[string]string{
				"foo": "1",
				"bar": "2",
			},
			output: []string{
				"'foo' unknown",
				"'bar' unknown",
			},
		},
		{
			name: "single input and script",
			scriptArgs: []config.ShuttleScriptArgs{
				{
					Name: "foo",
				},
			},
			inputArgs: map[string]string{
				"foo": "1",
			},
			output: nil,
		},
		{
			name: "multple input and script",
			scriptArgs: []config.ShuttleScriptArgs{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			},
			inputArgs: map[string]string{
				"foo": "1",
				"bar": "2",
			},
			output: nil,
		},
		{
			name: "multple input and script with one unknown",
			scriptArgs: []config.ShuttleScriptArgs{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			},
			inputArgs: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			output: []string{
				"'baz' unknown",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := validateUnknownArgs(tc.scriptArgs, tc.inputArgs)
			assert.Equal(t, tc.output, output, "output not as expected")
		})
	}
}
