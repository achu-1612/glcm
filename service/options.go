package service

import (
	"github.com/achu-1612/glcm/hook"
)

// Option defines a way to mutate the service configuration while registeration.
type Option func(opts *Wrapper)

// WithPreHooks sets the pre-hooks for the service.
func WithPreHooks(hooks ...hook.Handler) Option {
	return func(opts *Wrapper) {
		opts.preHooks = hooks
	}
}

// WithPostHooks sets the post-hooks for the service.
func WithPostHooks(hooks ...hook.Handler) Option {
	return func(opts *Wrapper) {
		opts.postHooks = hooks
	}
}

// WithAutoRestart sets the auto-restart option for the service.
// AutoRestart will only happen if the service is exited and
// not stopped by the base runner as a result of runner shutdown or stop call..
func WithAutoRestart() Option {
	return func(opts *Wrapper) {
		opts.autoRestart = true
	}
}

// WithBackoff sets the backoff option for the service.
func WithBackoff() Option {
	return func(opts *Wrapper) {
		opts.backoff = true
	}
}

// WithMaxRetries sets the maximum number of retries for the service.
func WithMaxRetries(maxRetries int) Option {
	return func(opts *Wrapper) {
		opts.maxRetries = maxRetries
	}
}
