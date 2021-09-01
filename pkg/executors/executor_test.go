package executors

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/lunarway/shuttle/pkg/config"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestValidateUnknownArgs(t *testing.T) {
	tt := []struct {
		name       string
		scriptArgs []config.ShuttleScriptArgs
		inputArgs  map[string]string
		output     []validationError
	}{
		{
			name:       "no input or script",
			scriptArgs: nil,
			inputArgs:  nil,
			output:     nil,
		},
		{
			name:       "single input without script",
			scriptArgs: nil,
			inputArgs: map[string]string{
				"foo": "1",
			},
			output: []validationError{
				{"foo", "unknown"},
			},
		},
		{
			name:       "multiple input without script",
			scriptArgs: nil,
			inputArgs: map[string]string{
				"foo": "1",
				"bar": "2",
			},
			output: []validationError{
				{"bar", "unknown"},
				{"foo", "unknown"},
			},
		},
		{
			name: "single input and script",
			scriptArgs: []config.ShuttleScriptArgs{
				{
					Name: "foo",
				},
			},
			inputArgs: map[string]string{
				"foo": "1",
			},
			output: nil,
		},
		{
			name: "multple input and script",
			scriptArgs: []config.ShuttleScriptArgs{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			},
			inputArgs: map[string]string{
				"bar": "2",
				"foo": "1",
			},
			output: nil,
		},
		{
			name: "multple input and script with one unknown",
			scriptArgs: []config.ShuttleScriptArgs{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			},
			inputArgs: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			output: []validationError{
				{"baz", "unknown"},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := validateUnknownArgs(tc.scriptArgs, tc.inputArgs)
			// sort as the order is not guarenteed by validateUnknownArgs
			sortValidationErrors(output)
			assert.Equal(t, tc.output, output, "output not as expected")
		})
	}
}

func TestSortValidationErrors(t *testing.T) {
	tt := []struct {
		name   string
		input  []validationError
		output []validationError
	}{
		{
			name: "sorted",
			input: []validationError{
				{"bar", ""},
				{"baz", ""},
				{"foo", ""},
			},
			output: []validationError{
				{"bar", ""},
				{"baz", ""},
				{"foo", ""},
			},
		},
		{
			name: "not sorted",
			input: []validationError{
				{"baz", ""},
				{"foo", ""},
				{"bar", ""},
			},
			output: []validationError{
				{"bar", ""},
				{"baz", ""},
				{"foo", ""},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sortValidationErrors(tc.input)
			assert.Equal(t, tc.output, tc.input, "output not as expected")
		})
	}
}

func TestExecute(t *testing.T) {
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
			verboseUI := ui.Create()
			verboseUI.SetUserLevel(ui.LevelVerbose)

			err := Execute(context.Background(), config.ShuttleProjectContext{
				ProjectPath: ".",
				UI:          verboseUI,
				Scripts: map[string]config.ShuttlePlanScript{
					"test": {
						Actions: []config.ShuttleAction{
							{
								Shell: tc.script,
							},
						},
					},
				},
			}, "test", nil, true)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestExecute_contextCancellation tests that scripts are closed when the
// context is cancelled.
func TestExecute_contextCancellation(t *testing.T) {
	imageName := fmt.Sprintf("shuttle-test-execute-cancellation-%d", time.Now().UnixNano())
	t.Logf("Starting image %s", imageName)
	verboseUI := ui.Create()
	verboseUI.SetUserLevel(ui.LevelVerbose)
	projectContext := config.ShuttleProjectContext{
		UI: verboseUI,
		Scripts: map[string]config.ShuttlePlanScript{
			"serve": {
				Description: "",
				Actions: []config.ShuttleAction{
					{
						Shell: fmt.Sprintf("docker run --rm -i --name %s nginx", imageName),
					},
				},
				Args: []config.ShuttleScriptArgs{},
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// let the container start before stopping it
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	err := Execute(ctx, projectContext, "serve", nil, true)
	assert.EqualError(t, err, context.Canceled.Error())

	// sadly we need to give the docker some time before "docker ps" shows the
	// containers
	time.Sleep(500 * time.Millisecond)
	images := runningDockerImages(t, imageName)
	assert.Len(t, images, 0, "expected no images to be running")
}

func runningDockerImages(t *testing.T, imageName string) []string {
	t.Helper()
	cmd := cmd.NewCmd("docker", "ps", "-a", "--format", "{{ .Names }}")
	status := <-cmd.Start()
	t.Logf("docker ps: stderr: %v", status.Stderr)

	t.Logf("Docker containers")
	for _, container := range status.Stdout {
		t.Logf("- %s", container)
		if container == imageName {
			t.Errorf("Container '%s still exists in docker but shouldn't", container)
		}
	}
	if status.Exit != 0 {
		t.Fatal("Failed to check running docker images")
	}
	return nil
}
