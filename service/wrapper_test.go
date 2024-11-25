package service

import (
	"testing"
)

func TestNewWrapper(t *testing.T) {
	service := &MockService{}
	context := &Context{}

	wrapper := NewWrapper(service, context)

	if wrapper.Service() != service {
		t.Errorf("expected service to be %v, got %v", service, wrapper.Service())
	}

	if wrapper.Context() != context {
		t.Errorf("expected context to be %v, got %v", context, wrapper.Context())
	}
}
