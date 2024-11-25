package service

// ServiceOptions defines the options that can be passed while registering the service.
type serviceOptions struct {
	// PreHooks are the hooks that will be executed before starting the service.
	preHooks []func()

	// IgnorePreRunHooksError is a flag to indicate if the pre-run hooks error should be ignored or not.
	ignorePreRunHooksError bool

	// PostHooks are the hooks that will be executed after stopping the service.
	postHooks []func()

	// IgnorePostRunHooksError is a flag to indicate if the post-run hooks error should be ignored or not.
	ignorePostRunHooksError bool
}

// Option defines a way to mutate the service configuration while registeration.
type Option func(opts *serviceOptions)

// WithPreHooks sets the pre-hooks for the service.
func WithPreHooks(hooks ...func()) Option {
	return func(opts *serviceOptions) {
		opts.preHooks = hooks
	}
}

// WithIgnorePreRunHooksError sets the ignorePreRunHooksError flag for the service.
func WithIgnorePreRunHooksError(ignore bool) Option {
	return func(opts *serviceOptions) {
		opts.ignorePreRunHooksError = ignore
	}
}

// WithPostHooks sets the post-hooks for the service.
func WithPostHooks(hooks ...func()) Option {
	return func(opts *serviceOptions) {
		opts.postHooks = hooks
	}
}

// WithIgnorePostRunHooksError sets the ignorePostRunHooksError flag for the service.
func WithIgnorePostRunHooksError(ignore bool) Option {
	return func(opts *serviceOptions) {
		opts.ignorePostRunHooksError = ignore
	}
}
