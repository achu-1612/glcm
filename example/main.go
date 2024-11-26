package main

import (
	"log"
	"time"

	"github.com/achu-1612/glcm"
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

func main() {
	base := glcm.NewRunner()
	base.RegisterService(&serviceA{})
	base.BootUp(nil)
	base.Wait()
}
