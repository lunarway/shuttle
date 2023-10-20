package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	variables := map[string]string{
		"VAR1": "TEST1",
		"VAR2": "TEST2",
		"VAR3": "TEST3",
	}

	for n, v := range variables {
		t.Setenv(n, v)
	}

	testCases := []testCase{
		{
			name:      "No exclude should display VAR1, VAR2 and VAR3 for Environment",
			input:     args("config"),
			stdoutput: "Version <dev-version>\nPlan:\nfalse\nEnvironment:\nVAR1=TEST1\nVAR2=TEST2\nVAR3=TEST3\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "with exlcude VAR2 and VAR3 should only display VAR1 for Environment",
			input:     args("config", "--exclude-env-vars", "VAR2,VAR3"),
			stdoutput: "Version <dev-version>\nPlan:\nfalse\nEnvironment:\nVAR1=TEST1\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "git plan",
			input:     args("-p", "testdata/project-git", "config"),
			stdoutput: "Version <dev-version>\nPlan:\nhttps://github.com/lunarway/shuttle-example-go-plan.git master\nEnvironment:\nVAR1=TEST1\nVAR2=TEST2\nVAR3=TEST3\n",
			erroutput: "",
			err:       nil,
		},
	}

	executeTestCasesWithCustomAssertion(
		t,
		testCases,
		func(t *testing.T, tc testCase, stdout, stderr string) {
			t.Helper()

			for _, outputLine := range strings.Split(tc.stdoutput, "\n") {
				assert.Contains(t, stdout, outputLine, "one std output not as expected")
			}

			assert.Equal(t, tc.erroutput, stderr, "err output not as expected")
		},
	)
}
