package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGitPlan(t *testing.T) {
	tt := []struct {
		input  string
		output gitPlan
	}{
		{
			input: "git://git@github.com:lunarway/some-plan.git#some-branch",
			output: gitPlan{
				isGitPlan:  true,
				protocol:   "ssh",
				user:       "git",
				host:       "github.com",
				repository: "github.com:lunarway/some-plan.git",
				head:       "some-branch",
			},
		},
		{
			input: "git://git@github.com:lunarway/some-plan.git",
			output: gitPlan{
				isGitPlan:  true,
				protocol:   "ssh",
				user:       "git",
				host:       "github.com",
				repository: "github.com:lunarway/some-plan.git",
				head:       "master",
			},
		},
		{
			input: "https://github.com/lunarway/some-plan.git#some-branch",
			output: gitPlan{
				isGitPlan:  true,
				protocol:   "https",
				user:       "",
				host:       "",
				repository: "github.com/lunarway/some-plan.git",
				head:       "some-branch",
			},
		},
		{
			input: "https://github.com/lunarway/some-plan.git",
			output: gitPlan{
				isGitPlan:  true,
				protocol:   "https",
				user:       "",
				host:       "",
				repository: "github.com/lunarway/some-plan.git",
				head:       "master",
			},
		},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("can parse %s", tc.input), func(t *testing.T) {
			output := parseGitPlan(tc.input)

			assert.Equal(t, tc.output, output, "output does not match the expected")
		})
	}
}
