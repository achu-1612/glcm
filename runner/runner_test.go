package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRunning(t *testing.T) {
	r := NewRunner()

	// Test when runner is not running
	assert.False(t, r.IsRunning(), "Expected runner to not be running")

	// Start the runner
	r.BootUp(context.TODO())

	// Test when runner is running
	assert.True(t, r.IsRunning(), "Expected runner to be running")

	// Shutdown the runner
	r.Shutdown()

	// Test when runner is not running after shutdown
	assert.False(t, r.IsRunning(), "Expected runner to not be running after shutdown")
}
