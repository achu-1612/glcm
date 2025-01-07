package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/achu-1612/glcm"
	"github.com/achu-1612/glcm/hook"
	svc "github.com/achu-1612/glcm/service"
)

/*
    This example demonstrates how to use resutl of the hook execution in the service.
	The ServiceA will implement the hook.Handler interface and the Service interface.
	The prehook execution result will be stored in the struct which implements the service interface.
*/

var _ svc.Service = &ServiceA{}
var _ hook.Handler = &ServiceA{}

// ServiceA will imeplement the Service interface as well as the hook Handler interface.
type ServiceA struct {
	PreHookResult string
}

func (s *ServiceA) Start(ctx svc.Terminator) {
	for {
		<-time.After(time.Second * 2)

		select {
		case <-ctx.TermCh():
			return

		default:
			log.Println("ServiceA is running ", time.Now())
			log.Println("ServiceA PreHookResult: ", s.PreHookResult)
		}
	}
}

func (s *ServiceA) Status() string {
	return ""
}

func (s *ServiceA) Name() string {
	return "ServiceA"
}

func (s *ServiceA) Execute() error {
	log.Println("ServiceA PreHook executed")
	s.PreHookResult = "ServiceA PreHook executed successfully"

	return nil
}

func main() {
	base := glcm.NewRunner()

	sA := &ServiceA{}

	if err := base.RegisterService(
		sA,
		svc.WithPreHooks(
			sA,
		),
	); err != nil {
		log.Fatal(err)
	}

	go func() {
		<-time.After(time.Second * 10)

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
