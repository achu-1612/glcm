package glcm

// hookHandler implements the Handler interface.
type hookHandler struct {
	f    func(...interface{}) error
	args []interface{}
	name string
}

// NewHook returns a new instance of the Hook.
func NewHook(name string, f func(...interface{}) error, args ...interface{}) Hook {
	return &hookHandler{
		f:    f,
		args: args,
		name: name,
	}
}

// Execute executes the hook method.
func (h *hookHandler) Execute() error {
	return h.f(h.args...)
}

// Name returns the name of the hook.
func (h *hookHandler) Name() string {
	return h.name
}
