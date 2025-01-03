package service

import (
	"testing"

	"github.com/achu-1612/glcm/hook"
	"github.com/stretchr/testify/assert"
)

func TestWithPostHooks(t *testing.T) {
	// Create a mock post hook
	mockHook := &hook.MockHandler{}

	// Create a Wrapper instance
	wrapper := &Wrapper{}

	// Apply the WithPostHooks option
	option := WithPostHooks(mockHook)
	option(wrapper)

	// Assert that the postHooks field is set correctly
	assert.Equal(t, []hook.Handler{mockHook}, wrapper.postHooks, "Expected postHooks to be set correctly")
}

func TestWithPostHooks_MultipleHooks(t *testing.T) {
	// Create multiple mock post hooks
	mockHook1 := &hook.MockHandler{}
	mockHook2 := &hook.MockHandler{}

	// Create a Wrapper instance
	wrapper := &Wrapper{}

	// Apply the WithPostHooks option with multiple hooks
	option := WithPostHooks(mockHook1, mockHook2)
	option(wrapper)

	// Assert that the postHooks field is set correctly with multiple hooks
	assert.Equal(t, []hook.Handler{mockHook1, mockHook2}, wrapper.postHooks, "Expected postHooks to be set correctly with multiple hooks")
}

func TestWithPostHooks_NoHooks(t *testing.T) {
	// Create a Wrapper instance
	wrapper := &Wrapper{}

	// Apply the WithPostHooks option with no hooks
	option := WithPostHooks()
	option(wrapper)

	// Assert that the postHooks field is set correctly with no hooks
	assert.Equal(t, []hook.Handler(nil), wrapper.postHooks, "Expected postHooks to be set correctly with no hooks")
}

func TestWithPreHooks(t *testing.T) {
	// Create a mock pre hook
	mockHook := &hook.MockHandler{}

	// Create a Wrapper instance
	wrapper := &Wrapper{}

	// Apply the WithPreHooks option
	option := WithPreHooks(mockHook)
	option(wrapper)

	// Assert that the preHooks field is set correctly
	assert.Equal(t, []hook.Handler{mockHook}, wrapper.preHooks, "Expected preHooks to be set correctly")
}

func TestWithPreHooks_MultipleHooks(t *testing.T) {
	// Create multiple mock pre hooks
	mockHook1 := &hook.MockHandler{}
	mockHook2 := &hook.MockHandler{}

	// Create a Wrapper instance
	wrapper := &Wrapper{}

	// Apply the WithPreHooks option with multiple hooks
	option := WithPreHooks(mockHook1, mockHook2)
	option(wrapper)

	// Assert that the preHooks field is set correctly with multiple hooks
	assert.Equal(t, []hook.Handler{mockHook1, mockHook2}, wrapper.preHooks, "Expected preHooks to be set correctly with multiple hooks")
}

func TestWithPreHooks_NoHooks(t *testing.T) {
	// Create a Wrapper instance
	wrapper := &Wrapper{}

	// Apply the WithPreHooks option with no hooks
	option := WithPreHooks()
	option(wrapper)

	// Assert that the preHooks field is set correctly with no hooks
	assert.Equal(t, []hook.Handler(nil), wrapper.preHooks, "Expected preHooks to be set correctly with no hooks")
}
