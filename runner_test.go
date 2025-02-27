package glcm

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIsRunning(t *testing.T) {
	r := NewRunner(context.Background(), RunnerOptions{})

	// Test when runner is not running
	assert.False(t, r.IsRunning(), "Expected runner to not be running")

	// Start the runner
	go func() {
		if err := r.BootUp(); err != nil {
			t.Errorf("Error while booting up the runner: %v", err)
		}
	}()

	<-time.After(time.Second * 10)

	// Test when runner is running
	assert.True(t, r.IsRunning(), "Expected runner to be running")

	// Shutdown the runner
	r.Shutdown()

	// Test when runner is not running after shutdown
	assert.False(t, r.IsRunning(), "Expected runner to not be running after shutdown")
}

func TestStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService1 := NewMockService(ctrl)
	mockService1.EXPECT().Name().Return("mockService1").AnyTimes()
	mockService1.EXPECT().Start(gomock.Any()).Times(1)

	mockService2 := NewMockService(ctrl)
	mockService2.EXPECT().Name().Return("mockService2").AnyTimes()
	mockService2.EXPECT().Start(gomock.Any()).Times(1)

	r := NewRunner(context.Background(), RunnerOptions{})

	err := r.RegisterService(mockService1, ServiceOptions{})
	assert.Nil(t, err, "Expected no error for registering service")

	err = r.RegisterService(mockService2, ServiceOptions{})
	assert.Nil(t, err, "Expected no error for registering service")

	status := r.Status()

	// drain the uptime for all services
	for k := range status.Services {
		x := status.Services[k]
		x.Uptime = 0
		status.Services[k] = x
	}

	assert.False(t, status.IsRunning, "Expected runner to not be running")
	assert.Equal(t, ServiceInfo{Status: ServiceStatusRegistered, Uptime: 0, Restarts: 0}, status.Services["mockService1"], "Expected mockService1 to be registered")
	assert.Equal(t, ServiceInfo{Status: ServiceStatusRegistered, Uptime: 0, Restarts: 0}, status.Services["mockService2"], "Expected mockService2 to be registered")

	// Start the runner
	go func() {
		if err := r.BootUp(); err != nil {
			t.Errorf("Error while booting up the runner: %v", err)
		}
	}()

	<-time.After(time.Second * 3)

	status = r.Status()

	// drain the uptime for all services
	for k := range status.Services {
		x := status.Services[k]
		x.Uptime = 0
		status.Services[k] = x
	}

	assert.True(t, status.IsRunning, "Expected runner to be running")
	// As we are using mock services, the status of the services will be exited once they are started.
	assert.Equal(t, ServiceInfo{Status: ServiceStatusExited, Uptime: 0, Restarts: 0}, status.Services["mockService1"], "Expected mockService1 to be exited")
	assert.Equal(t, ServiceInfo{Status: ServiceStatusExited, Uptime: 0, Restarts: 0}, status.Services["mockService2"], "Expected mockService2 to be exited")

	// Shutdown the runner
	r.Shutdown()

	status = r.Status()
	assert.False(t, status.IsRunning, "Expected runner to not be running after shutdown")
}
func TestRestartAllServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper1 := NewMockWrapper(ctrl)
	mockWrapper1.EXPECT().Status().Return(ServiceStatusRunning).Times(1)
	mockWrapper1.EXPECT().Stop().Times(1)
	mockWrapper1.EXPECT().Start().Times(1)

	mockWrapper2 := NewMockWrapper(ctrl)
	mockWrapper2.EXPECT().Status().Return(ServiceStatusRunning).Times(1)
	mockWrapper2.EXPECT().Stop().Times(1)
	mockWrapper2.EXPECT().Start().Times(1)

	mockWrapper3 := NewMockWrapper(ctrl)
	mockWrapper3.EXPECT().Status().Return(ServiceStatusStopped).Times(1)

	r := NewRunner(context.Background(), RunnerOptions{})
	ri := r.(*runner)

	ri.svc = map[string]Wrapper{
		"mockService1": mockWrapper1,
		"mockService2": mockWrapper2,
		"mockService3": mockWrapper3,
	}

	r.RestartAllServices()
	<-time.After(time.Second * 1)
}

func TestRestartService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper1 := NewMockWrapper(ctrl)
	mockWrapper1.EXPECT().Status().Return(ServiceStatusRunning).Times(1)
	mockWrapper1.EXPECT().Stop().Times(1)
	mockWrapper1.EXPECT().Start().Times(1)

	mockWrapper2 := NewMockWrapper(ctrl)

	r := NewRunner(context.Background(), RunnerOptions{})
	ri := r.(*runner)

	ri.svc = map[string]Wrapper{
		"mockService1": mockWrapper1,
		"mockService2": mockWrapper2,
	}

	_ = r.RestartService("mockService1")
	<-time.After(time.Second * 1)
}

func TestStopAllServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper1 := NewMockWrapper(ctrl)
	mockWrapper1.EXPECT().Status().Return(ServiceStatusRunning).Times(1)
	mockWrapper1.EXPECT().Stop().Times(1)

	mockWrapper2 := NewMockWrapper(ctrl)
	mockWrapper2.EXPECT().Status().Return(ServiceStatusRunning).Times(1)
	mockWrapper2.EXPECT().Stop().Times(1)

	mockWrapper3 := NewMockWrapper(ctrl)
	mockWrapper3.EXPECT().Status().Return(ServiceStatusStopped).Times(1)

	r := NewRunner(context.Background(), RunnerOptions{})
	ri := r.(*runner)

	ri.svc = map[string]Wrapper{
		"mockService1": mockWrapper1,
		"mockService2": mockWrapper2,
		"mockService3": mockWrapper3,
	}

	r.StopAllServices()
	<-time.After(time.Second * 1)
}

func TestStopService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper1 := NewMockWrapper(ctrl)
	mockWrapper1.EXPECT().Status().Return(ServiceStatusRunning).Times(1)
	mockWrapper1.EXPECT().Stop().Times(1)

	mockWrapper2 := NewMockWrapper(ctrl)

	mockWrapper3 := NewMockWrapper(ctrl)
	mockWrapper3.EXPECT().Status().Return(ServiceStatusStopped).Times(1)

	r := NewRunner(context.Background(), RunnerOptions{})
	ri := r.(*runner)

	ri.svc = map[string]Wrapper{
		"mockService1": mockWrapper1,
		"mockService2": mockWrapper2,
		"mockService3": mockWrapper3,
	}

	_ = r.StopService("mockService1", "mockService3")
	<-time.After(time.Second * 1)
}

func TestShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper1 := NewMockWrapper(ctrl)
	mockWrapper1.EXPECT().Status().Return(ServiceStatusRunning).Times(1)
	mockWrapper1.EXPECT().Stop().Times(1)

	mockWrapper2 := NewMockWrapper(ctrl)
	mockWrapper2.EXPECT().Status().Return(ServiceStatusStopped).Times(1)

	r := NewRunner(context.Background(), RunnerOptions{})
	ri := r.(*runner)

	ri.svc = map[string]Wrapper{
		"mockService1": mockWrapper1,
		"mockService2": mockWrapper2,
	}

	r.Shutdown()

	// Test when runner is not running after shutdown
	assert.False(t, r.IsRunning(), "Expected runner to not be running after shutdown")
}

func TestDeregisterService(t *testing.T) {
	r := NewRunner(context.Background(), RunnerOptions{})
	// Test deregistering a non-existent service
	err := r.DeregisterService("nonExistentService")
	assert.Equal(t, ErrDeregisterServiceNotFound, err, "Expected error for deregistering non-existent service")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	mockService.EXPECT().Name().Return("mockService").Times(1)

	err = r.RegisterService(mockService, ServiceOptions{})
	assert.Nil(t, err, "Expected no error for registering service")

	// Test deregistering a registered service
	err = r.DeregisterService("mockService")
	assert.Nil(t, err, "Expected no error for deregistering service")

	ri := r.(*runner)

	// Test if the service is deregistered
	_, ok := ri.svc["mockService"]
	assert.False(t, ok, "Expected service to be deregistered")
}

func TestDeregisterRunningService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper := NewMockWrapper(ctrl)
	mockWrapper.EXPECT().Status().Return(ServiceStatusRunning).Times(1)
	mockWrapper.EXPECT().Stop().Times(1)

	r := NewRunner(context.Background(), RunnerOptions{})
	ri := r.(*runner)

	ri.svc = map[string]Wrapper{
		"mockService": mockWrapper,
	}

	// Test deregistering a running service
	err := r.DeregisterService("mockService")
	assert.Nil(t, err, "Expected no error for deregistering running service")

	// Test if the service is deregistered
	_, ok := ri.svc["mockService"]
	assert.False(t, ok, "Expected running service to be deregistered")
}

func TestRegisterService(t *testing.T) {
	r := NewRunner(context.Background(), RunnerOptions{})

	// Test registering a nil service
	err := r.RegisterService(nil, ServiceOptions{})
	assert.Equal(t, ErrRegisterNilService, err, "Expected error for registering nil service")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	mockService.EXPECT().Name().Return("mockService").Return("mockService").Times(1)

	err = r.RegisterService(mockService, ServiceOptions{})
	assert.Nil(t, err, "Expected no error for registering service")

	mockService1 := NewMockService(ctrl)
	mockService1.EXPECT().Name().Return("mockService").Return("mockService").Times(1)

	// Test registering a service successfully
	err = r.RegisterService(mockService1, ServiceOptions{})
	assert.Equal(t, ErrRegisterServiceAlreadyExists, err, "Expected error for registering service that already exists")

	ri := r.(*runner)

	// Test if the service is registered
	_, ok := ri.svc["mockService"]
	assert.True(t, ok, "Expected service to be registered")
}
