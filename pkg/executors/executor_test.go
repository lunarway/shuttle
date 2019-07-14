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
		output     []validationError
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
			output: []validationError{
				{"foo", "unknown"},
			},
		},
		{
			name:       "multiple input without script",
			scriptArgs: nil,
			inputArgs: map[string]string{
				"foo": "1",
				"bar": "2",
			},
			output: []validationError{
				{"bar", "unknown"},
				{"foo", "unknown"},
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
				"bar": "2",
				"foo": "1",
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
			output: []validationError{
				{"baz", "unknown"},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := validateUnknownArgs(tc.scriptArgs, tc.inputArgs)
			// sort as the order is not guarenteed by validateUnknownArgs
			sortValidationErrors(output)
			assert.Equal(t, tc.output, output, "output not as expected")
		})
	}
}

func TestSortValidationErrors(t *testing.T) {
	tt := []struct {
		name   string
		input  []validationError
		output []validationError
	}{
		{
			name: "sorted",
			input: []validationError{
				{"bar", ""},
				{"baz", ""},
				{"foo", ""},
			},
			output: []validationError{
				{"bar", ""},
				{"baz", ""},
				{"foo", ""},
			},
		},
		{
			name: "not sorted",
			input: []validationError{
				{"baz", ""},
				{"foo", ""},
				{"bar", ""},
			},
			output: []validationError{
				{"bar", ""},
				{"baz", ""},
				{"foo", ""},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sortValidationErrors(tc.input)
			assert.Equal(t, tc.output, tc.input, "output not as expected")
		})
	}
}
