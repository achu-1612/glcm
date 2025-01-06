package service

import (
	"log"
	"time"

	svc "github.com/achu-1612/glcm/service"
)

var _ svc.Service = &ServiceA{}
var _ svc.Service = &ServiceB{}

type ServiceA struct{}

func (s *ServiceA) Start(ctx svc.Terminator) {
	for {
		<-time.After(time.Second * 2)

		select {
		case <-ctx.TermCh():
			return

		default:
			log.Println("ServiceA is running ", time.Now())
		}
	}
}

func (s *ServiceA) Status() string {
	return ""
}

func (s *ServiceA) Name() string {
	return "ServiceA"
}

type ServiceB struct{}

func (s *ServiceB) Start(ctx svc.Terminator) {
	for {
		<-time.After(time.Second * 2)

		select {
		case <-ctx.TermCh():
			return

		default:
			log.Println("ServiceB is running ", time.Now())
		}
	}
}

func (s *ServiceB) Status() string {
	return ""
}

func (s *ServiceB) Name() string {
	return "ServiceB"
}
