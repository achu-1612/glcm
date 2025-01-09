# glcm: Go Routine Lifecycle Management
```

      ___           ___       ___           ___
     /\  \         /\__\     /\  \         /\__\
    /::\  \       /:/  /    /::\  \       /::|  |
   /:/\:\  \     /:/  /    /:/\:\  \     /:|:|  |
  /:/  \:\  \   /:/  /    /:/  \:\  \   /:/|:|__|__
 /:/__/_\:\__\ /:/__/    /:/__/ \:\__\ /:/ |::::\__\
 \:\  /\ \/__/ \:\  \    \:\  \  \/__/ \/__/~~/:/  /
  \:\ \:\__\    \:\  \    \:\  \             /:/  /
   \:\/:/  /     \:\  \    \:\  \           /:/  /
    \::/  /       \:\__\    \:\__\         /:/  /
     \/__/         \/__/     \/__/         \/__/

```



`glcm` is a Go package designed to manage the complete lifecycle of goroutines, providing a structured approach to starting, stopping, and monitoring services within your Go applications.

## Features

- **Service Registration**: Register multiple services to be managed concurrently.
- **Lifecycle Management**: Control the startup and shutdown sequences of all registered services.
- **Service Control**: Individually start, stop, and restart services as needed.
- **Hooks Integration**: Define pre-run and post-run hooks for services to execute custom logic before starting or after stopping a service.

## Installation

To install the package, run:

```bash
go get github.com/achu-1612/glcm
```

## Usage

Here's how to use `glcm` in your project:

### 1. Import the package

```go
import "github.com/achu-1612/glcm"
```

### 2. Create a new runner

```go
runner := glcm.NewRunner()
```

### 3. Define a service

Implement the `Service` interface for your service. This interface requires the following methods:

- `Start(service.Terminator)`: Defines the startup logic for the service.
- `Status() string`: Status returns the status of the service.
- `Name() string`: Returns the name of the service.

Example:

```go
import "github.com/achu-1612/glcm/service"

type MyService struct{}

func (m *MyService) Start(ctx service.Terminator) {
    // Initialization logic here
    // Start should be a blocking call.
    // On closing of the ctx.TermCh() channel, the method should return.

    // example: 
    for {
		<-time.After(time.Second * 2)

		select {
		case <-ctx.TermCh():
			return
		default:
			log.Println("service is running ", time.Now())
	}
	
}

func (m *MyService) Name() error {
    return "MyService"
}

func (m *MyService) Status() string {
    return ""
}
```

### 4. Register the service

```go
err := runner.RegisterService(&MyService{})
if err != nil {
    // Handle error
}
```

### 5. Boot up the runner

```go
ctx := context.Background()

// BootUp boots up the runner. This will start all the registered services.
runner.BootUp(ctx)
```

### 6. Wait for the runner to stop

```go
// Wait waits for the runner to stop.
// Note: This is a blocking call. It is to be called after BootUp.
// Only a ShutDown() call will stop the runner.
// Even after all the registered services are stopped, runner would still be running.
runner.Wait()
```

### 7. Shutdown the runner

```go
// Shutdown shuts down the runner. This will stop all the registered services.
runner.Shutdown()
```

### 8. Stop service(s)

```go
// StopService stops the given list of services.
runner.StopService("MyService1", "MyService2")

// StopAllServices stops all the registered/running services.
runner.StopAllServices()
```

### 8. Restart service(s)

```go
// RestartService restarts the given list of services.
runner.RestartService("MyService1", "MyService2")

// RestartAllServices restarts all the registered/running services.
runner.RestartAllServices()
```

## Service Hooks

The `hook` package allows you to define hooks that execute before or after a service starts.

### 1. Import the hook package

```go
import "github.com/achu-1612/glcm/hook"
```

### 2. Create a new hook handler

```go
preRunHook := hook.NewHandler("PreRunHook", func(args ...interface{}) error {
    // Pre-run logic here
    return nil
})
```

### 3. Register the service with hooks

```go
err := runner.RegisterService(
    &MyService{}, 
    service.WithPreRunHook(preRunHook),
    service.WithPostRunHook(postRunHook)
    )
if err != nil {
    // Handle error
}
```

## Contributing

Contributions are welcome! Please submit issues and pull requests for any improvements or bug fixes.

## License

This project is licensed under the MIT License.

## TODO
- Support for timeout for go-routine shutdowns (if possible)
- Better error handling for the pre and post hooks for service
- Service dependency

