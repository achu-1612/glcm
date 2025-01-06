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

func TestSupressLog(t *testing.T) {
	// Create a runner instance
	r := &runner{}

	// Apply the WithSupressLog option with true
	option := WithSupressLog(true)
	option(r)

	// Assert that the supressLog field is set correctly
	assert.Equal(t, true, r.supressLog, "Expected supressLog to be set to true")

	// Apply the WithSupressLog option with false
	option = WithSupressLog(false)
	option(r)

	// Assert that the supressLog field is set correctly
	assert.Equal(t, false, r.supressLog, "Expected supressLog to be set to false")
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
