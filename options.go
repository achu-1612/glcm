package glcm

import (
	"time"
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
}
