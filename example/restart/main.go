package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/achu-1612/glcm"
	"github.com/achu-1612/glcm/example/hook"
	"github.com/achu-1612/glcm/example/service"
	svc "github.com/achu-1612/glcm/service"
)

func main() {
	base := glcm.NewRunner()

	if err := base.RegisterService(
		&service.ServiceA{},
		svc.WithPreHooks(
			hook.NewHookHandler("h1", "pre", "ServiceA"),
			hook.NewHookHandler("h2", "pre", "ServiceA"),
		),
		svc.WithPostHooks(
			hook.NewHookHandler("h3", "post", "ServiceA"),
			hook.NewHookHandler("h4", "post", "ServiceA"),
		),
	); err != nil {
		log.Fatal(err)
	}

	if err := base.RegisterService(
		&service.ServiceB{},
		svc.WithPreHooks(
			hook.NewHookHandler("h1", "pre", "ServiceB"),
			hook.NewHookHandler("h2", "pre", "ServiceB"),
		),
		svc.WithPostHooks(
			hook.NewHookHandler("h3", "post", "ServiceB"),
			hook.NewHookHandler("h4", "post", "ServiceB"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Create a thread which will restart ServiceA after 10 seconds.
	go func() {
		<-time.After(time.Second * 10)

		if err := base.RestartService("ServiceA"); err != nil {
			log.Printf("Error while restarting ServiceA: %v", err)
		}
	}()

	// Create a thread which will restart all the running service after 20 seconds.
	// But the baser runner will still be running.
	go func() {
		<-time.After(time.Second * 20)
		base.RestartAllServices()
	}()

	go func() {
		<-time.After(time.Second * 30)

		process, err := os.FindProcess(os.Getpid())
		if err != nil {
			fmt.Printf("Error finding process: %s\n", err)
			return
		}

		if err := process.Signal(syscall.SIGTERM); err != nil {
			fmt.Printf("Error sending termination signal: %s\n", err)
		}

	}()

	base.BootUp(context.TODO())
}
