package group

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/achu-1612/glcm/log"
	"github.com/achu-1612/glcm/service"
)

type Base interface {
	// BottUp boots up the runner. This will start all the registered services.
	BootUp(context.Context)

	// Shutdown shuts down the runner. This will stop all the registered services.
	Shutdown()

	// RegisterService registers a service with the runner.
	RegisterService(service.Service, ...service.Option) error

	// IsRunning returns true if the runner is running, otherwise false.
	IsRunning() bool

	// Wait waits for the runner to stop.
	// Note: This is a blocking call. It is to be called after BootUp.
	Wait()
}

type serviceWrapper struct {
	s         service.Service
	preHooks  []func()
	postHooks []func()
}

// NewRunner returns a new instance of the runner.
func NewRunner() Base {
	return &runner{
		wg: &sync.WaitGroup{},
		mu: &sync.Mutex{},
	}
}

// runner implements the Base interface.
type runner struct {
	wg  *sync.WaitGroup
	mu  *sync.Mutex
	svc map[string]service.Wrapper

	// isRunning is a flag to indicate if the runner is running or not.
	isRunning bool
}

// IsRunning returns true if the runner is running, otherwise false.
func (r *runner) IsRunning() bool {
	return r.isRunning
}

// RegisterService registers a service with the runner.
func (r *runner) RegisterService(svc service.Service, opts ...service.Option) error {
	if r.isRunning {
		return ErrRunnerAlreadyRunning
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.svc[svc.Name()]; ok {
		return ErrServiceAlreadyExists
	}

	r.svc[svc.Name()] = *service.NewWrapper(svc, opts...)

	return nil
}

// BootUp boots up the runner. This will start all the registered services.
func (r *runner) BootUp(ctx context.Context) {
	if ctx == nil {
		log.Warn("context is nil. Using the background context.")

		ctx = context.Background()
	}

	if r.isRunning {
		log.Info("runner is already running. Doing nothing.")

		return
	}

	r.wg.Add(1)

	log.Info("Booting up Runner ...")
	r.isRunning = true
}

func (r *runner) Wait() {
	log.Info("waiting to catch shutdown signal ...")

	catchShutdownSignal()

	r.Shutdown()

	r.wg.Wait()

	r.isRunning = false
}

// Shutdown shuts down the runner. This will stop all the registered services.
func (r *runner) Shutdown() {
	log.Info("shutting down Runner ...")

	r.wg.Done()
}

// catchShutdownSignal - catch a shutdown signal.
func catchShutdownSignal() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit,
		syscall.SIGTERM, syscall.SIGINT,
		syscall.SIGQUIT, syscall.SIGHUP)

	<-quit
}
