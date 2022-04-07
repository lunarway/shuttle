package cmd

import (
	"errors"
	"testing"
)

func TestLs(t *testing.T) {
	testCases := []testCase{
		{
			name:      "invalid shuttle.yaml file",
			input:     args("-p", "testdata/invalid-yaml", "ls"),
			stdoutput: "",
			erroutput: "Error: exit code 2 - Failed to parse shuttle configuration: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `plan` into config.ShuttleConfig\n\nMake sure your 'shuttle.yaml' is valid.\n",
			err:       errors.New("exit code 2 - Failed to parse shuttle configuration: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `plan` into config.ShuttleConfig\n\nMake sure your 'shuttle.yaml' is valid."),
		},
		{
			name:      "list one action",
			input:     args("-p", "../examples/no-plan-project", "ls"),
			stdoutput: "Available Scripts:\n  hello        \n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "list actions",
			input:     args("-p", "../examples/repo-project/", "ls"),
			stdoutput: "Pulling latest plan changes on master\nAvailable Scripts:\n  build        Build the docker image\n  deploy       Deploys the image to a kubernetes environment\n  push         Push the docker image\n  say          Say something\n  test         Run test for the project\n",
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
