package shuttlefile_test

import (
	"context"
	"testing"

	"github.com/kjuulh/shuttletask/pkg/shuttlefile"
	shuttleconfig "github.com/lunarway/shuttle/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {
	filePath := "testdata/shuttle.yaml"

	config, err := shuttlefile.ParseFile(context.Background(), filePath)
	assert.NoError(t, err)
	assert.Equal(t, shuttleconfig.ShuttleConfig{
		Plan:    "someplan",
		PlanRaw: "someplan",
		Variables: map[string]interface{}{
			"someVar": "var",
			"someNestedVar": map[string]interface{}{
				"nestedVar": true,
			},
		},
		Scripts: map[string]shuttleconfig.ShuttlePlanScript{
			"someAction": {
				Description: "",
				Actions: []shuttleconfig.ShuttleAction{
					{
						Shell: "action",
					},
				},
			},
		},
	}, *config)
}
