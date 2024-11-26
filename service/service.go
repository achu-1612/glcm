package service

// Service defines an interface which represents a single service and the
// operations that can be performed on the service.
// Note:
// The service should be able to handle the termination signal.
// At this point, I don't think we need to have a Stop or Restart method.
// Once the termincation channel is closed, the service should stop.
// If the service needs to be restarted, the runner should take care of it, internally.
type Service interface {
	// Name returns the name of the service.
	Name() string

	// Start executes/boots-up/starts a service.
	Start(*Context)

	// Status returns the status of the service.
	Status() string
}
