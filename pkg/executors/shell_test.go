package executors

import (
	"context"
	"testing"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestExecuteShell(t *testing.T) {
	tt := []struct {
		name   string
		script string
		err    error
	}{
		{
			name:   "cat file with line over 80k characters",
			script: "cat testdata/large.log",
			err:    nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			err := executeShell(context.Background(), ActionExecutionContext{
				ScriptContext: ScriptExecutionContext{
					Project: config.ShuttleProjectContext{
						ProjectPath: ".",
					},
					ScriptName: tc.name,
				},
				Action: config.ShuttleAction{
					Shell: tc.script,
				},
			})

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
