package cmd

import (
	"testing"
)

func TestPlan(t *testing.T) {
	testCases := []testCase{
		{
			name:      "no plan",
			input:     args("-p", "testdata/project", "plan"),
			stdoutput: "",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "git plan",
			input:     args("-p", "testdata/project-git", "plan"),
			stdoutput: "https://github.com/lunarway/shuttle-example-go-plan.git",
			erroutput: "Cloning plan https://github.com/lunarway/shuttle-example-go-plan.git\n",
			err:       nil,
		},
		{
			name:      "no plan with template",
			input:     args("-p", "testdata/project", "plan", "--template", "{{.PlanRaw}}"),
			stdoutput: "false",
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
