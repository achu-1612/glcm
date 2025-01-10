package runner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsRunning(t *testing.T) {
	r := NewRunner()

	// Test when runner is not running
	assert.False(t, r.IsRunning(), "Expected runner to not be running")

	// Start the runner
	go func() {
		if err := r.BootUp(context.Background()); err != nil {
			t.Errorf("Error while booting up the runner: %v", err)
		}
	}()

	<-time.After(time.Second * 10)

	// Test when runner is running
	assert.True(t, r.IsRunning(), "Expected runner to be running")

	// Shutdown the runner
	r.Shutdown()

	// Test when runner is not running after shutdown
	assert.False(t, r.IsRunning(), "Expected runner to not be running after shutdown")
}
