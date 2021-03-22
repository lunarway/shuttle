package config_test

import (
	"errors"
	"testing"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestShuttleProjectContext_Documentation(t *testing.T) {
	tt := []struct {
		name    string
		planRef string
		docsRef string

		result string
		err    error
	}{
		{
			name:    "empty context",
			planRef: "",
			docsRef: "",
			result:  "",
			err:     errors.New("exit code 1 - Could not find any plan documentation"),
		},
		{
			name:    "unknown plan protocol",
			planRef: "something-odd",
			docsRef: "",
			result:  "",
			err:     errors.New("exit code 1 - Could not detect protocol for plan 'something-odd'"),
		},
		{
			name:    "unknown git plan protocol",
			planRef: "something-odd",
			docsRef: "",
			result:  "",
			err:     errors.New("exit code 1 - Could not detect protocol for plan 'something-odd'"),
		},
		{
			name:    "explicit HTTP docs",
			planRef: "",
			docsRef: "http://github.com/lunarway/shuttle",
			result:  "http://github.com/lunarway/shuttle",
			err:     nil,
		},
		{
			name:    "explicit HTTPS docs",
			planRef: "",
			docsRef: "https://github.com/lunarway/shuttle",
			result:  "https://github.com/lunarway/shuttle",
			err:     nil,
		},
		{
			name:    "no explicit docs and git plan ssh reference",
			planRef: "git://git@github.com:lunarway/shuttle-example-go-plan.git",
			docsRef: "",
			result:  "https://github.com/lunarway/shuttle-example-go-plan.git",
			err:     nil,
		},
		{
			name:    "no explicit docs and git plan http reference",
			planRef: "http://github.com/lunarway/shuttle-example-go-plan.git",
			docsRef: "",
			result:  "http://github.com/lunarway/shuttle-example-go-plan.git",
			err:     nil,
		},
		{
			name:    "no explicit docs and git plan https reference",
			planRef: "https://github.com/lunarway/shuttle-example-go-plan.git",
			docsRef: "",
			result:  "https://github.com/lunarway/shuttle-example-go-plan.git",
			err:     nil,
		},
		{
			name:    "no explicit docs and git plan has branch reference",
			planRef: "https://github.com/lunarway/shuttle-example-go-plan.git#branch",
			docsRef: "",
			result:  "https://github.com/lunarway/shuttle-example-go-plan.git",
			err:     nil,
		},
		{
			name:    "absolute local file path",
			planRef: "/plan",
			docsRef: "",
			result:  "",
			err:     errors.New("exit code 2 - Local plan has no documentation"),
		},
		{
			name:    "local file path",
			planRef: "./plan",
			docsRef: "",
			result:  "",
			err:     errors.New("exit code 2 - Local plan has no documentation"),
		},
		{
			name:    "local file path in parent path",
			planRef: "../plan",
			docsRef: "",
			result:  "",
			err:     errors.New("exit code 2 - Local plan has no documentation"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := config.ShuttleProjectContext{
				Config: config.ShuttleConfig{
					Plan: tc.planRef,
				},
				Plan: config.ShuttlePlanConfiguration{
					Documentation: tc.docsRef,
				},
			}
			result, err := p.DocumentationURL()

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error(), "error not as expected")
			} else {
				assert.NoError(t, err, "unexpected error")
			}
			assert.Equal(t, tc.result, result, "result not as expected")
		})
	}
}
