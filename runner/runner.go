package runner

import (
	"context"
	"io"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/achu-1612/glcm/log"
	"github.com/achu-1612/glcm/service"

	fig "github.com/common-nighthawk/go-figure"
)

// runner implements the Base interface.
type runner struct {
	// swg is a wait group to wait for all the services to stop.
	swg *sync.WaitGroup

	// mu is a mutex to protect the runner state.
	mu *sync.Mutex

	// svc is a map of services registered with the runner.
	svc map[string]*service.Wrapper

	// isRunning is a flag to indicate if the runner is running or not.
	isRunning bool

	// ctx is the base context for the runner.
	ctx context.Context

	// hideBanner is a flag to indicate if the banner should be hidden or not.
	hideBanner bool

	// suppressLog is a flag to indicate if the logs should be suppressed or not.
	suppressLog bool
}

// NewRunner returns a new instance of the runner.
func NewRunner(opts ...Options) Base {
	r := &runner{
		svc: make(map[string]*service.Wrapper),
		mu:  &sync.Mutex{},
		swg: &sync.WaitGroup{},
	}

	// mutate the runner with the options.
	for _, opt := range opts {
		opt(r)
	}

	if r.suppressLog {
		log.SetOutput(io.Discard)
	}

	return r
}

// IsRunning returns true if the runner is running, otherwise false.
func (r *runner) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.isRunning
}

// RegisterService registers a service with the runner.
func (r *runner) RegisterService(svc service.Service, opts ...service.Option) error {
	if svc == nil {
		return ErrRegisterNilService
	}

	if r.IsRunning() {
		return ErrRunnerAlreadyRunning
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.svc[svc.Name()]; ok {
		return ErrRegisterServiceAlreadyExists
	}

	r.svc[svc.Name()] = service.NewWrapper(svc, r.swg, opts...)

	return nil
}

// BootUp boots up the runner. This will start all the registered services.
func (r *runner) BootUp(ctx context.Context) error {
	if r.IsRunning() {
		return ErrRunnerAlreadyRunning
	}

	if !r.hideBanner {
		fig.NewColorFigure("GLCM", "isometric1", "green", true).Print()
	}

	r.ctx = ctx

	if r.ctx == nil {
		log.Warn("Base Context is empty. Using the background context.")

		r.ctx = context.Background()
	}

	log.Info("Booting up the Runner ...")

	r.isRunning = true

	quit := make(chan os.Signal, 1)

	signal.Notify(quit,
		syscall.SIGTERM, syscall.SIGINT,
		syscall.SIGQUIT, syscall.SIGHUP)

	// TODO: run the reconciler only if there is service state change.

	func() {
		t := time.NewTicker(time.Second)

		for {
			select {
			case <-quit:
				return
			case <-r.ctx.Done():
				return
			case <-t.C:
				r.reconcile()
			}
		}
	}()

	log.Info("Received shutdown signal !!!")

	r.Shutdown()

	log.Info("All services stopped. Exiting ...")

	return nil
}

// reconcile takes necessary actions on the services based on their state.
func (r *runner) reconcile() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range r.svc {
		// The services are expected to be in the registered state at first.
		// If the service is registered, then start the service on fist rec cycle.
		if svc.Status == service.StatusRegistered {
			go svc.Start()
		}

		// skip the service if it is already pending start.
		if svc.AutoRestart.PendingStart.Load() {
			continue
		}

		// auto restart the service if it is exited (not stopped) and auto-restart is enabled for the service
		// the service will not be started automatically if it stopped by the runner.
		if svc.Status == service.StatusExited && svc.AutoRestart.Enabled {
			if svc.AutoRestart.RetryCount >= svc.AutoRestart.MaxRetries {
				log.Infof("Service %s reached max retries. Not restarting ...", svc.Name())

				continue
			}

			backoffDuration := time.Duration(0)

			if svc.AutoRestart.Backoff {
				backoffDuration = time.Duration(
					math.Pow(float64(svc.AutoRestart.BackOffExponent), float64(svc.AutoRestart.RetryCount)),
				) * time.Second
			}

			svc.AutoRestart.RetryCount++

			// using same flow for both immediate and backoff restarts.
			svc.AutoRestart.PendingStart.Store(true)

			go func() {
				if backoffDuration > 0 {
					log.Infof("Service %s backing-off. Restarting in %s ...", svc.Name(), backoffDuration)
					<-time.After(backoffDuration)
				}

				svc.Start()
			}()
		}
	}
}

// Shutdown shuts down the runner. This will stop all the registered services.
func (r *runner) Shutdown() {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Info("Shutting down Runner...")

	for _, svc := range r.svc {
		if svc.Status == service.StatusRunning {
			svc.Stop()
		}
	}

	log.Infof("Waiting for %d service(s) to stop ...", len(r.svc))

	r.swg.Wait()

	r.isRunning = false
}

// StopAllServices stops all the registered/running services.
func (r *runner) StopAllServices() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range r.svc {
		if svc.Status == service.StatusRunning {
			svc.Stop()
		}
	}

	r.swg.Wait()
}

// StopService stops the given list of services.
func (r *runner) StopService(name ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, n := range name {
		if svc, ok := r.svc[n]; ok && svc.Status == service.StatusRunning {
			svc.StopAndWait()
		}
	}

	return nil
}

// RestartService restarts the given list of services.
func (r *runner) RestartService(name ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, n := range name {
		if svc, ok := r.svc[n]; ok {
			if svc.Status == service.StatusRunning {
				svc.StopAndWait()
			}

			go svc.Start()
		}
	}

	return nil
}

// RestartAllServices restarts all the registered/running services.
func (r *runner) RestartAllServices() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range r.svc {
		if svc.Status == service.StatusRunning {
			svc.StopAndWait()
		}

		go svc.Start()
	}
}
