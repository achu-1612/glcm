package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/achu-1612/glcm"
)

/*
    This example demonstrates how to use resutl of the hook execution in the service.
	The ServiceA will implement the hook.Handler interface and the Service interface.
	The prehook execution result will be stored in the struct which implements the service interface.
*/

var _ glcm.Service = &ServiceA{}
var _ glcm.Hook = &ServiceA{}

// ServiceA will imeplement the Service interface as well as the hook Handler interface.
type ServiceA struct {
	PreHookResult string
}

func (s *ServiceA) Start(ctx glcm.Terminator) {
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
	base := glcm.NewRunner(context.Background(), glcm.RunnerOptions{})

	sA := &ServiceA{}

	if err := base.RegisterService(
		sA,
		glcm.ServiceOptions{
			PreHooks: []glcm.Hook{sA},
		},
	); err != nil {
		log.Fatal(err)
	}

	go func() {
		<-time.After(time.Second * 10)

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
