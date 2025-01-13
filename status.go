package glcm

// ServiceStatus represents the status of the service.
type ServiceStatus string

// Status options for the service.
const (
	ServiceStatusRegistered          ServiceStatus = "registered"
	ServiceStatusRunning             ServiceStatus = "running"
	ServiceStatusExited              ServiceStatus = "exited"
	ServiceStatusStopped             ServiceStatus = "stopped"
	ServiceStatusScheduled           ServiceStatus = "scheduled"
	ServiceStatusScheduledForRestart ServiceStatus = "scheduled-for-restart"
	ServiceStatusExhausted           ServiceStatus = "exhausted"
)
