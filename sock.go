package glcm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/achu-1612/glcm/log"
)

const (
	defaultSocketPath = "/tmp/glcm.sock"
)

// socketAction is a type for socket actions.
type socketAction string

// list of supported socket action commands
const (
	SocketActionStopAllServices socketAction = "stopAll"
	SocketActionStopService     socketAction = "stop"
	SocketActionRestartAll      socketAction = "restartAll"
	SocketActionRestartService  socketAction = "restart"
	SocketActionStatus          socketAction = "status"
	// socketActionBootup          socketAction = "bootup"
	// socketActionShutdown        socketAction = "shutdown"
)

// socketCommandStatus is a type for socket action status.
type socketCommandStatus string

// list of supported socket command status
const (
	Success socketCommandStatus = "success"
	Failure socketCommandStatus = "failure"
)

// SocketResponse represents the response from the socket.
type SocketResponse struct {
	Result interface{}         `json:"resulth"`
	Status socketCommandStatus `json:"status"`
}

// newSocket returns a new instance of the socket.
func newSocket(
	r Runner,
	socketPath string,
	allowedUIDs []int,
) (*socket, error) {
	s := &socket{
		r:          r,
		socketPath: socketPath,
		allowedUID: allowedUIDs,
	}

	return s, nil
}

// socket implements basic socket operations
type socket struct {
	r          Runner
	allowedUID []int
	socketPath string
}

// stopService stops the service with the given name(s).
func (s *socket) stopService(name ...string) *SocketResponse {
	if len(name) == 0 {
		return &SocketResponse{
			Result: "no service name provided",
			Status: Failure,
		}
	}

	if err := s.r.StopService(name...); err != nil {
		return &SocketResponse{
			Result: fmt.Sprintf("failed to stop service(s)- %v: %v", name, err),
			Status: Failure,
		}
	}

	return &SocketResponse{
		Result: fmt.Sprintf("service(s) stopped successfully: %v", name),
		Status: Success,
	}
}

// stopAllServices stops all the services.
func (s *socket) stopAllServices() *SocketResponse {
	s.r.StopAllServices()

	return &SocketResponse{
		Result: "All services stopped successfully",
		Status: Success,
	}
}

// restartService restarts the service with the given name(s).
func (s *socket) restartService(name ...string) *SocketResponse {
	if len(name) == 0 {
		return &SocketResponse{
			Result: "no service name provided",
			Status: Failure,
		}
	}

	if err := s.r.RestartService(name...); err != nil {
		return &SocketResponse{
			Result: fmt.Sprintf("failed to restart service(s)- %v: %v", name, err),
			Status: Failure,
		}
	}

	return &SocketResponse{
		Result: fmt.Sprintf("service(s) restarted successfully: %v", name),
		Status: Success,
	}
}

// restartAllServices restarts all the services.
func (s *socket) restartAllServices() *SocketResponse {
	s.r.RestartAllServices()

	return &SocketResponse{
		Result: "All services restarted successfully",
		Status: Success,
	}
}

// status returns the status of the runner along with the status of each registered service.
func (s *socket) status() *SocketResponse {
	return &SocketResponse{
		Result: s.r.Status(),
		Status: Success,
	}
}

// start starts the socket server.
// Note: This will be a blocking call. Once the quit signal is received,
// the server will be stopped and the docker file will be cleaned up
func (s *socket) start(done <-chan os.Signal) error {
	if s.socketPath == "" {
		s.socketPath = defaultSocketPath

		log.Warnf("No socket path provided, using default path: %s", s.socketPath)
	}

	// Remove existing socket/pipe if present
	if _, err := os.Stat(s.socketPath); err == nil {
		log.Warnf("Removing existing socket file: %s", s.socketPath)

		os.Remove(s.socketPath)
	}

	sock, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("creating socket listener: %w", err)
	}

	if err := os.Chmod(s.socketPath, 0700); err != nil {
		_ = sock.Close()

		return fmt.Errorf("setting file permission for the socket file: %w", err)
	}

	// on shutdown, close the socket and remove the socket file
	defer func() {
		if err := sock.Close(); err != nil {
			log.Errorf("close socket: %v", err)
		}

		if err := os.Remove(s.socketPath); err != nil {
			log.Errorf("remove socket file: %v", err)
		}
	}()

	log.Infof("Listening on %s. Permitted Access for user: %v", s.socketPath, s.allowedUID)

	for {
		select {
		case <-done:
			return nil
		default:
			conn, err := sock.Accept()
			if err != nil {
				return fmt.Errorf("accepting connection: %w", err)
			}

			if err := validateSocketAccess(conn, s.allowedUID); err != nil {
				log.Errorf("validate socket access: %v", err)
				_ = conn.Close()

				continue
			}

			// Not handling the command inside a go-routine.
			// This is to ensure that the commands are executed sequentially.
			if err := s.handler(conn); err != nil {
				log.Errorf("handle incomfing connection: %v", err)
			}

			_ = conn.Close()
		}
	}
}

func (s *socket) handler(conn net.Conn) error {
	reader := bufio.NewReader(conn)

	message, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	if message == "" {
		return fmt.Errorf("empty message received")
	}

	split := strings.Split(message, " ")

	command := split[0]
	args := split[1:]

	// remove /t /n from command and arguments
	command = strings.TrimSpace(command)
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}

	log.Infof("received command: %s with args: %v", command, args)

	var res *SocketResponse

	switch socketAction(command) {
	case SocketActionStopAllServices:
		res = s.stopAllServices()

	case SocketActionStopService:
		res = s.stopService(args...)

	case SocketActionRestartAll:
		res = s.restartAllServices()

	case SocketActionRestartService:
		res = s.restartService(args...)

	case SocketActionStatus:
		res = s.status()

	default:
		res = &SocketResponse{
			Result: fmt.Sprintf("unknown command: %s", command),
			Status: Failure,
		}
	}

	b, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("marshal response: %v", err)
	}

	if _, err := conn.Write(b); err != nil {
		return fmt.Errorf("write to socket response: %v", err)
	}

	return nil
}

// getUID returns the UID of the user who connected to the socket.
func getUID(conn *net.UnixConn) (int, error) {
	// Get file descriptor
	file, err := conn.File()
	if err != nil {
		return -1, fmt.Errorf("get file descriptor of the socket file: %w", err)
	}
	defer file.Close()

	ucred, err := syscall.GetsockoptUcred(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
	if err != nil {
		return -1, fmt.Errorf("getsockopt SO_PEERCRED: %w", err)
	}

	return int(ucred.Uid), nil
}

// validateSocketAccess checks if the user who connected to the socket is allowed to access it.
func validateSocketAccess(conn net.Conn, allowedUIDs []int) error {
	if len(allowedUIDs) == 0 {
		return nil // No restrictions
	}

	uid, err := getUID(conn.(*net.UnixConn))
	if err != nil {
		return err
	}

	for _, allowedUID := range allowedUIDs {
		if uid == allowedUID {
			return nil
		}
	}

	return fmt.Errorf("access denied for uid: %d", uid)
}
