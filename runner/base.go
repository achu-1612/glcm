package runner

import (
	"context"

	"github.com/achu-1612/glcm/service"
)

//go:generate mockgen -package runner -destination base.mock.go -source base.go -self_package "github.com/achu-1612/glcm/runner"

// Base is the blueprint for the runner.
type Base interface {
	// BootUp boots up the runner. This will start all the registered services.
	// Note: This is a blocking call. It is to be called after BootUp.
	// Only a ShutDown() call will stop the runner.
	// Even after all the registered services are stopped, runner would still be running.
	BootUp(context.Context) error

	// Shutdown shuts down the runner. This will stop all the registered services.
	Shutdown()

	// RegisterService registers a service with the runner.
	RegisterService(service.Service, ...service.Option) error

	// IsRunning returns true if the runner is running, otherwise false.
	IsRunning() bool

	// RestartService restarts the given list of services.
	RestartService(...string) error

	// RestartAllServices restarts all the registered/running services.
	RestartAllServices()

	// StopService stops the given list of services.
	StopService(...string) error

	// StopAllServices stops all the registered/running services.
	StopAllServices()
}
