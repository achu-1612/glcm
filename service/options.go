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
