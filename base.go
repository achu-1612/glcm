package glcm

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/achu-1612/glcm/log"
	"github.com/achu-1612/glcm/service"

	fig "github.com/common-nighthawk/go-figure"
)

// Base is the blueprint for the runner.
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
	// Only a ShutDown() call will stop the runner.
	// Even after all the registered services are stopped, runner would still be running.
	Wait()

	// restartService(...string) error
	// restartAllServices() error

	// stopService(...string) error
	// stopAllServices() error
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
	ctx       context.Context
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
	if r.isRunning {
		log.Info("Runner is already running. Doing nothing.")

		return
	}

	fig.NewColorFigure("GLCM", "isometric1", "green", true).Print()

	r.ctx = ctx

	if r.ctx == nil {
		log.Warn("Base Context is empty. Using the background context.")

		r.ctx = context.Background()
	}

	// Adding the base runner to the wait group.
	// This is to keep the runner running even
	// if all the services are stopped.
	// if no service has been registered.
	r.wg.Add(1)

	log.Info("Booting up Base Runner ...")

	for _, svc := range r.svc {
		// Adding the service to the wait group.
		// r.wg.Add(1)

		go svc.Start()
	}

	r.isRunning = true
}

// Wait waits for the runner to stop.
func (r *runner) Wait() {
	log.Info("Waiting to catch shutdown signal...")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit,
		syscall.SIGTERM, syscall.SIGINT,
		syscall.SIGQUIT, syscall.SIGHUP)

	func() {
		select {
		case <-quit:
			return
		case <-r.ctx.Done():
			return
		}
	}()

	log.Info("Received shutdown signal !!!")

	r.Shutdown()

	log.Infof("Waiting for %d service(s) to stop ...", len(r.svc))

	r.wg.Wait()

	log.Info("All services stopped. Exiting ...")

	r.isRunning = false
}

// Shutdown shuts down the runner. This will stop all the registered services.
func (r *runner) Shutdown() {
	log.Info("Shutting down Runner...")

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
