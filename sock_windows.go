//go:build windows

package glcm

import (
	"net"

	"github.com/Microsoft/go-winio"
	"github.com/achu-1612/glcm/log"
)

func validateSocketAccess(_ net.Conn, _ []int) error {
	return nil
}

func getSocket(path string) (net.Listener, error) {
	if path == "" {
		path = "\\\\.\\pipe\\glcm"
		log.Warnf("No socket path provided, using default path: %s", path)
	}

	return winio.ListenPipe(path, nil)
}
