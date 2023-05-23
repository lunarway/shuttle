package cmd

import (
	"testing"
)

func TestTemplate(t *testing.T) {
	testCases := []testCase{
		{
			name:  "local path",
			input: args("-p", "testdata/project", "template", "../custom-template.tmpl", "GO_VERSION=1.17"),
			stdoutput: `# Custom docker file template not located inside a project
FROM golang:1.17-alpine
`,
			erroutput: "",
			err:       nil,
		},
		{
			name:  "alternative delimiters",
			input: args("-p", "testdata/project", "template", "../custom-template-alternative-delims.tmpl", "--delims", ">>,<<"),
			stdoutput: `FROM golang:1.17-alpine
LABEL svc=shuttle
`,
			erroutput: "",
			err:       nil,
		},
		{
			name:  "image tag should not be parsed as float",
			input: args("-p", "testdata/project", "template", "../Dockerfile.tmpl"),
			stdoutput: `FROM golang:1.20
`,
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
