package cmd

import (
	"testing"
)

func TestGet(t *testing.T) {
	t.Cleanup(func() {
		removeShuttleDirectories(t)
	})

	testCases := []testCase{
		{
			name:      "local plan",
			input:     args("-p", "testdata/project", "get", "service"),
			stdoutput: "shuttle",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "local plan with templating",
			input:     args("-p", "testdata/project", "get", "nested", "--template", "{{ range $k, $v := . }}{{ $k }}{{ end }}"),
			stdoutput: "subvar",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "local plan with templating function",
			input:     args("-p", "testdata/project", "get", "nested", "--template", `{{ range objectArray "sub" . }}{{ .Key }}{{ end }}`),
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
	}
	executeTestCases(t, testCases)
}
