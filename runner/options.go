package runner

import "context"

type Options func(opts *runner)

// WithHideBanner sets the hide banner flag for the runner.
func WithHideBanner(hideBanner bool) Options {
	return func(opts *runner) {
		opts.hideBanner = hideBanner
	}
}

// WithSupressLog sets the supress log flag for the runner.
func WithSupressLog(supressLog bool) Options {
	return func(opts *runner) {
		opts.supressLog = supressLog
	}
}

// WithContext sets the context for the runner.
func WithContext(ctx context.Context) Options {
	return func(opts *runner) {
		opts.ctx = ctx
	}
}

