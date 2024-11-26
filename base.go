package glcm

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

// NewRunner returns a new instance of the runner.
func NewRunner() Base {
	return &runner{
		svc: make(map[string]*service.Wrapper),
		wg:  &sync.WaitGroup{},
		mu:  &sync.Mutex{},
	}
}

// runner implements the Base interface.
type runner struct {
	wg  *sync.WaitGroup
	mu  *sync.Mutex
	svc map[string]*service.Wrapper

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

	r.svc[svc.Name()] = service.NewWrapper(svc, r.wg, opts...)

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

	// Adding the base runner to the wait group.
	// This is to keep the runner running even
	// if all the services are stopped.
	// if no service has been registered.
	r.wg.Add(1)

	log.Info("Booting up Runner ...")

	for _, svc := range r.svc {

		// Adding the service to the wait group.
		r.wg.Add(1)

		go func(svc *service.Wrapper) {
			defer svc.Context().Done()
			defer log.Infof("service %s stopped", svc.Service().Name())

			log.Infof("starting service %s ...", svc.Service().Name())
			svc.Service().Start(svc.Context())
		}(svc)
	}

	r.isRunning = true
}

// Wait waits for the runner to stop.
func (r *runner) Wait() {
	log.Info("waiting to catch shutdown signal ...")

	catchShutdownSignal()

	log.Info("received shutdown signal ...")

	r.Shutdown()

	r.wg.Wait()

	r.isRunning = false
}

// Shutdown shuts down the runner. This will stop all the registered services.
func (r *runner) Shutdown() {
	log.Info("shutting down Runner ...")

	for _, svc := range r.svc {
		close(svc.Context().TermCh())
	}

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
