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

func TestRegisterService(t *testing.T) {
	r := NewRunner(context.Background(), RunnerOptions{})

	ri := r.(*runner)

	// Test registering a nil service
	err := r.RegisterService(nil, ServiceOptions{})
	assert.Equal(t, ErrRegisterNilService, err, "Expected error for registering nil service")

	// Test registering a service when runner is running
	ri.isRunning = true

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	mockService.EXPECT().Name().Return("mockService").Return("mockService").Times(2)

	err = r.RegisterService(mockService, ServiceOptions{})
	assert.Equal(t, ErrRunnerAlreadyRunning, err, "Expected error for registering service when runner is running")

	ri.isRunning = false

	// Test registering a service successfully
	err = r.RegisterService(mockService, ServiceOptions{})
	assert.Nil(t, err, "Expected no error for registering service")

	// Test registering a service that already exists
	err = r.RegisterService(mockService, ServiceOptions{})
	assert.Equal(t, ErrRegisterServiceAlreadyExists, err, "Expected error for registering service that already exists")
}

// func TestStatus(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockService1 := NewMockService(ctrl)
// 	mockService1.EXPECT().Name().Return("mockService1").AnyTimes()
// 	mockService1.EXPECT().Status().Return(ServiceStatusRunning).AnyTimes()

// 	mockService2 := NewMockService(ctrl)
// 	mockService2.EXPECT().Name().Return("mockService2").AnyTimes()
// 	mockService2.EXPECT().Status().Return(ServiceStatusStopped).AnyTimes()

// 	r := NewRunner(context.Background(), RunnerOptions{})

// 	err := r.RegisterService(mockService1, ServiceOptions{})
// 	assert.Nil(t, err, "Expected no error for registering service")

// 	err = r.RegisterService(mockService2, ServiceOptions{})
// 	assert.Nil(t, err, "Expected no error for registering service")

// 	status := r.Status()
// 	assert.False(t, status.IsRunning, "Expected runner to not be running")
// 	assert.Equal(t, ServiceStatusRunning, status.Services["mockService1"], "Expected mockService1 to be running")
// 	assert.Equal(t, ServiceStatusStopped, status.Services["mockService2"], "Expected mockService2 to be stopped")

// 	// Start the runner
// 	go func() {
// 		if err := r.BootUp(); err != nil {
// 			t.Errorf("Error while booting up the runner: %v", err)
// 		}
// 	}()

// 	<-time.After(time.Second * 10)

// 	status = r.Status()
// 	assert.True(t, status.IsRunning, "Expected runner to be running")
// 	assert.Equal(t, ServiceStatusRunning, status.Services["mockService1"], "Expected mockService1 to be running")
// 	assert.Equal(t, ServiceStatusStopped, status.Services["mockService2"], "Expected mockService2 to be stopped")

// 	// Shutdown the runner
// 	r.Shutdown()

// 	status = r.Status()
// 	assert.False(t, status.IsRunning, "Expected runner to not be running after shutdown")
// }
// func TestRegisterService(t *testing.T) {
// 	r := NewRunner(context.Background(), RunnerOptions{})

// 	// Test registering a nil service
// 	err := r.RegisterService(nil, ServiceOptions{})
// 	assert.Equal(t, ErrRegisterNilService, err, "Expected error for registering nil service")

// 	// Test registering a service when runner is running
// 	go func() {
// 		if err := r.BootUp(); err != nil {
// 			t.Errorf("Error while booting up the runner: %v", err)
// 		}
// 	}()

// 	<-time.After(time.Second * 10)

// 	mockService := &MockService{}
// 	err = r.RegisterService(mockService, ServiceOptions{})
// 	assert.Equal(t, ErrRunnerAlreadyRunning, err, "Expected error for registering service when runner is running")

// 	// Shutdown the runner
// 	r.Shutdown()

// 	// Test registering a service successfully
// 	err = r.RegisterService(mockService, ServiceOptions{})
// 	assert.Nil(t, err, "Expected no error for registering service")

// 	// Test registering a service that already exists
// 	err = r.RegisterService(mockService, ServiceOptions{})
// 	assert.Equal(t, ErrRegisterServiceAlreadyExists, err, "Expected error for registering service that already exists")
// }
