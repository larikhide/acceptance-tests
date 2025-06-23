package integration

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const baseBinName = "temp-bestbinary"

// Build the program
//
// Run it (and wait for it listen on 8080)
//
// # Send an HTTP request to the server
//
// # Before the server has a chance to send an HTTP response, send SIGTERM
//
// See if we still get a response
func LauncTestProgram(port string) (cleanupFunc func(), sendInterrupt func() error, err error) {
	binName, err := buildBinary()
	if err != nil {
		return nil, nil, err
	}

	sendInterrupt, kill, err := runServer(binName, port)

	cleanupFunc = func() {
		if kill != nil {
			kill()
		}

		os.Remove(binName)
	}

	if err != nil {
		cleanupFunc()
		return nil, nil, err
	}

	return cleanupFunc, sendInterrupt, nil
}

func buildBinary() (string, error) {
	binName := randomString(10) + "-" + baseBinName

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		return "", fmt.Errorf("cannot build tool %s: %s", binName, err)
	}
	return binName, nil
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for k := range s {
		s[k] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

func runServer(binName string, port string) (sendInterrupt func() error, kill func(), err error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	cmdPath := filepath.Join(dir, binName)

	cmd := exec.Command(cmdPath)

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("cannot run temp converter: %s", err)
	}

	kill = func() {
		_ = cmd.Process.Kill()
	}

	sendInterrupt = func() error {
		return cmd.Process.Signal(syscall.SIGTERM)
	}

	err = waitForServerListening(port)

	return sendInterrupt, kill, err
}

func waitForServerListening(port string) error {
	for i := 0; i < 30; i++ {
		conn, err := net.Dial("tcp", net.JoinHostPort("localhost", port))
		if err != nil {
			return fmt.Errorf("cannot create conn :%v", err)
		}
		if conn != nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("nothing seems to be listening on port: %s", port)
}
