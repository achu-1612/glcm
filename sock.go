package glcm

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"encoding/json"

	"github.com/achu-1612/glcm/log"
)

// newSocket returns a new instance of the socket.
func newSocket(
	r Runner,
	socketPath string,
	allowedUIDs []int,
) (*socket, error) {
	// if the socket path is empty,
	// the default path will be used while creating the connection; no need to check

	s := &socket{
		r:               r,
		socketPath:      socketPath,
		allowedUID:      allowedUIDs,
		permittedAction: []string{},
	}

	return s, nil
}

// socket implmentes basic socker operations
type socket struct {
	r Runner
	// list of allowed user ids
	// Note: only valid for Liux systems.
	allowedUID      []int
	socketPath      string
	permittedAction []string
}

// stopService stops the service with the given name(s).
func (s *socket) stopService(name ...string) []byte {
	if err := s.r.StopService(name...); err != nil {
		return []byte(fmt.Sprintf("failed to stop service(s): %v", err))
	}

	return []byte("service(s) stopped successfully")
}

// stopAllServices stops all the services.
func (s *socket) stopAllServices() []byte {
	log.Infof("stopping all services")
	s.r.StopAllServices()

	return []byte("all services stopped successfully")
}

// restartService restarts the service with the given name(s).
func (s *socket) restartService(name ...string) []byte {
	if err := s.r.RestartService(name...); err != nil {
		return []byte(fmt.Sprintf("failed to restart service(s): %v", err))
	}

	return []byte("service(s) restarted successfully")
}

// restartAllServices restarts all the services.
func (s *socket) restartAllServices() []byte {
	s.r.RestartAllServices()

	return []byte("all services restarted successfully")
}

// listServices lists all the services and their current status.
func (s *socket) listServices() []byte {
	b, err := json.Marshal(s.r.ListServices())
	if err != nil {
		return []byte(fmt.Sprintf("failed to list services: json encoding of the values: %v", err))
	}

	return b
}

// start starts the socket server.
// Note: This will be a blocking call. Once the quit signal is received,
// the server will be stopped and the docker file will be cleaned up
func (s *socket) start(done <-chan os.Signal) error {
	// based on the operation systm we are running this on, get the appropriate socket
	sock, err := getSocket(s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	// on shutdown, close the socket and remove the socket file
	defer func() {
		if err := sock.Close(); err != nil {
			log.Errorf("failed to close socket: %v", err)
		}

		if err := os.Remove(s.socketPath); err != nil {
			log.Errorf("failed to remove socket file: %v", err)
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
				return fmt.Errorf("failed to accept connection: %w", err)
			}

			if err := validateSocketAccess(conn, s.allowedUID); err != nil {
				log.Errorf("failed to validate socket access: %v", err)
				_ = conn.Close()

				continue
			}

			// Not handling the command inside a go-routine.
			// This is to ensure that the commands are executed sequentially.
			if err := s.handler(conn); err != nil {
				log.Errorf("failed to handle connection: %v", err)
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

	var res []byte

	switch command {
	case "stopAll":
		res = s.stopAllServices()
	case "stop":
		if len(args) == 0 {
			return fmt.Errorf("no service name provided")
		}

		res = s.stopService(args...)
	case "restartAll":
		res = s.restartAllServices()

	case "restart":
		if len(args) == 0 {
			return fmt.Errorf("no service name provided")
		}

		res = s.restartService(args...)
	case "list":
		res = s.listServices()
	}

	if _, err := conn.Write(res); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}
