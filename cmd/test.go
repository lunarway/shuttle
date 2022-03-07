package cmd

import (
	"bytes"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func args(s ...string) []string {
	return s
}

type testCase struct {
	name      string
	input     []string
	stdoutput string
	erroutput string
	err       error
}

func executeTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdBuf := new(bytes.Buffer)
			errBuf := new(bytes.Buffer)

			rootCmd, _ := initializedRoot(stdBuf, errBuf)
			rootCmd.SetArgs(tc.input)

			err := rootCmd.Execute()
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
			t.Logf("STDOUT: %s", stdBuf.String())
			t.Logf("STDERR: %s", errBuf.String())
			assert.Equal(t, tc.stdoutput, stdBuf.String(), "std output not as expected")
			assert.Equal(t, tc.erroutput, errBuf.String(), "err output not as expected")
		})
	}
}

func removeShuttleDirectories(t *testing.T) {
	t.Helper()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	var directoriesToRemove []string
	err = fs.WalkDir(os.DirFS(pwd), "testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".shuttle" {
			directoriesToRemove = append(directoriesToRemove, path)
		}
		return nil
	})
	if err != nil {
		t.Errorf("Failed to cleanup .shuttle files: %v", err)
	}

	for _, d := range directoriesToRemove {
		err := os.RemoveAll(d)
		if err != nil {
			t.Errorf("Failed to cleanup '%s': %v", d, err)
		}
	}
}
