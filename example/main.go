package main

import (
	"log"
	"time"

	"github.com/achu-1612/glcm"
	"github.com/achu-1612/glcm/hook"
	"github.com/achu-1612/glcm/service"
)

type serviceA struct{}

func (s *serviceA) Start(ctx service.Terminator) {
	for {
		<-time.After(time.Second * 2)

		select {
		case <-ctx.TermCh():
			return
		default:
			log.Println("serviceA is running ", time.Now())
		}
	}
}

func (s *serviceA) Status() string {
	return "status"
}

func (s *serviceA) Name() string {
	return "serviceA"
}

func preHook1(args ...interface{}) error {
	log.Println("pre-hook 1")

	return nil
}

func preHook2(args ...interface{}) error {
	log.Println("pre-hook 2")

	return nil
}

func postHook1(args ...interface{}) error {
	log.Println("post-hook 1")

	return nil
}

func postHook2(args ...interface{}) error {
	log.Println("post-hook 2")

	return nil
}

func main() {
	base := glcm.New()
	base.RegisterService(
		&serviceA{},
		service.WithPreHooks(
			hook.NewHandler("h1", preHook1, nil),
			hook.NewHandler("h2", preHook2, nil),
		),
		service.WithPostHooks(
			hook.NewHandler("h3", postHook1, nil),
			hook.NewHandler("h4", postHook2, nil),
		),
	)

	base.BootUp(nil)

	go func() {
		<-time.After(time.Second * 5)
		base.RestartService("serviceA")
	}()

	go func() {
		<-time.After(time.Second * 15)
		base.RestartAllServices()
	}()

	go func() {
		<-time.After(time.Second * 25)
		base.StopAllServices()
	}()

	base.Wait()
}
