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
		&service.ServiceC{},
		svc.WithPreHooks(
			hook.NewHookHandler("h1", "pre", "ServiceC"),
			hook.NewHookHandler("h2", "pre", "ServiceC"),
		),
		svc.WithPostHooks(
			hook.NewHookHandler("h3", "post", "ServiceC"),
			hook.NewHookHandler("h4", "post", "ServiceC"),
		),
		svc.WithAutoRestart(),
	); err != nil {
		log.Fatal(err)
	}

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

	if err := base.BootUp(context.TODO()); err != nil {
		log.Fatalf("Error while booting up the runner: %v", err)
	}
}
