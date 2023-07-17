package telemetry

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	t.Run("run without inherititng context_id", func(t *testing.T) {
		ctx := context.Background()

		t.Setenv("SHUTTLE_CONTEXT_ID", "")

		ctx = WithContextID(ctx)

		value := ctx.Value(telemetryContextID).(string)

		assert.NotEmpty(t, value)
	})

	t.Run("run inherititng context_id", func(t *testing.T) {
		ctx := context.Background()

		expected := uuid.New().String()
		t.Setenv("SHUTTLE_CONTEXT_ID", expected)

		ctx = WithContextID(ctx)

		value := ctx.Value(telemetryContextID).(string)

		assert.Equal(t, expected, value)
	})

}
