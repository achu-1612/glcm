package main

import (
	"log"
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

	// Create a thread which will restart ServiceA after 10 seconds.
	go func() {
		<-time.After(time.Second * 10)

		base.RestartService("ServiceA")
	}()

	// Create a thread which will restart all the running service after 20 seconds.
	// But the baser runner will still be running.
	go func() {
		<-time.After(time.Second * 20)
		base.RestartAllServices()
	}()

	base.BootUp(nil)
	base.Wait()
}
