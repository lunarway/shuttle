package executer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsActionsEnabled(t *testing.T) {
	t.Run("default is enabled", func(t *testing.T) {
		t.Setenv("SHUTTLE_GOLANG_ACTIONS", "")

		assert.True(t, isActionsEnabled())
	})

	t.Run("set false is not enabled", func(t *testing.T) {
		t.Setenv("SHUTTLE_GOLANG_ACTIONS", "false")

		assert.False(t, isActionsEnabled())
	})

	t.Run("set true is enabled", func(t *testing.T) {
		t.Setenv("SHUTTLE_GOLANG_ACTIONS", "true")

		assert.True(t, isActionsEnabled())
	})

	t.Run("set any other value is enabled", func(t *testing.T) {
		t.Setenv("SHUTTLE_GOLANG_ACTIONS", "blabla")

		assert.True(t, isActionsEnabled())
	})
}
