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

	// call the pre exec hooks
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

	// start the service
	log.Infof("starting service %s ...", w.Service().Name())
	w.Service().Start(w.Context())

	// call the post exec hooks.
	// Note: we don't really need the ignore flag here,,
	// as there is nothing for us to do, if the post hooks fail.
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

/*

Have a primary wait group for the base runner.
This will make sure the base runner is always running even if the
services are stopped.
Have a second wait group for the services.

This will help to wait for the services to end the execution,
when we call the stop all services.

each service we will have a channel, which will be closed when the service is done,
This will help us to wait for one services to end.

When the service will call done,
two things we will do:
1. close the channel
2. call the done on the wait group.

When the service is started, we will add the service to the wait group.


NO need to have a context as separate struct,
Just pull everything to the wrapper.

Have start stop and restart methods on the wrapper to handle things easily outside thw pkg.const
wrapper will consider all the cases like, if the service is already running,
olr stopped when we call the start methods, or stop all, or start all, etc.

Calling stop on the wrapper will wait on the channel whcih will idncate the at,
the service is done,
it wil be useful when we need to start a service right after we call a stop.const
Basically restart.


no need to track errors for the post hooks executions
if the post hook fails, we can log it and move on.

wraper will take care of the pre hooks and post hooks execution. already implemented.


*/
