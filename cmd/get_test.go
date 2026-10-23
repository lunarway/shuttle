package cmd

import (
	"testing"
)

func TestGet(t *testing.T) {
	testCases := []testCase{
		{
			name:      "local plan",
			input:     args("-p", "testdata/project", "get", "service"),
			stdoutput: "shuttle",
			erroutput: "",
			err:       nil,
		},
		{
			name: "local plan with templating",
			input: args(
				"-p",
				"testdata/project",
				"get",
				"nested",
				"--template",
				"{{ range $k, $v := . }}{{ $k }}{{ end }}",
			),
			stdoutput: "subvar",
			erroutput: "",
			err:       nil,
		},
		{
			name: "local plan with templating function",
			input: args(
				"-p",
				"testdata/project",
				"get",
				"nested",
				"--template",
				`{{ range objectArray "sub" . }}{{ .Key }}{{ end }}`,
			),
			stdoutput: "field",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "bool",
			input:     args("-p", "testdata/project", "get", "boolVar"),
			stdoutput: "false",
			erroutput: "",
			err:       nil,
		},
		{
			name: "falls back to plan variables",
			input: args(
				"--project",
				"./testdata/project-local/service",
				"--plan",
				"./testdata/project-local/plan",
				"get",
				"language",
			),
			stdoutput: "go",
			erroutput: "Using overloaded plan ./testdata/project-local/plan\n",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
