package cmd

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	variables := map[string]string{
		"VAR1": "TEST1",
		"VAR2": "TEST2",
		"VAR3": "TEST3",
	}
	t.Cleanup(func() {
		for n := range variables {
			os.Unsetenv(n)
		}
	})
	for n, v := range variables {
		os.Setenv(n, v)
	}
	testCases := []testCase{
		{
			name:      "No exlcude should display VAR1, VAR2 and VAR3 for Environment",
			input:     args("config"),
			stdoutput: "(?s)^Version <dev-version>\n\nPlan:\nfalse\n\nEnvironment:\n.*VAR1=TEST1\nVAR2=TEST2\nVAR3=TEST3\n$",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "with exlcude VAR2 and VAR3 should only display VAR1 for Environment",
			input:     args("config", "--exclude-env-vars", "VAR2,VAR3"),
			stdoutput: "(?s)^Version <dev-version>\n\nPlan:\nfalse\n\nEnvironment:\n.*VAR1=TEST1\n$",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "git plan",
			input:     args("-p", "testdata/project-git", "config"),
			stdoutput: "(?s)^Version <dev-version>\n\nPlan:\nhttps://github.com/lunarway/shuttle-example-go-plan.git master\n\nEnvironment:\n.*VAR1=TEST1\nVAR2=TEST2\nVAR3=TEST3\n",
			erroutput: "",
			err:       nil,
		},
	}

	executeTestCasesWithCustomAssertion(t, testCases, func(t *testing.T, tc testCase, stdout, stderr string) {
		t.Helper()
		assert.Regexp(t, regexp.MustCompile(tc.stdoutput), stdout, "std output not as expected")
		assert.Equal(t, tc.erroutput, stderr, "err output not as expected")
	})
}
