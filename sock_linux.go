//go:build linux

package glcm

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/achu-1612/glcm/log"
)

func getUID(conn *net.UnixConn) (int, error) {
	// Get file descriptor
	file, err := conn.File()
	if err != nil {
		return -1, err
	}
	defer file.Close()

	ucred, err := syscall.GetsockoptUcred(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
	if err != nil {
		return -1, err
	}

	return int(ucred.Uid), nil
}

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

func getSocket(path string) (net.Listener, error) {
	if path == "" {
		path = "/tmp/glcm.sock"

		log.Warnf("No socket path provided, using default path: %s", path)
	}

	// Remove existing socket/pipe if present
	if _, err := os.Stat(path); err == nil {
		log.Warnf("Removing existing socket file: %s", path)

		os.Remove(path)
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	if err := os.Chmod(path, 0700); err != nil {
		_ = listener.Close()
		
		return nil, fmt.Errorf("failed to set file permissions: %w", err)
	}

	return listener, nil
}
