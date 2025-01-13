package glcm

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

	fig "github.com/common-nighthawk/go-figure"
)

// runner implements the Base interface.
type runner struct {
	// swg is a wait group to wait for all the services to stop.
	swg *sync.WaitGroup

	// mu is a mutex to protect the runner state.
	mu *sync.Mutex

	// svc is a map of services registered with the runner.
	svc map[string]*Wrapper

	// isRunning is a flag to indicate if the runner is running or not.
	isRunning bool

	// ctx is the base context for the runner.
	ctx context.Context

	// hideBanner is a flag to indicate if the banner should be hidden or not.
	hideBanner bool

	// Verbose represents if the logs should be suppressed or not.
	verbose bool
}

// New returns a new instance of the runner.
func NewRunner(ctx context.Context, opts RunnerOptions) Runner {
	r := &runner{
		svc:        make(map[string]*Wrapper),
		mu:         &sync.Mutex{},
		swg:        &sync.WaitGroup{},
		ctx:        ctx,
		verbose:    opts.Verbose,
		hideBanner: opts.HideBanner,
	}

	if !r.verbose {
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
func (r *runner) RegisterService(svc Service, opts ServiceOptions) error {
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

	r.svc[svc.Name()] = NewWrapper(svc, r.swg, opts)

	return nil
}

// BootUp boots up the runner.
func (r *runner) BootUp() error {
	if r.IsRunning() {
		return ErrRunnerAlreadyRunning
	}

	if !r.hideBanner {
		fig.NewColorFigure("GLCM", "isometric1", "green", true).Print()
	}

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

	for _, w := range r.svc {
		// The services are expected to be in the registered state at first.
		// If the service is registered, then start the service on fist rec cycle.
		if w.Status == ServiceStatusRegistered {
			go w.Start()
		}

		// skip the service if it is already pending start.
		if w.AutoRestartPendingStart.Load() {
			continue
		}

		// auto restart the service if it is exited (not stopped) and auto-restart is enabled for the service
		// the service will not be started automatically if it stopped by the runner.
		if w.Status == ServiceStatusExited && w.AutoRestartEnabled {
			if w.AutoRestartRetryCount >= w.AutoRestartMaxRetries {
				log.Infof("Service %s reached max retries. Not restarting ...", w.Name())

				continue
			}

			backoffDuration := time.Duration(0)

			if w.AutoRestartBackoff {
				backoffDuration = time.Duration(
					math.Pow(float64(w.AutoRestartBackoffExponent), float64(w.AutoRestartRetryCount)),
				) * time.Second
			}

			w.AutoRestartRetryCount++

			// using same flow for both immediate and backoff restarts.
			w.AutoRestartPendingStart.Store(true)

			go func() {
				if backoffDuration > 0 {
					log.Infof("Service %s backing-off. Restarting in %s ...", w.Name(), backoffDuration)
					<-time.After(backoffDuration)
				}

				w.Start()
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
		if svc.Status == ServiceStatusRunning {
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
		if svc.Status == ServiceStatusRunning {
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
		if svc, ok := r.svc[n]; ok && svc.Status == ServiceStatusRunning {
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
			if svc.Status == ServiceStatusRunning {
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
		if svc.Status == ServiceStatusRunning {
			svc.StopAndWait()
		}

		go svc.Start()
	}
}
