// pakcage hook allows user to define hooks which can be executed before or after the service is started.
package hook

//go:generate mockgen -package hook -destination hook.mock.go -source hook.go -self_package "github.com/achu-1612/glcm/hook"

// Handler is an interface which represents a single hook.
// When a servcice is regsited in the runner, implementations of the Hndler interface can be registered too.
// Based on the nature of the hook (pre-run or post-run), the hook will be executed.
type Handler interface {
	// Execute executes the hook method.
	Execute() error

	// Name returns the name of the hook.
	Name() string
}
