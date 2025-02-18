package glcm

import (
	"time"

	"github.com/achu-1612/glcm/log"
)

const (
	defaultSocketPath      = "/tmp/glcm.sock"
	defaultShutdownTimeout = time.Second * 30
	defaultMaxRetries      = 10
	defaultBackoffExp      = 2
)

// ServiceOptions represents the options for a service.
type ServiceOptions struct {
	// PreHooks are the hooks that are executed before the service is started.
	PreHooks []Hook

	// PostHooks are the hooks that are executed after the service is stopped.
	PostHooks []Hook

	// AutoStart represents the options for auto-starting the service.
	AutoStart AutoRestartOptions

	// Schedule represents the options for scheduling the service.
	Schedule SchedulingOptions
}

// Sanitize fills the default values for the service options.
func (s *ServiceOptions) Sanitize() {
	if s.AutoStart.MaxRetries == 0 {
		log.Warnf("MaxRetries is not set for service. Setting it to default value %d", defaultMaxRetries)

		s.AutoStart.MaxRetries = defaultMaxRetries
	}

	if s.AutoStart.BackOffExponent == 0 {
		log.Warnf("BackoffExponent is not set for service. Setting it to default value %d", defaultBackoffExp)

		s.AutoStart.BackOffExponent = defaultBackoffExp
	}
}

// AutoRestartOptions represents the options for auto-restarting the service.
type AutoRestartOptions struct {
	// Enabled represents if the auto-restart is enabled.
	Enabled bool

	// MaxRetries represents the maximum number of retries.
	MaxRetries int

	// Backoff represents if the backoff is enabled.
	Backoff bool

	// BackOffExponent represents the exponent for the backoff.
	BackOffExponent int
}

// SchedulingOptions represents the options for scheduling the service.
type SchedulingOptions struct {
	// Enabled represents if the auto-restart is enabled.
	Enabled bool

	// Cron represents the cron expression for scheduling the service.
	Cron string

	// TimeOut represents the timeout for the service.
	// After the timeout, the service will be sent a termination signal.
	TimeOut time.Duration

	// MaxRuns represents the maximum number of runs for the service.
	MaxRuns int
}

// RunnerOptions represents the options for the runner.
type RunnerOptions struct {
	// HideBanner represents if the banner should be hidden.
	HideBanner bool

	// Verbose represents if the logs should be suppressed or not.
	Verbose bool

	// Socket represents if the socket should be enabled or not.
	Socket bool

	// SocketPath represents the path to the socket file.
	SocketPath string

	// AllowedUID represents the allowed user ids to interact with the socket.
	AllowedUID []int

	// ShutdownTimeout represents the timeout for shutting down the runner.
	ShutdownTimeout time.Duration
}

// Santizie fills the default values for the runner options.
func (r *RunnerOptions) Sanitize() {
	if r.ShutdownTimeout == 0 {
		log.Warnf("ShutdownTimeout is not set for runner. Setting it to default value %v", defaultShutdownTimeout)

		r.ShutdownTimeout = defaultShutdownTimeout
	}

	if r.SocketPath == "" {
		log.Warnf("SocketPath is not set for runner. Setting it to default value %s", defaultSocketPath)

		r.SocketPath = defaultSocketPath
	}
}
