package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/achu-1612/glcm"
	"github.com/achu-1612/glcm/example/service"
)

func main() {
	base := glcm.NewRunner(context.Background(), glcm.RunnerOptions{})

	if err := base.RegisterService(
		&service.ServiceC{},
		glcm.ServiceOptions{
			AutoStart: glcm.AutoRestartOptions{
				Enabled:    true,
				Backoff:    true,
				MaxRetries: 3,
			},
		},
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

	if err := base.BootUp(); err != nil {
		log.Fatalf("Error while booting up the runner: %v", err)
	}
}
