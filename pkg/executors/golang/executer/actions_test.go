package executer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionsMerge(t *testing.T) {
	t.Parallel()

	t.Run("can merge multiple other actions", func(t *testing.T) {
		sut := NewActions()

		actual := sut.Merge(
			&Actions{
				Actions: map[string]Action{
					"someAction": {},
				},
			},
			&Actions{
				Actions: map[string]Action{
					"someOtherAction": {},
				},
			},
		)

		assert.Len(t, actual.Actions, 2)
	})

	t.Run("can override", func(t *testing.T) {
		sut := NewActions()

		sut.Actions["someAction"] = Action{
			Args: []ActionArg{
				{
					Name: "someArg",
				},
			},
		}

		actual := sut.Merge(
			&Actions{
				Actions: map[string]Action{
					"someAction": {},
				},
			},
		)

		assert.Len(t, actual.Actions["someAction"].Args, 0)
	})

	t.Run("ignores nil", func(t *testing.T) {
		sut := NewActions()

		actual := sut.Merge(
			nil,
			nil,
		)

		assert.Len(t, actual.Actions, 0)
	})

}

func TestActionsExecute(t *testing.T) {
	t.Parallel()

	t.Run("finds action executes closure", func(t *testing.T) {
		sut := NewActions().Merge(&Actions{
			Actions: map[string]Action{
				"action": {},
			},
		})

		called := false

		ran, err := sut.Execute("action", func() error {
			called = true
			return nil
		})
		assert.NoError(t, err)

		assert.True(t, ran)
		assert.True(t, called)
	})

	t.Run("does not find an action does not execute", func(t *testing.T) {
		sut := NewActions().Merge(&Actions{
			Actions: map[string]Action{},
		})

		called := false

		ran, err := sut.Execute("action", func() error {
			called = true
			return nil
		})
		assert.NoError(t, err)

		assert.False(t, ran)
		assert.False(t, called)
	})

	t.Run("action is null", func(t *testing.T) {
		var action *Actions
		action = nil

		ran, err := action.Execute("something", func() error { return nil })

		assert.False(t, ran)
		assert.NoError(t, err)
	})
}
