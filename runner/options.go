package runner

import "context"

type Options func(opts *runner)

// WithHideBanner sets the hide banner flag for the runner.
func WithHideBanner(hideBanner bool) Options {
	return func(opts *runner) {
		opts.hideBanner = hideBanner
	}
}

// WithSuppressLog sets the suppress log flag for the runner.
func WithSuppressLog(suppressLog bool) Options {
	return func(opts *runner) {
		opts.suppressLog = suppressLog
	}
}

// WithContext sets the context for the runner.
func WithContext(ctx context.Context) Options {
	return func(opts *runner) {
		opts.ctx = ctx
	}
}
