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
	svc map[string]Wrapper

	// isRunning is a flag to indicate if the runner is running or not.
	isRunning bool

	// ctx is the base context for the runner.
	ctx context.Context

	// hideBanner is a flag to indicate if the banner should be hidden or not.
	hideBanner bool

	// socket holds the object of the socket.
	socket *socket

	// socketPath represents the path to the socket file.
	socketPath string

	// allowedUIDs represents the allowed user ids to interact with the socket.
	allowedUIDs []int

	// shutdownTimeout represents the timeout for shutting down the runner.
	shutdownTimeout time.Duration
}

// NewRunner returns a new instance of the runner.
func NewRunner(ctx context.Context, opts RunnerOptions) Runner {
	opts.Sanitize()

	r := &runner{
		svc:             make(map[string]Wrapper),
		mu:              &sync.Mutex{},
		swg:             &sync.WaitGroup{},
		ctx:             ctx,
		hideBanner:      opts.HideBanner,
		socketPath:      opts.SocketPath,
		allowedUIDs:     opts.AllowedUID,
		shutdownTimeout: opts.ShutdownTimeout,
	}

	if opts.Verbose {
		log.SetOutput(io.Discard)
	}

	if ctx == nil {
		log.Warn("Base Context is empty. Using the background context.")

		r.ctx = context.Background()
	}

	if opts.Socket {
		socket, err := newSocket(r, opts.SocketPath, opts.AllowedUID)
		if err != nil {
			log.Errorf("creating a socket: %v", err)
		}

		r.socket = socket
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

	r.mu.Lock()
	defer r.mu.Unlock()

	sName := svc.Name()

	if _, ok := r.svc[sName]; ok {
		return ErrRegisterServiceAlreadyExists
	}

	opts.Sanitize()

	r.svc[sName] = NewWrapper(svc, r.swg, opts)

	return nil
}

// DeregisterService deregisters a service from the runner.
// If the service is running, it will be stopped before deregistering.
func (r *runner) DeregisterService(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.svc[name]; !ok {
		return ErrDeregisterServiceNotFound
	}

	// stop the service if it is running.
	if r.svc[name].Status() == ServiceStatusRunning {
		r.svc[name].Stop()
	}

	delete(r.svc, name)

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

	log.Info("Booting up the Runner ...")

	r.isRunning = true

	quit := make(chan os.Signal, 1)

	signal.Notify(quit,
		syscall.SIGTERM, syscall.SIGINT,
		syscall.SIGQUIT, syscall.SIGHUP)

	if r.socket != nil {
		// start the socket inside a go-routine, as Start is a blocking call.,
		// it will be shutdown automatically when we receive the signal.
		go func() {
			if err := r.socket.start(); err != nil {
				log.Errorf("failed to start socket: %v", err)
			}
		}()

		defer r.socket.shutdown()
	}

	t := time.NewTicker(time.Second)

	for {
		select {
		case <-quit:
			log.Info("Received shutdown signal. Shutting down the runner ...")
			r.Shutdown()

			return nil
		case <-r.ctx.Done():
			log.Info("Received shutdown signal. Shutting down the runner ...")
			r.Shutdown()

			return nil
		case <-t.C:
			r.reconcile()
		}
	}
}

// reconcile takes necessary actions on the services based on their state.
func (r *runner) reconcile() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, w := range r.svc {
		log.Infof("Reconciling service: %s, current status: %s", w.Name(), w.Status())

		// The services are expected to be in the registered state at first.
		// If the service is registered, then start the service on first rec cycle.
		if w.Status() == ServiceStatusRegistered {
			log.Infof("Service %s is registered. Starting service ...", w.Name())

			go w.Start()
		}

		// skip the service if it is already pending start.
		if w.AutoRestart().PendingStart.Load() {
			log.Infof("Service %s is pending start. Skipping ...", w.Name())

			continue
		}

		// auto restart the service if it is exited (not stopped) and auto-restart is enabled for the service
		// the service will not be started automatically if it stopped by the runner.
		if w.Status() == ServiceStatusExited && w.AutoRestart().Enabled {
			if w.AutoRestart().RetryCount >= w.AutoRestart().MaxRetries {
				log.Infof("Service %s reached max retries. Not restarting ...", w.Name())

				continue
			}

			backoffDuration := time.Duration(0)

			if w.AutoRestart().Backoff {
				backoffDuration = time.Duration(
					math.Pow(float64(w.AutoRestart().BackoffExponent), float64(w.AutoRestart().RetryCount)),
				) * time.Second
			}

			w.AutoRestart().RetryCount++

			// using same flow for both immediate and backoff restarts.
			w.AutoRestart().PendingStart.Store(true)

			go func() {
				if backoffDuration > 0 {
					log.Infof("Service %s backing-off. Restarting in %s ...", w.Name(), backoffDuration)

					<-time.After(backoffDuration)
				}

				log.Infof("Service %s restarting now ...", w.Name())

				w.Start()
			}()
		}
	}
}

// Shutdown shuts down the runner. This will stop all the registered services.
func (r *runner) Shutdown() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRunning {
		log.Warn("Runner is not running. Skipping shutdown ...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.shutdownTimeout)
	defer cancel()

	log.Info("Shutting down Runner...")

	for _, svc := range r.svc {
		if svc.Status() == ServiceStatusRunning {
			go func(svc Wrapper) {
				svc.Stop()
			}(svc)
		}
	}

	gracefulShutdown := make(chan struct{})

	go func() {
		log.Infof("Waiting for %d service(s) to stop ...", len(r.svc))

		r.swg.Wait()

		close(gracefulShutdown)
	}()

	select {
	case <-ctx.Done():
		log.Infof("Graceful shutdown timed out. Forcing shutdown ...")
	case <-gracefulShutdown:
		log.Infof("All services stopped gracefully.")
	}

	r.isRunning = false
}

// StopAllServices stops all the registered/running services.
func (r *runner) StopAllServices() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range r.svc {
		if svc.Status() == ServiceStatusRunning {
			go func(svc Wrapper) {
				svc.Stop()
			}(svc)
		}
	}

	r.swg.Wait()
}

// StopService stops the given list of services.
func (r *runner) StopService(name ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, n := range name {
		if svc, ok := r.svc[n]; ok && svc.Status() == ServiceStatusRunning {
			svc.Stop()
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
			if svc.Status() == ServiceStatusRunning {
				svc.Stop()
				go svc.Start()
			}
		}
	}

	return nil
}

// RestartAllServices restarts all the registered/running services.
func (r *runner) RestartAllServices() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range r.svc {
		if svc.Status() == ServiceStatusRunning {
			svc.Stop()
			go svc.Start()
		}
	}
}

func (r *runner) Status() *RunnerStatus {
	r.mu.Lock()
	defer r.mu.Unlock()

	status := &RunnerStatus{
		IsRunning: r.isRunning,
		Services:  make(map[string]ServiceInfo),
	}

	for _, svc := range r.svc {
		status.Services[svc.Name()] = ServiceInfo{
			Status: svc.Status(),
			Uptime: svc.Uptime(),
		}
	}

	return status
}
