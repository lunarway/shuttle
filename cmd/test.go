package cmd

import (
	"bytes"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func args(s ...string) []string {
	return s
}

type testCase struct {
	name      string
	input     []string
	initErr   error
	stdoutput string
	erroutput string
	err       error
}

func executeTestCasesWithCustomAssertion(
	t *testing.T,
	testCases []testCase,
	assertion func(t *testing.T, tc testCase, stdout, stderr string),
) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// remove any .shuttle files up front and after each test to make sure the
			// runs are deterministic
			t.Cleanup(func() {
				removeShuttleDirectories(t)
			})
			removeShuttleDirectories(t)

			stdBuf := new(bytes.Buffer)
			errBuf := new(bytes.Buffer)

			rootCmd, _, err := initializedRootFromArgs(stdBuf, errBuf, tc.input)
			if err != nil {
				require.Error(t, err)
				require.Error(t, tc.initErr)
				require.Equal(t, tc.initErr.Error(), err.Error())
				return
			}
			rootCmd.SetArgs(tc.input)

			err = rootCmd.Execute()
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
			assertion(t, tc, stdBuf.String(), errBuf.String())
		})
	}
}

func executeTestCases(t *testing.T, testCases []testCase) {
	executeTestCasesWithCustomAssertion(
		t,
		testCases,
		func(t *testing.T, tc testCase, stdout, stderr string) {
			assert.Equal(t, tc.stdoutput, stdout, "std output not as expected")
			assert.Equal(t, tc.erroutput, stderr, "err output not as expected")
		},
	)
}

func executeTestContainsCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// remove any .shuttle files up front and after each test to make sure the
			// runs are deterministic
			t.Cleanup(func() {
				removeShuttleDirectories(t)
			})
			removeShuttleDirectories(t)

			stdBuf := new(bytes.Buffer)
			errBuf := new(bytes.Buffer)
			rootCmd, _, err := initializedRootFromArgs(stdBuf, errBuf, tc.input)
			if err != nil {
				if tc.initErr == nil {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
					require.Error(t, tc.initErr)
					require.Equal(t, tc.initErr.Error(), err.Error())
					return
				}
			}
			rootCmd.SetArgs(tc.input)

			err = rootCmd.Execute()
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err.Error())
			}
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
