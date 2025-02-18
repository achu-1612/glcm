package glcm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/achu-1612/glcm/log"
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
	Result interface{}         `json:"result"`
	Status socketCommandStatus `json:"status"`
}

// socket implements basic socket operations
type socket struct {
	r          Runner
	allowedUID []int
	socketPath string
	shutdownCh chan struct{}
	doneCh     chan struct{}
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
		shutdownCh: make(chan struct{}),
		doneCh:     make(chan struct{}),
	}

	return s, nil
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

// func (s *socket) shutdownRunner() *SocketResponse {
// 	s.r.Shutdown()
// 	return &SocketResponse{
// 		Result: "Shutting down the runner",
// 		Status: Success,
// 	}
// }

// shutdown stops the socket server.
// Note: This is a blocking call. It will wait for the server to stop.
func (s *socket) shutdown() {
	close(s.shutdownCh)
	<-s.doneCh // wait for the socket to stop
}

// start starts the socket server.
// Note: This will be a blocking call. Once the shutdown on the socket is called,
// the server will be stopped and the socket file will be removed.
func (s *socket) start() error {
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

	// on shutdown, close the socket and remove the socket file
	go func() {
		<-s.shutdownCh

		log.Info("Closing the socket listener")

		if err := sock.Close(); err != nil {
			log.Errorf("Close socket listener: %v", err)
		}

		log.Info("Removing socket file")

		if err := os.Remove(s.socketPath); err != nil {
			log.Errorf("Remove socket file: %v", err)
		}

		log.Info("Socket closed and file removed")

		close(s.doneCh) // notifiy the shutdown call that the socket is closed.
	}()

	if err := os.Chmod(s.socketPath, 0600); err != nil {
		return fmt.Errorf("setting file permission for the socket file: %v", err)
	}

	log.Infof("Listening on %s. Permitted Access for user: %v", s.socketPath, s.allowedUID)

	for {
		conn, err := sock.Accept()
		if err != nil {
			select {
			case <-s.shutdownCh:
				return nil
			default:
				log.Errorf("Accepting connection: %v", err)
				<-time.After(time.Second * 5)
				continue
			}
		}

		go func(conn net.Conn) {
			defer func() {
				if err := conn.Close(); err != nil {
					log.Errorf("Close connection: %v", err)
				}
			}()

			if err := validateSocketAccess(conn, s.allowedUID); err != nil {
				log.Errorf("Validate socket access: %v", err)
				return
			}

			if err := s.handler(conn); err != nil {
				log.Errorf("Handle incoming connection: %v", err)
			}
		}(conn)
	}
}

func (s *socket) handler(conn net.Conn) error {
	reader := bufio.NewReader(conn)

	message, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading message from socket: %w", err)
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

	log.Infof("Received command: %s with args: %v", command, args)

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
