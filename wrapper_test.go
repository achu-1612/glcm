package glcm

import (
	"sync"
	"testing"
	"time"
)

func TestWrapper_Start(t *testing.T) {
	tests := []struct {
		name      string
		preHooks  []Hook
		postHooks []Hook
	}{
		{
			name:      "No hooks",
			preHooks:  nil,
			postHooks: nil,
		},
		{
			name: "With pre-hooks",
			preHooks: []Hook{
				&mockHook{name: "pre-hook-1"},
				&mockHook{name: "pre-hook-2"},
			},
			postHooks: nil,
		},
		{
			name:     "With post-hooks",
			preHooks: nil,
			postHooks: []Hook{
				&mockHook{name: "post-hook-1"},
				&mockHook{name: "post-hook-2"},
			},
		},
		{
			name: "With pre and post-hooks",
			preHooks: []Hook{
				&mockHook{name: "pre-hook-1"},
				&mockHook{name: "pre-hook-2"},
			},
			postHooks: []Hook{
				&mockHook{name: "post-hook-1"},
				&mockHook{name: "post-hook-2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			svc := &mockService{}
			w := NewWrapper(svc, wg, ServiceOptions{
				PreHooks:  tt.preHooks,
				PostHooks: tt.postHooks,
			})

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
