package glcm

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func setupSocket() (*net.UnixConn, func(), error) {
	socketPath := "/tmp/test.sock"
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		l.Close()
		os.Remove(socketPath)
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return conn.(*net.UnixConn), cleanup, nil
}

func TestGetUID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (*net.UnixConn, func(), error)
		wantUID int
		wantErr bool
	}{
		{
			name:    "Valid UID",
			setup:   setupSocket,
			wantUID: os.Getuid(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, cleanup, err := tt.setup()
			if cleanup != nil {
				defer cleanup()
			}
			if err != nil {
				t.Fatalf("setup() error = %v", err)
			}

			gotUID, err := getUID(conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUID != tt.wantUID {
				t.Errorf("getUID() = %v, want %v", gotUID, tt.wantUID)
			}
		})
	}
}

func TestValidateSocketAccess(t *testing.T) {
	tests := []struct {
		name        string
		allowedUIDs []int
		setup       func() (*net.UnixConn, func(), error)
		wantErr     bool
	}{
		{
			name:        "No restrictions",
			allowedUIDs: []int{},
			setup:       setupSocket,
			wantErr:     false,
		},
		{
			name:        "Allowed UID",
			allowedUIDs: []int{os.Getuid()},
			setup:       setupSocket,
			wantErr:     false,
		},
		{
			name:        "Denied UID",
			allowedUIDs: []int{99999}, // Assuming 99999 is not a valid UID
			setup:       setupSocket,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, cleanup, err := tt.setup()
			if cleanup != nil {
				defer cleanup()
			}
			if err != nil {
				t.Fatalf("setup() error = %v", err)
			}

			err = validateSocketAccess(conn, tt.allowedUIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSocketAccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestSocketRestartService(t *testing.T) {
	tests := []struct {
		name      string
		service   []string
		setupMock func(mockRunner *MockRunner)
		want      *SocketResponse
	}{
		{
			name:      "No service name provided",
			service:   []string{},
			setupMock: func(mockRunner *MockRunner) {},
			want: &SocketResponse{
				Result: "no service name provided",
				Status: Failure,
			},
		},
		{
			name:    "Service restart success",
			service: []string{"service1"},
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().RestartService("service1").Return(nil).Times(1)
			},
			want: &SocketResponse{
				Result: "service(s) restarted successfully: [service1]",
				Status: Success,
			},
		},
		{
			name:    "Service restart failure",
			service: []string{"service1"},
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().RestartService("service1").Return(fmt.Errorf("failed to restart")).Times(1)
			},
			want: &SocketResponse{
				Result: "failed to restart service(s)- [service1]: failed to restart",
				Status: Failure,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRunner := NewMockRunner(ctrl)

			tt.setupMock(mockRunner)

			s := &socket{
				r: mockRunner,
			}

			got := s.restartService(tt.service...)
			if got.Result != tt.want.Result || got.Status != tt.want.Status {
				t.Errorf("restartService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSocketStopService(t *testing.T) {
	tests := []struct {
		name      string
		service   []string
		setupMock func(mockRunner *MockRunner)
		want      *SocketResponse
	}{
		{
			name:      "No service name provided",
			service:   []string{},
			setupMock: func(mockRunner *MockRunner) {},
			want: &SocketResponse{
				Result: "no service name provided",
				Status: Failure,
			},
		},
		{
			name:    "Service stop success",
			service: []string{"service1"},
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().StopService("service1").Return(nil).Times(1)
			},
			want: &SocketResponse{
				Result: "service(s) stopped successfully: [service1]",
				Status: Success,
			},
		},
		{
			name:    "Service stop failure",
			service: []string{"service1"},
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().StopService("service1").Return(fmt.Errorf("failed to stop")).Times(1)
			},
			want: &SocketResponse{
				Result: "failed to stop service(s)- [service1]: failed to stop",
				Status: Failure,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRunner := NewMockRunner(ctrl)

			tt.setupMock(mockRunner)

			s := &socket{
				r: mockRunner,
			}

			got := s.stopService(tt.service...)
			if got.Result != tt.want.Result || got.Status != tt.want.Status {
				t.Errorf("stopService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSocketStopAllServices(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(mockRunner *MockRunner)
		want      *SocketResponse
	}{
		{
			name: "Stop all services success",
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().StopAllServices().Times(1)
			},
			want: &SocketResponse{
				Result: "All services stopped successfully",
				Status: Success,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRunner := NewMockRunner(ctrl)

			tt.setupMock(mockRunner)

			s := &socket{
				r: mockRunner,
			}

			got := s.stopAllServices()
			if got.Result != tt.want.Result || got.Status != tt.want.Status {
				t.Errorf("stopAllServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSocketRestartAllServices(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(mockRunner *MockRunner)
		want      *SocketResponse
	}{
		{
			name: "Restart all services success",
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().RestartAllServices().Times(1)
			},
			want: &SocketResponse{
				Result: "All services restarted successfully",
				Status: Success,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRunner := NewMockRunner(ctrl)

			tt.setupMock(mockRunner)

			s := &socket{
				r: mockRunner,
			}

			got := s.restartAllServices()
			if got.Result != tt.want.Result || got.Status != tt.want.Status {
				t.Errorf("restartAllServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSocketServiceStatus(t *testing.T) {
	status := &RunnerStatus{
		IsRunning: true,
		Services: map[string]ServiceStatus{
			"service1": ServiceStatusRunning,
		},
	}

	tests := []struct {
		name      string
		service   string
		setupMock func(mockRunner *MockRunner)
		want      *SocketResponse
	}{
		{
			name:    "Service status success",
			service: "service1",
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().Status().Return(status).Times(1)
			},
			want: &SocketResponse{
				Result: status,
				Status: Success,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRunner := NewMockRunner(ctrl)

			tt.setupMock(mockRunner)

			s := &socket{
				r: mockRunner,
			}

			got := s.status()
			if got.Result != tt.want.Result || got.Status != tt.want.Status {
				t.Errorf("serviceStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestSocketHandler(t *testing.T) {
	// status := &RunnerStatus{
	// 	IsRunning: true,
	// 	Services: map[string]ServiceStatus{
	// 		"service1": ServiceStatusRunning,
	// 	},
	// }

	tests := []struct {
		name      string
		command   string
		setupMock func(mockRunner *MockRunner)
		want      *SocketResponse
	}{
		{
			name:    "Stop all services",
			command: "stopAll\n",
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().StopAllServices().Times(1)
			},
			want: &SocketResponse{
				Result: "All services stopped successfully",
				Status: Success,
			},
		},
		{
			name:    "Stop specific service",
			command: "stop service1\n",
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().StopService("service1").Return(nil).Times(1)
			},
			want: &SocketResponse{
				Result: "service(s) stopped successfully: [service1]",
				Status: Success,
			},
		},
		{
			name:    "Restart all services",
			command: "restartAll\n",
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().RestartAllServices().Times(1)
			},
			want: &SocketResponse{
				Result: "All services restarted successfully",
				Status: Success,
			},
		},
		{
			name:    "Restart specific service",
			command: "restart service1\n",
			setupMock: func(mockRunner *MockRunner) {
				mockRunner.EXPECT().RestartService("service1").Return(nil).Times(1)
			},
			want: &SocketResponse{
				Result: "service(s) restarted successfully: [service1]",
				Status: Success,
			},
		},
		// {
		// 	name:    "Get status",
		// 	command: "status\n",
		// 	setupMock: func(mockRunner *MockRunner) {
		// 		mockRunner.EXPECT().Status().Return(status).Times(1)
		// 	},
		// 	want: &SocketResponse{
		// 		Result: status,
		// 		Status: Success,
		// 	},
		// },
		{
			name:      "Unknown command",
			command:   "unknown\n",
			setupMock: func(mockRunner *MockRunner) {},
			want: &SocketResponse{
				Result: "unknown command: unknown",
				Status: Failure,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRunner := NewMockRunner(ctrl)
			tt.setupMock(mockRunner)

			s := &socket{
				r: mockRunner,
			}

			conn := &mockConn{
				readBuffer:  strings.NewReader(tt.command),
				writeBuffer: &strings.Builder{},
			}

			err := s.handler(conn)
			if err != nil {
				t.Fatalf("handler() error = %v", err)
			}

			got := &SocketResponse{}
			if err = json.Unmarshal([]byte(conn.writeBuffer.String()), &got); err != nil {
				t.Fatalf("unmarshal response error = %v", err)
			}

			if got.Result != tt.want.Result || got.Status != tt.want.Status {
				t.Errorf("handler() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockConn struct {
	readBuffer  *strings.Reader
	writeBuffer *strings.Builder
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.UnixAddr{Name: "/tmp/mock.sock", Net: "unix"}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.UnixAddr{Name: "/tmp/mock.sock", Net: "unix"}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
