package service

import (
	"log"
	"time"

	"github.com/achu-1612/glcm"
)

var _ glcm.Service = &ServiceA{}
var _ glcm.Service = &ServiceB{}
var _ glcm.Service = &ServiceC{}

type ServiceA struct{}

func (s *ServiceA) Start(ctx glcm.Terminator) {
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

func (s *ServiceB) Start(ctx glcm.Terminator) {
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

type ServiceC struct{}

func (s *ServiceC) Start(ctx glcm.Terminator) {
	for {
		<-time.After(time.Second * 5)

		select {
		case <-ctx.TermCh():
			return

		default:
			log.Println("ServiceC is exiting on its own  ", time.Now())

			return
		}
	}
}

func (s *ServiceC) Status() string {
	return ""
}

func (s *ServiceC) Name() string {
	return "ServiceC"
}
