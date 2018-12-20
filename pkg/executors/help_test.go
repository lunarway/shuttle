package executors_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/executors"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	scriptMap := func(args ...interface{}) map[string]config.ShuttlePlanScript {
		m := make(map[string]config.ShuttlePlanScript)
		for i := 0; i < len(args); i = i + 2 {
			m[args[i].(string)] = args[i+1].(config.ShuttlePlanScript)
		}
		return m
	}
	scripts := func(description string, args ...config.ShuttleScriptArgs) config.ShuttlePlanScript {
		return config.ShuttlePlanScript{
			Description: description,
			Args:        args,
		}
	}
	arg := func(name string, required bool, description string) config.ShuttleScriptArgs {
		return config.ShuttleScriptArgs{
			Name:        name,
			Required:    required,
			Description: description,
		}
	}
	tt := []struct {
		name    string
		scripts map[string]config.ShuttlePlanScript
		script  string
		err     error
		output  string
	}{
		{
			name:    "no scripts",
			scripts: nil,
			script:  "test",
			err:     errors.New("unrecognized script"),
		},
		{
			name:    "no script matches",
			scripts: scriptMap("build", scripts("build stuff")),
			script:  "test",
			err:     errors.New("unrecognized script"),
		},
		{
			name:    "script without arguments",
			scripts: scriptMap("build", scripts("A script to build stuff")),
			script:  "build",
			output: `A script to build stuff

`,
		},
		{
			name:    "script with argument",
			scripts: scriptMap("test", scripts("A script to test stuff", arg("long", false, "Run long running tests"))),
			script:  "test",
			output: `A script to test stuff

Available arguments:
  long        Run long running tests
`,
		},
		{
			name:    "script with required argument",
			scripts: scriptMap("test", scripts("A script to test stuff", arg("long", true, "Run long running tests"))),
			script:  "test",
			output: `A script to test stuff

Available arguments:
  long (required)       Run long running tests
`,
		},
		{
			name:    "script with multiple arguments",
			scripts: scriptMap("test", scripts("A script to test stuff", arg("long", true, "Run long running tests"), arg("short", false, "Run short tests"))),
			script:  "test",
			output: `A script to test stuff

Available arguments:
  long (required)       Run long running tests
  short                 Run short tests
`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := bytes.Buffer{}

			err := executors.Help(tc.scripts, tc.script, &output, "")

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error(), "output error not as expected")
			} else {
				assert.NoError(t, err, "no output error expected")
			}
			assert.Equal(t, tc.output, output.String(), "output not as expected")
		})
	}
}

func TestHelp_customTemplate(t *testing.T) {
	scriptMap := func(args ...interface{}) map[string]config.ShuttlePlanScript {
		m := make(map[string]config.ShuttlePlanScript)
		for i := 0; i < len(args); i = i + 2 {
			m[args[i].(string)] = args[i+1].(config.ShuttlePlanScript)
		}
		return m
	}
	scripts := func(description string, args ...config.ShuttleScriptArgs) config.ShuttlePlanScript {
		return config.ShuttlePlanScript{
			Description: description,
			Args:        args,
		}
	}
	arg := func(name string, required bool, description string) config.ShuttleScriptArgs {
		return config.ShuttleScriptArgs{
			Name:        name,
			Required:    required,
			Description: description,
		}
	}
	tt := []struct {
		name     string
		scripts  map[string]config.ShuttlePlanScript
		script   string
		template string
		err      error
		output   string
	}{
		{
			name:     "script without arguments",
			scripts:  scriptMap("build", scripts("A script to build stuff")),
			script:   "build",
			template: `{{.Args}}`,
			output:   `[]`,
		},
		{
			name:     "ranging args",
			scripts:  scriptMap("test", scripts("A script to test stuff", arg("long", false, "Run long running tests"))),
			script:   "test",
			template: `{{- range $i, $arg := .Args -}}{{$arg.Name}}{{end}}`,
			output:   `long`,
		},
		{
			name:     "invalid template",
			scripts:  scriptMap("test", scripts("A script to test stuff", arg("long", true, "Run long running tests"))),
			script:   "test",
			template: `{{.Args`,
			err:      errors.New("invalid template: template: runHelp:1: unclosed action"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := bytes.Buffer{}

			err := executors.Help(tc.scripts, tc.script, &output, tc.template)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error(), "output error not as expected")
			} else {
				assert.NoError(t, err, "no output error expected")
			}
			assert.Equal(t, tc.output, output.String(), "output not as expected")
		})
	}
}
