package glcm

// hook implements the Hook interface.
type hook struct {
	f    func(...interface{}) error
	args []interface{}
	name string
}

// NewHook returns a new instance of the Hook.
func NewHook(name string, f func(...interface{}) error, args ...interface{}) Hook {
	return &hook{
		f:    f,
		args: args,
		name: name,
	}
}

// Execute executes the hook method.
func (h *hook) Execute() error {
	return h.f(h.args...)
}

// Name returns the name of the hook.
func (h *hook) Name() string {
	return h.name
}
