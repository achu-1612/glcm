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
- **Auto-Restart with Backoff**: Automatically restart services with optional exponential backoff.

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
ctx := context.Background()
runner := glcm.NewRunner(ctx, glcm.RunnerOptions{})
```

### 3. Define a service

Implement the `Service` interface for your service. This interface requires the following methods:

- `Start(Terminator)`: Defines the startup logic for the service.
- `Status() string`: Status returns the status of the service.
- `Name() string`: Returns the name of the service.

Example:

```go
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
err := runner.RegisterService(&MyService{}, glcm.ServiceOptions{})
if err != nil {
    // Handle error
}
```

### 5. Boot up the runner

```go
// BootUp boots up the runner. This will start all the registered services.
//Note: This is a blocking call. It is to be called after BootUp.
// Only a ShutDown() call will stop the runner.
// Even after all the registered services are stopped, runner would 
if err := runner.BootUp(); err != nil {
    log.Fatalf("Error while booting up the runner: %v", err)
}
```

### 6. Shutdown the runner

```go
// Shutdown shuts down the runner. This will stop all the registered services.
runner.Shutdown()
```

### 7. Stop service(s)

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

## Auto-Restart with Backoff
To enable auto-restart with backoff for a service, use the following options during service registration:
Note: The service will be restarted automatically only when `service.WithAutoRestart()` options is given while service registration and when the service exits automatically not by runner shutting it down.

```go
err := runner.RegisterService(
    &MyService{},
    glcm.ServiceOptions{
        AutoStart: glcm.AutoRestartOptions{
            Enabled:    true,
            Backoff:    true,
            MaxRetries: 5, // Optional: Set maximum retries
            BackOffExponent: 2, // Optional: Set backoff exponent
        },
    },
)
if err != nil {
    // Handle error
}
```

## Service Hooks

The `hook` package allows you to define hooks that execute before or after a service starts.

### 1. Create a new hook handler

```go
preRunHook := hook.NewHook("PreRunHook", func(args ...interface{}) error {
    // Pre-run logic here
    return nil
})
```

### 2. Register the service with hooks

```go
err := runner.RegisterService(
    &MyService{}, 
    glcm.ServiceOptions{
        PreHooks: []glcm.Hook{preRunHook},
        PostHooks: []glcm.Hook{postRunHook},
    },
)
if err != nil {
    // Handle error
}
```

## Socket Usage

`glcm` supports socket communication for both Windows and Linux platforms. This allows you to send commands to control the lifecycle of services.

### Windows

On Windows, `glcm` uses named pipes for socket communication.

### Linux

On Linux, `glcm` uses Unix domain sockets for socket communication.

### Allowed Messages and Actions

The following messages can be sent to the socket to control the services:

- `restart <service_name>`: restart the specified service
- `stop <service_name>`: stop the specified service.
- `restartAll <service_name>`: restart all the services.
- `stopAll <service_name>`: stop all the services.
- `list`: list all the service and their current status.

### Example Usage
```sh
echo "restartAll" | socat - UNIX-CONNECT:/tmp/glcm.sock
```

```go
// Example for Linux Unix domain socket communication
socketPath := "/tmp/glcm_socket"
conn, err := net.Dial("unix", socketPath)
if err != nil {
    log.Fatalf("Failed to connect to socket: %v", err)
}
defer conn.Close()

_, err = conn.Write([]byte("start MyService"))
if err != nil {
    log.Fatalf("Failed to send message: %v", err)
}
```

## Contributing

Contributions are welcome! Please submit issues and pull requests for any improvements or bug fixes.

## License

This project is licensed under the MIT License.

## TODO
- Support for Job with scheduling.
- Support for timeout for go-routine shutdowns (if possible).
- Better error handling for the pre and post hooks for service.
- Service dependency.