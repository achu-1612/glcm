package service

import (
	"sync"
	"testing"
	"time"

	"github.com/achu-1612/glcm/hook"
)

func TestWrapper_Start(t *testing.T) {
	tests := []struct {
		name      string
		preHooks  []hook.Handler
		postHooks []hook.Handler
	}{
		{
			name:      "No hooks",
			preHooks:  nil,
			postHooks: nil,
		},
		{
			name: "With pre-hooks",
			preHooks: []hook.Handler{
				&mockHook{name: "pre-hook-1"},
				&mockHook{name: "pre-hook-2"},
			},
			postHooks: nil,
		},
		{
			name:     "With post-hooks",
			preHooks: nil,
			postHooks: []hook.Handler{
				&mockHook{name: "post-hook-1"},
				&mockHook{name: "post-hook-2"},
			},
		},
		{
			name: "With pre and post-hooks",
			preHooks: []hook.Handler{
				&mockHook{name: "pre-hook-1"},
				&mockHook{name: "pre-hook-2"},
			},
			postHooks: []hook.Handler{
				&mockHook{name: "post-hook-1"},
				&mockHook{name: "post-hook-2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			svc := &mockService{}
			w := NewWrapper(svc, wg, WithPreHooks(tt.preHooks...), WithPostHooks(tt.postHooks...))

			go w.Start()

			<-time.After(time.Second)

			if !svc.started {
				t.Errorf("Service was not started")
			}

			if svc.stopped {
				t.Errorf("Service was stopped prematurely")
			}

			w.StopAndWait()

			if !svc.stopped {
				t.Errorf("Service was not stopped")
			}
		})
	}
}

type mockService struct {
	started bool
	stopped bool
}

func (m *mockService) Start(t Terminator) {
	m.started = true
	<-t.TermCh()
	m.started = false
	m.stopped = true
}

func (m *mockService) Name() string {
	return "mockService"
}

func (m *mockService) Status() string {
	return "mockServiceStatus"
}

type mockHook struct {
	name string
}

func (m *mockHook) Execute() error {
	return nil
}

func (m *mockHook) Name() string {
	return m.name
}