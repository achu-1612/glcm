package service

import "sync"

type Context interface {
	// Done is to be called by the service when it completes and exits.
	// This will signal the runner that the service has completed.
	// If Done() is called then, the service will not be restarted after it stops.
	Done()

	// TerminationChan returns the termination channel.
	// The channel will be closed the service is to be stopped.
	// 1. The Runner is shutting down.
	// 2. The Stop() method is called on the service.
	TermCh() chan struct{}
}

// context holds all the lifecycle objects for the service.
type context struct {
	// PreHooks are the hooks that will be executed before starting the service.
	preHooks []func()

	// IgnorePreRunHooksError is a flag to indicate if the pre-run hooks error should be ignored or not.
	ignorePreRunHooksError bool

	// PostHooks are the hooks that will be executed after stopping the service.
	postHooks []func()

	// IgnorePostRunHooksError is a flag to indicate if the post-run hooks error should be ignored or not.
	ignorePostRunHooksError bool

	// terminationChan is a channel which will be used to direct the service to stop.
	terminationChan chan struct{}

	// wg is the wait group created by the base runner.
	wg *sync.WaitGroup
}

func (c *context) Done() {
	c.wg.Done()
}

func (c *context) TermCh() chan struct{} {
	return c.terminationChan
}

// Option defines a way to mutate the service configuration while registeration.
type Option func(opts *context)

// WithPreHooks sets the pre-hooks for the service.
func WithPreHooks(hooks ...func()) Option {
	return func(opts *context) {
		opts.preHooks = hooks
	}
}

// WithIgnorePreRunHooksError sets the ignorePreRunHooksError flag for the service.
func WithIgnorePreRunHooksError(ignore bool) Option {
	return func(opts *context) {
		opts.ignorePreRunHooksError = ignore
	}
}

// WithPostHooks sets the post-hooks for the service.
func WithPostHooks(hooks ...func()) Option {
	return func(opts *context) {
		opts.postHooks = hooks
	}
}

// WithIgnorePostRunHooksError sets the ignorePostRunHooksError flag for the service.
func WithIgnorePostRunHooksError(ignore bool) Option {
	return func(opts *context) {
		opts.ignorePostRunHooksError = ignore
	}
}
