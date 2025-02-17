package glcm

//go:generate mockgen -package glcm -destination spec.mock.go -source spec.go -self_package "github.com/achu-1612/glcm"

// Hook is an interface which represents a single hook.
// When a servcice is regsited in the runner, implementations of the Hndler interface can be registered too.
// Based on the nature of the hook (pre-run or post-run), the hook will be executed.
type Hook interface {
	// Execute executes the hook method.
	Execute() error

	// Name returns the name of the hook.
	Name() string
}

// Service defines an interface which represents a single service and the
// operations that can be performed on the service.
// Note:
// The service should be able to handle the termination signal.
// At this point, I don't think we need to have a Stop or Restart method.
// Once the termincation channel is closed, the service should stop.
// If the service needs to be restarted, the runner will take care of it, internally.
type Service interface {
	// Name returns the name of the service.
	Name() string

	// Start executes/boots-up/starts a service.
	Start(Terminator)
}

// Terminator defines an indicator to the service to stop.
type Terminator interface {
	// TermCh returns a channel which will be closed when the service should stop.
	TermCh() chan struct{}
}

// Runner represents the interface for the base runner methods.
type Runner interface {
	// IsRunning returns true if the runner is running, otherwise false.
	IsRunning() bool

	// RegisterService registers a service with the runner.
	RegisterService(Service, ServiceOptions) error

	// Shutdown stops all the services and the runner.
	Shutdown()

	// StopAllServices stops all the services.
	StopAllServices()

	// StopService stops the specified services.
	StopService(...string) error

	// RestartService restarts the specified services.
	RestartService(...string) error

	// RestartAllServices restarts all the services.
	RestartAllServices()

	// BootUp starts the runner.
	BootUp() error

	// Status returns the status of the runner along with the status of each registered service.
	Status() *RunnerStatus
}

// Wrapper is an interface which represents the wraper around the service.
type Wrapper interface {
	// Name returns the name of the service.
	Name() string

	// Status returns the status of the service/wrapper.
	Status() ServiceStatus

	// TermCh returns the termination channel for the service.
	TermCh() chan struct{}

	// Start starts the services in the wrapper.
	Start()

	// Stop stops the service in the wrapper.
	Stop()

	// StopAndWait stops the service (in the wrapper) and waits for it to stop.
	StopAndWait()

	// AutoRestart returns the auto-restart configuration for the wrapper.
	AutoRestart() *AutoRestart
}
