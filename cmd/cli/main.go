package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/achu-1612/glcm"
	"github.com/achu-1612/glcm/cmd/cli/display"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Usage = "Command line tool to interact with the glcm socket"
	app.Author = "achu-1612"

	app.Commands = runnerCommands()

	_ = app.Run(os.Args)
}

// runnerCommands returns the list of commands related to runner.
func runnerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "stopAll",
			Usage:  "Stop all services",
			Flags:  []cli.Flag{getSocketFlag()},
			Action: stopAllAction,
		},
		{
			Name:  "stop",
			Usage: "stop given list of sevices",
			Flags: []cli.Flag{
				getSocketFlag(),
				cli.StringFlag{
					Name:     "services",
					Usage:    "List of services to stop",
					Required: true,
				},
			},
			Action: stopAction,
		},
		{
			Name:   "restartAll",
			Usage:  "Restart all services",
			Flags:  []cli.Flag{getSocketFlag()},
			Action: restartAllAction,
		},
		{
			Name:  "restart",
			Usage: "Restart given list of sevices",
			Flags: []cli.Flag{
				getSocketFlag(),
				cli.StringFlag{
					Name:     "services",
					Usage:    "List of services to stop",
					Required: true,
				},
			},
			Action: restartAction,
		},
		{
			Name:  "status",
			Usage: "Get the status of the runner and services",
			Flags: []cli.Flag{
				getSocketFlag(),
			},
			Action: statusAction,
		},
	}
}

// getSocketFlag returns the socket flag.
func getSocketFlag() cli.Flag {
	return cli.StringFlag{
		Name:  "socket",
		Usage: "Path to the socket",
		Value: "/tmp/glcm.sock",
	}

}

// validateServiceNameList validates the service name list.
func validateServiceNameList(s string) error {
	if len(strings.Split(s, ",")) == 0 {
		return fmt.Errorf("service name list cannot be empty")
	}

	return nil
}

// sendMessageOnSocket sends the message on the socket and returns the response.
func sendMessageOnSocket(socketPath, msg string) (*glcm.SocketResponse, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("connect to socket: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		return nil, fmt.Errorf("write to socket: %v", err)
	}

	sr := &glcm.SocketResponse{}

	if err := json.NewDecoder(conn).Decode(sr); err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}

	return sr, nil
}

// stopAllAction stops all the services.
func stopAllAction(c *cli.Context) {
	res, err := sendMessageOnSocket(
		c.String("socket"),
		fmt.Sprintf("%s\n", glcm.SocketActionStopAllServices),
	)
	if err != nil {
		display.Fatalf("stop all services: %v", err)
	}

	display.Printf(res)
}

// stopAction stops the given list of services.
func stopAction(c *cli.Context) {
	services := c.String("services")

	if err := validateServiceNameList(services); err != nil {
		display.Fatalf("validate service name list: %v", err)
	}

	res, err := sendMessageOnSocket(
		c.String("socket"),
		fmt.Sprintf("%s %s\n", glcm.SocketActionStopService, services),
	)
	if err != nil {
		display.Fatalf("stop given service(s): %v", err)
	}

	display.Printf(res)
}

// restartAllAction restarts all the services.
func restartAllAction(c *cli.Context) {
	res, err := sendMessageOnSocket(
		c.String("socket"),
		fmt.Sprintf("%s\n", glcm.SocketActionRestartAll),
	)
	if err != nil {
		display.Fatalf("restart all services: %v", err)
	}

	display.Printf(res)
}

// restartAction restarts the given list of services.
func restartAction(c *cli.Context) {
	services := c.String("services")

	if err := validateServiceNameList(services); err != nil {
		display.Fatalf("validate service name list: %v", err)
	}

	res, err := sendMessageOnSocket(
		c.String("socket"),
		fmt.Sprintf("%s %s\n", glcm.SocketActionRestartService, services),
	)
	if err != nil {
		display.Fatalf("restart given service(s): %v", err)
	}

	display.Printf(res)
}

// statusAction gets the status of the runner and services.
func statusAction(c *cli.Context) {
	res, err := sendMessageOnSocket(
		c.String("socket"),
		fmt.Sprintf("%s\n", glcm.SocketActionStatus),
	)
	if err != nil {
		display.Fatalf("stop all services: %v", err)
	}

	display.PrintStatus(res)
}
