package service

// Service defines an interface which represents a single service and the
// operations that can be performed on the service.
type Service interface {
	// Name returns the name of the service.
	Name() string

	// Start executes/boots-up/starts a service.
	Start()

	// Restart restarts the service.
	Restart() // TODO: do we really need this or we can just use Stop() and Start() ?

	// Stop stops the service.
	Stop()

	// Status returns the status of the service.
	Status() string
}
