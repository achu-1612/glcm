package service

import (
	"sync"

	"github.com/achu-1612/glcm/log"
)

// Wrapper is a wrapper around the service and its context.
type Wrapper struct {
	s    Service
	sCtx *Context
}

// Service returns the service.
func (w *Wrapper) Service() Service {
	return w.s
}

// context returns the context.
func (w *Wrapper) Context() *Context {
	return w.sCtx
}

// reallocate the chan before starting if it is nil
func (w *Wrapper) Start() {
	if w.sCtx.isRunning() {
		log.Infof("Service %s is already running", w.s.Name())

		return
	}

	preHookError := false

	w.sCtx.wg.Add(1)
	defer func() {
		w.Context().Done() // indicate the worker group that the service has stopped.
		log.Infof("service %s stopped", w.Service().Name())
	}()

	func() {
		preHookErrorIgnore := w.Context().IgnorePreRunHooksError()

		log.Infof("Executing pre-hooks for service %s ...", w.Service().Name())

		for _, h := range w.Context().PreHooks() {
			log.Infof("executing pre-hook %s for service %s ...", h.Name(), w.Service().Name())

			hErr := h.Execute()
			if hErr != nil {
				preHookError = true

				log.Errorf("pre-hook %s failed for service %s", h.Name(), w.Service().Name())

				if !preHookErrorIgnore {
					return
				}
			}
		}
	}()

	if preHookError && !w.Context().IgnorePreRunHooksError() {
		log.Errorf("pre-hooks failed for service %s. Not starting the service", w.Service().Name())

		return
	}

	log.Infof("starting service %s ...", w.Service().Name())
	w.Service().Start(w.Context())

	func() {
		postHookErrorIgnore := w.Context().IgnorePostRunHooksError()

		log.Infof("Executing post-hooks for service %s ...", w.Service().Name())

		for _, h := range w.Context().PostHooks() {
			log.Infof("executing post-hook %s for service %s ...", h.Name(), w.Service().Name())

			hErr := h.Execute()
			if hErr != nil {
				log.Errorf("post-hook %s failed for service %s", h.Name(), w.Service().Name())

				if !postHookErrorIgnore {
					return
				}
			}
		}
	}()

}

// implement the stat stop logic here. as we have access to the service and  its context.

// NewWrapper returns a new instance of the wrapper.
func NewWrapper(s Service, wg *sync.WaitGroup, opts ...Option) *Wrapper {
	sCtx := &Context{
		terminationChan: make(chan struct{}),
		wg:              wg,
	}

	for _, opt := range opts {
		opt(sCtx)
	}

	return &Wrapper{
		s:    s,
		sCtx: sCtx,
	}
}
