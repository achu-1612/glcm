package glcm

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/achu-1612/glcm/log"
)

const (
	defaultMaxRetries = 10
	defaultBackoffExp = 2
)

// Wrapper is a wrapper around the service and its context.
type wrapper struct {
	s Service

	// preHooks are the hooks that will be executed before starting the service.
	preHooks []Hook

	// postHooks are the hooks that will be executed after stopping the service.
	postHooks []Hook

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

	// shutdownRequest is a flag to indicate if the service is requested to stop by the runner.
	shutdownRequest atomic.Bool

	// status is the current status of the service.
	status ServiceStatus

	// autorestart related configuration.
	autoRestart AutoRestart

	// scheduling related configuration.

	ScheduleEnabled        bool          // flag to indicate if scheduling is enabled.
	ScheduleCronExpression string        // cron expression for scheduling the service.
	ScheduleTimeOut        time.Duration // execution timeout for the service.
	ScheduleMaxRuns        int           // maximum number of runs for the service.
}

// AutoRestart is the configuration set for auto-restart.
type AutoRestart struct {
	Enabled         bool        // flag to indicate if auto-restart is enabled.
	MaxRetries      int         // maximum number of retries.
	Backoff         bool        // flag to indicate if backoff is enabled.
	BackoffExponent int         // exponent for the backoff.
	RetryCount      int         // current number of retries for the service.
	PendingStart    atomic.Bool // flag to indicate if the service is pending for a start after the backoff.
}

// NewWrapper returns a new instance of the service Wrapper.
func NewWrapper(s Service, wg *sync.WaitGroup, opts ServiceOptions) Wrapper {
	w := &wrapper{
		s:         s,
		wg:        wg,
		preHooks:  opts.PreHooks,
		postHooks: opts.PostHooks,
		status:    ServiceStatusRegistered,
		autoRestart: AutoRestart{
			RetryCount:      0,
			Enabled:         opts.AutoStart.Enabled,
			MaxRetries:      opts.AutoStart.MaxRetries,
			Backoff:         opts.AutoStart.Backoff,
			BackoffExponent: opts.AutoStart.BackOffExponent,
			PendingStart:    atomic.Bool{},
		},
		ScheduleEnabled:        opts.Schedule.Enabled,
		ScheduleCronExpression: opts.Schedule.Cron,
		ScheduleTimeOut:        opts.Schedule.TimeOut,
		ScheduleMaxRuns:        opts.Schedule.MaxRuns,
	}

	// sanitize the auto-restart configuration.

	if w.autoRestart.MaxRetries == 0 {
		log.Warnf("MaxRetries is not set for service %s. Setting it to default value %d", w.s.Name(), defaultMaxRetries)

		w.autoRestart.MaxRetries = defaultMaxRetries
	}

	if w.autoRestart.BackoffExponent == 0 {
		log.Warnf("BackoffExponent is not set for service %s. Setting it to default value %d", w.s.Name(), defaultBackoffExp)
		w.autoRestart.BackoffExponent = defaultBackoffExp
	}

	return w
}

func (w *wrapper) AutoRestart() *AutoRestart {
	return &w.autoRestart
}

func (w *wrapper) Name() string {
	return w.s.Name()
}

func (w *wrapper) Status() ServiceStatus {
	return w.status
}

// Done marks the services as done in the workergroup and closes the indication channel.
func (w *wrapper) done() {
	// indicate whether the service has stopped by runner or exited on its own.
	// if the service is stopped by the runner (shudownRequest will be set to true), then the status will be stopped.
	// if the service has exited on its own, then the status will be exited.
	if w.shutdownRequest.Load() {
		w.status = ServiceStatusStopped
	} else {
		w.status = ServiceStatusExited
	}

	// clearing the shutdown request flag.
	w.shutdownRequest.Store(false)

	w.wg.Done()

	close(w.dic)
}

// Wait blocks the caller until the service is stopped.
func (w *wrapper) wait() {
	<-w.dic
}

// TermCh returns the termination channel for the service.
// The service implmentation is expected to listen to this channel and
// stop the service when it is closed.
func (w *wrapper) TermCh() chan struct{} {
	return w.tc
}

// reallocate the chan before starting if it is nil
func (w *wrapper) Start() {
	if w.status == ServiceStatusRunning {
		log.Infof("Service %s is already running", w.s.Name())

		return
	}

	// we don't know if this is the first time the service is getting started.
	// So, we need to reallocate the channels.
	w.dic = make(chan struct{})
	w.tc = make(chan struct{})

	w.wg.Add(1)

	defer func() {
		w.done() // finalizer for the service wrapper.

		log.Infof("service %s status [%s]", w.s.Name(), w.status)
	}()

	// call the pre exec hooks
	func() {
		log.Infof("Executing pre-hooks for service %s ...", w.s.Name())

		for _, h := range w.preHooks {
			log.Infof("executing pre-hook %s for service %s ...", h.Name(), w.s.Name())

			hErr := h.Execute()
			if hErr != nil {
				log.Errorf("pre-hook %s failed for service %s: %v", h.Name(), w.s.Name(), hErr)
			}
		}
	}()

	// start the service
	log.Infof("starting service %s ...", w.s.Name())

	w.status = ServiceStatusRunning
	w.autoRestart.PendingStart.Store(false)
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
				log.Errorf("post-hook %s failed for service %s: %v", h.Name(), w.s.Name(), hErr)
			}
		}
	}()
}

// Stop stops the service and waits for it to exit.
func (w *wrapper) Stop() {
	if !(w.status == ServiceStatusRunning) {
		return
	}

	log.Infof("Stopping service %s ...", w.s.Name())

	close(w.tc)

	w.shutdownRequest.Store(true)

	log.Infof("Waiting for the service %s to exit ...", w.s.Name())

	w.wait()
}
