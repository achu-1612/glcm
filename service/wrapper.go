package service

// Wrapper is a wrapper around the service and its context.
type Wrapper struct {
	s    Service
	sCtx *Context
}

// Service returns the service.
func (w *Wrapper) Service() Service {
	return w.s
}

// Context returns the context.
func (w *Wrapper) Context() *Context {
	return w.sCtx
}

// NewWrapper returns a new instance of the wrapper.
func NewWrapper(s Service, sCtx *Context) *Wrapper {
	return &Wrapper{
		s:    s,
		sCtx: sCtx,
	}
}
