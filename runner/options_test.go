package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHideBanner(t *testing.T) {
	// Create a runner instance
	r := &runner{}

	// Apply the WithHideBanner option with true
	option := WithHideBanner(true)
	option(r)

	// Assert that the hideBanner field is set correctly
	assert.Equal(t, true, r.hideBanner, "Expected hideBanner to be set to true")

	// Apply the WithHideBanner option with false
	option = WithHideBanner(false)
	option(r)

	// Assert that the hideBanner field is set correctly
	assert.Equal(t, false, r.hideBanner, "Expected hideBanner to be set to false")
}

func TestSuppressLog(t *testing.T) {
	// Create a runner instance
	r := &runner{}

	// Apply the WithSuppressLog option with true
	option := WithSuppressLog(true)
	option(r)

	// Assert that the suppressLog field is set correctly
	assert.Equal(t, true, r.suppressLog, "Expected suppressLog to be set to true")

	// Apply the WithSuppressLog option with false
	option = WithSuppressLog(false)
	option(r)

	// Assert that the suppressLog field is set correctly
	assert.Equal(t, false, r.suppressLog, "Expected suppressLog to be set to false")
}

func TestContext(t *testing.T) {
	// Create a runner instance
	r := &runner{}

	// Create a context
	ctx := context.Background()

	// Apply the WithContext option
	option := WithContext(ctx)
	option(r)

	// Assert that the ctx field is set correctly
	assert.Equal(t, ctx, r.ctx, "Expected ctx to be set correctly")
}
