package adb

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ErrADBNotFound is returned when the ADB executable is not found.
var ErrADBNotFound = errors.New("ADB command not found on PATH")

// ErrDeviceUnauthorized is returned by ADB commands when the device has not
// authorized ADB debugging. Check the confirmation dialog on the device.
var ErrDeviceUnauthorized = errors.New("device unauthorized")

// The path to the adb executable, or an empty string if the adb executable was
// not found.
var adb string

func init() {
	// Search for ADB using ANDROID_HOME
	if home := os.Getenv("ANDROID_HOME"); home != "" {
		path, err := filepath.Abs(filepath.Join(home, "platform-tools", "adb")) //+ maker.HostExecutableExtension)
		if err == nil {
			if _, err := os.Stat(path); err == nil {
				adb = path
				return
			}
		}
	}
	// Fallback to searching on PATH.
	if p, err := exec.LookPath("adb"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			adb = p
		}
	}
}

// Cmd represents a command that can be run on an Android device.
type Cmd struct {
	// Path is the path of the command to run on the device.
	//
	// If the string is empty, the command is treated as a ADB command for Device.
	Path string
	// Args holds the command line arguments to pass to the command.
	Args []string
	// The device this command should be run on. If nil, then any one of the
	// attached devices will execute the command.
	Device *Device
	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If Stdout and Stderr are the same writer, at most one
	// goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer
}

// Run starts the specified command and waits for it to complete.
// The returned error is nil if the command runs, has no problems copying
// stdout and stderr, and exits with a zero exit status.
func (c *Cmd) Run() error {
	var args []string
	if c.Device != nil {
		ID := strconv.Itoa(c.Device.TransportId)
		args = append(args, "-t", ID)
	}
	if c.Path != "" {
		args = append(args, "shell", c.Path)
	}
	args = append(args, c.Args...)
	cmd := exec.Command(adb, args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd.Run()
}

// Call starts the specified command and waits for it to complete, returning the
// all stdout as a string.
// The returned error is nil if the command runs, has no problems copying
// stdout and stderr, and exits with a zero exit status.
func (c *Cmd) Call() (string, error) {
	clone := *c // Don't change c's Stdout
	stdout := &bytes.Buffer{}
	if clone.Stdout != nil {
		clone.Stdout = io.MultiWriter(clone.Stdout, stdout)
	} else {
		clone.Stdout = stdout
	}
	stderr := &bytes.Buffer{}
	if clone.Stdout != nil {
		clone.Stderr = io.MultiWriter(clone.Stdout, stderr)
	} else {
		clone.Stderr = stderr
	}
	err := clone.Run()
	if err != nil && strings.Contains(stderr.String(), "error: device unauthorized.") {
		err = ErrDeviceUnauthorized
	}
	return stdout.String(), err
}
