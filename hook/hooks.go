// pakcage hook allows user to define hooks which can be executed before or after the service is started.
package hook

import "github.com/achu-1612/glcm/log"

// Handler is an interface which represents a single hook.
// When a servcice is regsited in the runner, implementations of the Hndler interface can be registered too.
// Based on the nature of the hook (pre-run or post-run), the hook will be executed.
type Handler interface {
	// Execute executes the hook method.
	Execute() error
}

// handler implements the Handler interface.
type handler struct {
	f    func(...interface{}) error
	args []interface{}
	name string
}

// NewHandler returns a new instance of the handler.
func NewHandler(name string, f func(...interface{}) error, args ...interface{}) Handler {
	return &handler{
		f:    f,
		args: args,
		name: name,
	}
}

// Execute executes the hook method.
func (h *handler) Execute() error {
	log.Infof("Executing hook: %s", h.name)

	err := h.f(h.args...)

	log.Infof("Excuted hook: %s", h.name)

	return err
}
