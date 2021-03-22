package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePlan(t *testing.T) {
	tt := []struct {
		input  string
		output Plan
	}{
		{
			input: "git://git@github.com:lunarway/some-plan.git#some-branch",
			output: Plan{
				IsGitPlan:  true,
				Protocol:   "ssh",
				User:       "git",
				Host:       "github.com",
				Repository: "github.com:lunarway/some-plan.git",
				Head:       "some-branch",
			},
		},
		{
			input: "git://git@github.com:lunarway/some-plan.git",
			output: Plan{
				IsGitPlan:  true,
				Protocol:   "ssh",
				User:       "git",
				Host:       "github.com",
				Repository: "github.com:lunarway/some-plan.git",
				Head:       "master",
			},
		},
		{
			input: "https://github.com/lunarway/some-plan.git#some-branch",
			output: Plan{
				IsGitPlan:  true,
				Protocol:   "https",
				User:       "",
				Host:       "",
				Repository: "github.com/lunarway/some-plan.git",
				Head:       "some-branch",
			},
		},
		{
			input: "https://github.com/lunarway/some-plan.git",
			output: Plan{
				IsGitPlan:  true,
				Protocol:   "https",
				User:       "",
				Host:       "",
				Repository: "github.com/lunarway/some-plan.git",
				Head:       "master",
			},
		},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("can parse %s", tc.input), func(t *testing.T) {
			output := ParsePlan(tc.input)

			assert.Equal(t, tc.output, output, "output does not match the expected")
		})
	}
}
