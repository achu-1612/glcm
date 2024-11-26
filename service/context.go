package service

import "sync"

// context holds all the lifecycle objects for the service.
type Context struct {
	// PreHooks are the hooks that will be executed before starting the service.
	preHooks []func()

	// IgnorePreRunHooksError is a flag to indicate if the pre-run hooks error should be ignored or not.
	ignorePreRunHooksError bool

	// PostHooks are the hooks that will be executed after stopping the service.
	postHooks []func()

	// IgnorePostRunHooksError is a flag to indicate if the post-run hooks error should be ignored or not.
	ignorePostRunHooksError bool

	// terminationChan is a channel which will be used to direct the service to stop.
	// The channel will be closed the service is to be stopped.
	// 1. The Runner is shutting down.
	// 2. The Stop() method is called on the service.
	terminationChan chan struct{}

	// wg is the wait group created by the base runner.
	wg *sync.WaitGroup
}

// PreHooks returns the pre-hooks for the service.
func (c *Context) PreHooks() []func() {
	return c.preHooks
}

// IgnorePreRunHooksError returns the ignorePreRunHooksError flag for the service.
func (c *Context) IgnorePreRunHooksError() bool {
	return c.ignorePreRunHooksError
}

// PostHooks returns the post-hooks for the service.
func (c *Context) PostHooks() []func() {
	return c.postHooks
}

// IgnorePostRunHooksError returns the ignorePostRunHooksError flag for the service.
func (c *Context) IgnorePostRunHooksError() bool {
	return c.ignorePostRunHooksError
}

func (c *Context) Done() {
	c.wg.Done()
}

func (c *Context) TermCh() chan struct{} {
	return c.terminationChan
}

// Option defines a way to mutate the service configuration while registeration.
type Option func(opts *Context)

// WithPreHooks sets the pre-hooks for the service.
func WithPreHooks(hooks ...func()) Option {
	return func(opts *Context) {
		opts.preHooks = hooks
	}
}

// WithIgnorePreRunHooksError sets the ignorePreRunHooksError flag for the service.
func WithIgnorePreRunHooksError(ignore bool) Option {
	return func(opts *Context) {
		opts.ignorePreRunHooksError = ignore
	}
}

// WithPostHooks sets the post-hooks for the service.
func WithPostHooks(hooks ...func()) Option {
	return func(opts *Context) {
		opts.postHooks = hooks
	}
}

// WithIgnorePostRunHooksError sets the ignorePostRunHooksError flag for the service.
func WithIgnorePostRunHooksError(ignore bool) Option {
	return func(opts *Context) {
		opts.ignorePostRunHooksError = ignore
	}
}
