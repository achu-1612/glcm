package glcm

import "github.com/achu-1612/glcm/runner"

// NewRunner returns a new base runner for glcm
func NewRunner(opts ...runner.Options) runner.Base {
	return runner.NewRunner(opts...)
}
