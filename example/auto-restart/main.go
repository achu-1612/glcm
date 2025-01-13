package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/achu-1612/glcm"
	"github.com/achu-1612/glcm/example/hook"
	"github.com/achu-1612/glcm/example/service"
)

func main() {
	base := glcm.NewRunner(context.Background(), glcm.RunnerOptions{})

	if err := base.RegisterService(
		&service.ServiceA{},
		glcm.ServiceOptions{
			PreHooks: []glcm.Hook{
				hook.NewHookHandler("h1", "pre", "ServiceA"),
				hook.NewHookHandler("h2", "pre", "ServiceA"),
			},
			PostHooks: []glcm.Hook{
				hook.NewHookHandler("h3", "post", "ServiceA"),
				hook.NewHookHandler("h4", "post", "ServiceA"),
			},
		},
	); err != nil {
		log.Fatal(err)
	}

	if err := base.RegisterService(
		&service.ServiceB{},
		glcm.ServiceOptions{
			PreHooks: []glcm.Hook{
				hook.NewHookHandler("h1", "pre", "ServiceB"),
				hook.NewHookHandler("h2", "pre", "ServiceB"),
			},
			PostHooks: []glcm.Hook{
				hook.NewHookHandler("h3", "post", "ServiceB"),
				hook.NewHookHandler("h4", "post", "ServiceB"),
			},
		},
	); err != nil {
		log.Fatal(err)
	}

	// Create a thread which will stop ServiceA after 10 seconds.
	go func() {
		<-time.After(time.Second * 10)

		if err := base.StopService("ServiceA"); err != nil {
			log.Printf("Error while stopping ServiceA: %v", err)
		}
	}()

	// Create a thread which will stop all the running service after 20 seconds.
	// But the base runner will still be running.
	go func() {
		<-time.After(time.Second * 20)
		base.StopAllServices()
	}()

	go func() {
		<-time.After(time.Second * 30)

		process, err := os.FindProcess(os.Getpid())
		if err != nil {
			log.Printf("Error finding process: %s\n", err)
			return
		}

		if err := process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Error sending termination signal: %s\n", err)
		}

	}()

	if err := base.BootUp(); err != nil {
		log.Fatalf("Error while booting up the runner: %v", err)
	}
}
