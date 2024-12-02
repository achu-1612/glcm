package service

import (
	"sync"

	"github.com/achu-1612/glcm/hook"
	"github.com/achu-1612/glcm/log"
)

// Terminator interface abstract other implementation of the Wrapper.
// This is used as an indicator to the service to stop.
type Terminator interface {
	TermCh() chan struct{}
}

// Wrapper is a wrapper around the service and its context.
type Wrapper struct {
	s Service

	// PreHooks are the hooks that will be executed before starting the service.
	preHooks []hook.Handler

	// PostHooks are the hooks that will be executed after stopping the service.
	postHooks []hook.Handler

	// tc (termination channel) is a channel which will be used to direct the service to stop.
	// The channel will be closed the service is to be stopped.
	// 1. The Runner is shutting down.
	// 2. The Stop() method is called on the service.
	tc chan struct{}

	// dic (done indication channel) is a channel which will be close on calling Done() method.
	// This will indicate the runner that the service has stopped.
	dic chan struct{}

	// wg is the service wait group. Not the same as the base runner wait group.
	wg *sync.WaitGroup

	// isRunning is a flag to indicate if the service is running or not.
	isRunning bool

	// TODO: have a counter to indicate service restarts.
	// We can also provide a way for the user to specify whether they want the hooks
	// to get executed everytyhing the service stop and start of just for the first time.
}

// NewWrapper returns a new instance of the wrapper.
func NewWrapper(s Service, wg *sync.WaitGroup, opts ...Option) *Wrapper {
	w := &Wrapper{
		s:  s,
		wg: wg,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Done marks the services as done in the workergroup and clsoes the indication channel.
func (w *Wrapper) done() {
	w.wg.Done()

	close(w.dic)
}

// Wait blocks the caller until the service is stopped.
func (w *Wrapper) wait() {
	<-w.dic
}

// TermCh returns the termination channel for the service.
// The service implmentation is expected to listen to this channel and
// stop the service when it is closed.
func (w *Wrapper) TermCh() chan struct{} {
	return w.tc
}

// reallocate the chan before starting if it is nil
func (w *Wrapper) Start() {
	if w.isRunning {
		log.Infof("Service %s is already running", w.s.Name())

		return
	}

	// we don't know if this is the first time the service is getting started.
	// So, we need to reallocate the channels.
	w.dic = make(chan struct{})
	w.tc = make(chan struct{})

	w.wg.Add(1)
	defer func() {
		w.done() // indicate the worker group that the service has stopped.

		log.Infof("service %s stopped", w.s.Name())
	}()

	// call the pre exec hooks
	func() {
		log.Infof("Executing pre-hooks for service %s ...", w.s.Name())

		for _, h := range w.preHooks {
			log.Infof("executing pre-hook %s for service %s ...", h.Name(), w.s.Name())

			hErr := h.Execute()
			if hErr != nil {
				log.Errorf("pre-hook %s failed for service %s", h.Name(), w.s.Name())
			}
		}
	}()

	// start the service
	log.Infof("starting service %s ...", w.s.Name())
	w.isRunning = true
	w.s.Start(w)

	// call the post exec hooks.
	// Note: we don't really need the ignore flag here,,
	// as there is nothing for us to do, if the post hooks fail.
	func() {

		log.Infof("Executing post-hooks for service %s ...", w.s.Name())

		for _, h := range w.postHooks {
			log.Infof("executing post-hook %s for service %s ...", h.Name(), w.s.Name())

			hErr := h.Execute()
			if hErr != nil {
				log.Errorf("post-hook %s failed for service %s", h.Name(), w.s.Name())
			}
		}
	}()
}

// stop stops the service. It acts like a wrapper around the service's stop method.
// to be consumed by Stop() and StopAndWait() methods.
func (w *Wrapper) stop() error {
	if !w.isRunning {
		return ErrServiceNotRunning
	}

	log.Infof("Stopping service %s ...", w.s.Name())

	close(w.tc)

	w.isRunning = false

	return nil
}

// Stop stops the service.
func (w *Wrapper) Stop() {
	if err := w.stop(); err != nil {
		log.Errorf("Failed to stop service %s: %v", w.s.Name(), err)
	}
}

// StopAndWait stops the service and waits for it to exit.
func (w *Wrapper) StopAndWait() {
	if err := w.stop(); err != nil {
		log.Errorf("Failed to stop service %s: %v", w.s.Name(), err)

		return
	}

	log.Infof("Waiting for the service %s to exit ...", w.s.Name())

	w.wait()
}
