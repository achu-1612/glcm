package glcm

import (
	"errors"
	"testing"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name      string
		handlerFn func(...interface{}) error
		args      []interface{}
		wantErr   bool
	}{
		{
			name: "Handler executes successfully",
			handlerFn: func(args ...interface{}) error {
				return nil
			},
			args:    []interface{}{"arg1", "arg2"},
			wantErr: false,
		},
		{
			name: "Handler returns error",
			handlerFn: func(args ...interface{}) error {
				return errors.New("handler error")
			},
			args:    []interface{}{"arg1", "arg2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHook(tt.name, tt.handlerFn, tt.args...)

			err := h.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if h.Name() != tt.name {
				t.Errorf("Name() = %v, want %v", h.Name(), tt.name)
			}
		})
	}
}
