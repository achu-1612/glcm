package runner

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/achu-1612/glcm/log"
	"github.com/achu-1612/glcm/service"

	fig "github.com/common-nighthawk/go-figure"
)

// NewRunner returns a new instance of the runner.
func NewRunner(opts ...Options) Base {
	r := &runner{
		svc: make(map[string]*service.Wrapper),
		mu:  &sync.Mutex{},
		swg: &sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.suppressLog {
		log.SetOutput(io.Discard)
	}

	return r
}

// runner implements the Base interface.
type runner struct {
	swg *sync.WaitGroup
	mu  *sync.Mutex
	svc map[string]*service.Wrapper

	// isRunning is a flag to indicate if the runner is running or not.
	isRunning bool
	ctx       context.Context

	// hideBanner is a flag to indicate if the banner should be hidden or not.
	hideBanner bool

	// suppressLog is a flag to indicate if the logs should be suppressed or not.
	suppressLog bool
}

// IsRunning returns true if the runner is running, otherwise false.
func (r *runner) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.isRunning
}

// RegisterService registers a service with the runner.
func (r *runner) RegisterService(svc service.Service, opts ...service.Option) error {
	if r.IsRunning() {
		return ErrRunnerAlreadyRunning
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.svc[svc.Name()]; ok {
		return ErrServiceAlreadyExists
	}

	r.svc[svc.Name()] = service.NewWrapper(svc, r.swg, opts...)

	return nil
}

// BootUp boots up the runner. This will start all the registered services.
func (r *runner) BootUp(ctx context.Context) {
	if r.IsRunning() {
		log.Info("Runner is already running. Nothing to do.")

		return
	}

	if !r.hideBanner {
		fig.NewColorFigure("GLCM", "isometric1", "green", true).Print()
	}

	r.ctx = ctx

	if r.ctx == nil {
		log.Warn("Base Context is empty. Using the background context.")

		r.ctx = context.Background()
	}

	log.Info("Booting up Base Runner ...")

	for _, svc := range r.svc {
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

	r.swg.Wait()

	log.Info("All services stopped. Exiting ...")

	r.isRunning = false
}

// Shutdown shuts down the runner. This will stop all the registered services.
func (r *runner) Shutdown() {
	log.Info("Shutting down Runner...")

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range r.svc {
		svc.Stop()
	}

	r.isRunning = false
}

func (r *runner) RestartAllServices() {
	r.StopAllServices()

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range r.svc {
		go svc.Start()
	}
}

func (r *runner) StopAllServices() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range r.svc {
		svc.Stop()
	}

	r.swg.Wait()
}

func (r *runner) RestartService(name ...string) error {
	if err := r.StopService(name...); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, n := range name {
		if svc, ok := r.svc[n]; ok {
			go svc.Start()
		}
	}

	return nil
}

func (r *runner) StopService(name ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, n := range name {
		if svc, ok := r.svc[n]; ok {
			svc.StopAndWait()
		}
	}

	return nil
}
