package cmd

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func removeShuttleDirectories(t *testing.T, pwd string) {
	t.Helper()
	var directoriesToRemove []string
	err := fs.WalkDir(os.DirFS(pwd), "testdata", func(path string, d fs.DirEntry, err error) error {
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

func TestRun(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	t.Cleanup(func() {
		removeShuttleDirectories(t, pwd)
	})

	strings := func(s ...string) []string {
		return s
	}

	testCases := []testCase{
		{
			name:      "std out echo",
			input:     strings("-p", "testdata/project", "run", "hello_stdout"),
			stdoutput: "Hello stdout\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "std err echo",
			input:     strings("-p", "testdata/project", "run", "hello_stderr"),
			stdoutput: "",
			erroutput: "\x1b[31;1mHello stderr\x1b[0m\n",
			err:       nil,
		},
		{
			name:      "exit 0",
			input:     strings("-p", "testdata/project", "run", "exit_0"),
			stdoutput: "",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "exit 1",
			input:     strings("-p", "testdata/project", "run", "exit_1"),
			stdoutput: "",
			erroutput: "Error: exit code 4 - Failed executing script `exit_1`: shell script `exit 1`\nExit code: 1\n",
			err:       errors.New("exit code 4 - Failed executing script `exit_1`: shell script `exit 1`\nExit code: 1"),
		},
		{
			name:      "project with absolute path",
			input:     strings("-p", filepath.Join(pwd, "testdata/project"), "run", "hello_stdout"),
			stdoutput: "Hello stdout\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "project without shuttle.yaml",
			input:     strings("-p", "testdata/base", "run", "hello_stdout"),
			stdoutput: "",
			erroutput: `Error: exit code 2 - Failed to load shuttle configuration: shuttle.yaml file not found

Make sure you are in a project using shuttle and that a 'shuttle.yaml' file is available.
`,
			err: errors.New(`exit code 2 - Failed to load shuttle configuration: shuttle.yaml file not found

Make sure you are in a project using shuttle and that a 'shuttle.yaml' file is available.`),
		},
		{
			name:      "script fails when required argument is missing",
			input:     strings("-p", "testdata/project", "run", "required_arg"),
			stdoutput: "",
			erroutput: `Error: exit code 2 - Arguments not valid:
 'foo' not supplied but is required

Script 'required_arg' accepts the following arguments:
  foo (required)
`,
			err: errors.New(`exit code 2 - Arguments not valid:
 'foo' not supplied but is required

Script 'required_arg' accepts the following arguments:
  foo (required)`),
		},
		{
			name:      "script succeeds with required argument",
			input:     strings("-p", "testdata/project", "run", "required_arg", "foo=bar"),
			stdoutput: "bar\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "script succeeds with required argument missing and validation disabled",
			input:     strings("-p", "testdata/project", "run", "--validate=false", "required_arg"),
			stdoutput: "\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "script fails when validation is disabled and argument is not in valid format",
			input:     strings("-p", "testdata/project", "run", "--validate=false", "required_arg", "foo"),
			stdoutput: "",
			erroutput: `Error: exit code 2 - Arguments not valid:
 'foo' not <argument>=<value>

Script 'required_arg' accepts the following arguments:
  foo (required)
`,
			err: errors.New(`exit code 2 - Arguments not valid:
 'foo' not <argument>=<value>

Script 'required_arg' accepts the following arguments:
  foo (required)`),
		},
		{
			name:      "script fails on unkown argument",
			input:     strings("-p", "testdata/project", "run", "required_arg", "foo=bar", "a=b"),
			stdoutput: "",
			erroutput: `Error: exit code 2 - Arguments not valid:
 'a' unknown

Script 'required_arg' accepts the following arguments:
  foo (required)
`,
			err: errors.New(`exit code 2 - Arguments not valid:
 'a' unknown

Script 'required_arg' accepts the following arguments:
  foo (required)`),
		},
		{
			name:  "branched git plan",
			input: strings("-p", "testdata/project-git-branched", "run", "say"),
			stdoutput: `Cloning plan https://github.com/lunarway/shuttle-example-go-plan.git
something clever
`,
			erroutput: "",
			err:       nil,
		},
		{
			name:  "git plan",
			input: strings("-p", "testdata/project-git", "run", "say"),
			stdoutput: `Cloning plan https://github.com/lunarway/shuttle-example-go-plan.git
something masterly
`,
			erroutput: "",
			err:       nil,
		},
		{
			name:      "tagged git plan",
			input:     strings("-p", "testdata/project-git", "--plan", "#tagged", "run", "say"),
			stdoutput: "\x1b[032;1mOverload git plan branch/tag/sha with tagged\x1b[0m\n\x1b[032;1mSkipping plan pull because its running on detached head\x1b[0m\nsomething tagged\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "sha git plan",
			input:     strings("-p", "testdata/project-git", "--plan", "#2b52c21", "run", "say"),
			stdoutput: "\x1b[032;1mOverload git plan branch/tag/sha with 2b52c21\x1b[0m\n\x1b[032;1mSkipping plan pull because its running on detached head\x1b[0m\nsomething minor\n",
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
