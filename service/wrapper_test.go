package service

import (
	"testing"
)

func TestNewWrapper(t *testing.T) {
	mockService := &MockService{}
	option := func(ctx *Context) {
		ctx.terminationChan = make(chan struct{})
	}

	wrapper := NewWrapper(mockService, option)

	if wrapper.Service() != mockService {
		t.Errorf("expected service %v, got %v", mockService, wrapper.Service())
	}

	if wrapper.Context() == nil {
		t.Error("expected context to be non-nil")
	}

	if wrapper.Context().terminationChan == nil {
		t.Error("expected terminationChan to be non-nil")
	}
}
