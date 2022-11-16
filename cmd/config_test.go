package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	environmentStateBeforeTest := os.Environ()
	defer restoreEnvironment(environmentStateBeforeTest)

	os.Clearenv()
	os.Setenv("VAR1", "TEST1")
	os.Setenv("VAR2", "TEST2")
	os.Setenv("VAR3", "TEST3")
	testCases := []testCase{
		{
			name:      "No exlcude should display VAR1, VAR2 and VAR3 for Environment",
			input:     args("config"),
			stdoutput: "Version <dev-version>\n\nPlan: \n\nEnvironment: \nVAR1=TEST1\nVAR2=TEST2\nVAR3=TEST3\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "with exlcude VAR2 and VAR3 should only display VAR1 for Environment",
			input:     args("config", "--exclude-env-vars", "VAR2,VAR3"),
			stdoutput: "Version <dev-version>\n\nPlan: \n\nEnvironment: \nVAR1=TEST1\n",
			erroutput: "",
			err:       nil,
		},
	}

	executeTestCases(t, testCases)
}

func restoreEnvironment(originalState []string) {
	for _, envVar := range originalState {
		splitted := strings.Split(envVar, "=")
		key := splitted[0]
		val := splitted[1]

		os.Setenv(key, val)
	}
}
