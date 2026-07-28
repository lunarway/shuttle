package git

import (
	"io"
	"os"
	"testing"

	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestSkipPullFromEnv(t *testing.T) {
	tt := []struct {
		name      string
		skipPull  *string
		ci        *string
		skipsPull bool
	}{
		{
			name:      "nothing set",
			skipsPull: false,
		},
		{
			name:      "CI set",
			ci:        strPtr("true"),
			skipsPull: true,
		},
		{
			name:      "CI set to any non-empty value",
			ci:        strPtr("1"),
			skipsPull: true,
		},
		{
			name:      "CI set but empty",
			ci:        strPtr(""),
			skipsPull: false,
		},
		{
			name:      "SHUTTLE_SKIP_PULL set",
			skipPull:  strPtr("true"),
			skipsPull: true,
		},
		{
			name:      "SHUTTLE_SKIP_PULL set but empty is treated as set",
			skipPull:  strPtr(""),
			skipsPull: true,
		},
		{
			name:      "SHUTTLE_SKIP_PULL takes precedence over CI",
			skipPull:  strPtr("false"),
			ci:        strPtr("true"),
			skipsPull: false,
		},
		{
			name:      "SHUTTLE_SKIP_PULL with an unparsable value is treated as set",
			skipPull:  strPtr("yes-please"),
			skipsPull: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Unset by default so the developer's own environment does not leak
			// into the test.
			t.Setenv(skipPullKey, "")
			os.Unsetenv(skipPullKey)
			t.Setenv("CI", "")
			os.Unsetenv("CI")

			if tc.skipPull != nil {
				t.Setenv(skipPullKey, *tc.skipPull)
			}
			if tc.ci != nil {
				t.Setenv("CI", *tc.ci)
			}

			uii := ui.Create(io.Discard, io.Discard)

			assert.Equal(t, tc.skipsPull, skipPullFromEnv(uii))
		})
	}
}

func strPtr(s string) *string {
	return &s
}
