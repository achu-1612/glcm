package main

import (
	"context"
	"log"
	"time"

	"github.com/achu-1612/glcm"
	"github.com/achu-1612/glcm/hook"
	"github.com/achu-1612/glcm/service"
)

type serviceA struct{}

func (s *serviceA) Start(ctx *service.Context) {
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
	base := glcm.NewRunner()
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
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	base.BootUp(ctx)
	base.Wait()
}
