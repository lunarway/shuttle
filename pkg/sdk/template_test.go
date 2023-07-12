package sdk

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTemplate(t *testing.T) {
	tt := []struct {
		name         string
		templatePath string
		output       string
		templateCtx  TemplateContext
		err          error
	}{
		{
			name:         "valid",
			templatePath: "testdata/valid.yaml",
			templateCtx: TemplateContext{
				Vars: map[string]string{
					"bar": "bar",
				},
			},
			output: `foo: bar
`,
			err: nil,
		},
		{
			name:         "read content of unknown file",
			templatePath: "testdata/get_file_content.yaml",
			templateCtx: TemplateContext{
				Vars: map[string]string{
					"path": "unknown.yaml",
				},
			},
			output: ``,
			err: errors.New(
				"template: get_file_content.yaml:2:3: executing \"get_file_content.yaml\" at <getFileContent (string \"path\" .Vars)>: error calling getFileContent: open unknown.yaml: no such file or directory",
			),
		},
		{
			name:         "read content of file",
			templatePath: "testdata/get_file_content.yaml",
			templateCtx: TemplateContext{
				Vars: map[string]string{
					"path": "testdata/valid.yaml",
				},
			},
			output: `bar: baz
foo: {{ get "bar" .Vars }}

`,
			err: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			err := renderTemplate(tc.templatePath, tc.name, &output, tc.templateCtx, "{{", "}}")

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.output, output.String())
		})
	}
}
