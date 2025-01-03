package hook

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
	return h.f(h.args...)
}

// Name returns the name of the hook.
func (h *handler) Name() string {
	return h.name
}
