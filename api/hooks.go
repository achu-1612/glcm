package api

//go:generate mockgen -package api -destination hooks.mock.go -source hooks.go -self_package "github.com/achu-1612/glcm/api"

// Handler is an interface which represents a single hook.
// When a servcice is regsited in the runner, implementations of the Hndler interface can be registered too.
// Based on the nature of the hook (pre-run or post-run), the hook will be executed.
type Hook interface {
	// Execute executes the hook method.
	Execute() error

	// Name returns the name of the hook.
	Name()
}
