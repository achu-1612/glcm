package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/achu-1612/glcm"
	"github.com/achu-1612/glcm/example/service"
	svc "github.com/achu-1612/glcm/service"
)

func main() {
	base := glcm.NewRunner()

	if err := base.RegisterService(
		&service.ServiceC{},
		svc.WithAutoRestart(),
		svc.WithBackoff(),
		svc.WithMaxRetries(3),
	); err != nil {
		log.Fatal(err)
	}

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
	// base.Wait()
}
