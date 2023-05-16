package cmder_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lunarway/shuttle/pkg/executors/golang/cmder"
)

func TestCmderWithError(t *testing.T) {
	testFunc := cmder.NewCmd("test", func(ctx context.Context) error {
		return errors.New("some-error")
	})

	args := []string{
		"test",
	}

	err := cmder.NewRoot().AddCmds(testFunc).TryExecute(args)

	assert.ErrorContains(t, err, "some-error")
}

func TestCmderWithNoError(t *testing.T) {
	testFunc := cmder.NewCmd("test", func(ctx context.Context) error {
		return nil
	})

	args := []string{
		"test",
	}

	err := cmder.NewRoot().AddCmds(testFunc).TryExecute(args)

	assert.NoError(t, err)
}

func TestCmderWithMultipeReturns(t *testing.T) {
	testFunc := cmder.NewCmd("test", func(ctx context.Context) (string, error) {
		return "something", nil
	})

	args := []string{
		"test",
	}

	err := cmder.NewRoot().AddCmds(testFunc).TryExecute(args)

	assert.NoError(t, err)
}

func TestCmderWithMultipeReturnsErroring(t *testing.T) {
	testFunc := cmder.NewCmd("test", func(ctx context.Context) (string, error) {
		return "something", errors.New("some-error")
	})

	args := []string{
		"test",
	}

	err := cmder.NewRoot().AddCmds(testFunc).TryExecute(args)

	assert.ErrorContains(t, err, "some-error")
}

func TestCmderWithMultipeReturnsErroringInAnyPlace(t *testing.T) {
	testFunc := cmder.NewCmd("test", func(ctx context.Context) (error, string) {
		return errors.New("some-error"), "something"
	})

	args := []string{
		"test",
	}

	err := cmder.NewRoot().AddCmds(testFunc).TryExecute(args)

	assert.ErrorContains(t, err, "some-error")
}
