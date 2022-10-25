package adb

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ErrADBNotFound is returned when the ADB executable is not found.
var ErrADBNotFound = errors.New("ADB command not found on PATH")

// ErrDeviceUnauthorized is returned by ADB commands when the device has not
// authorized ADB debugging. Check the confirmation dialog on the device.
var ErrDeviceUnauthorized = errors.New("Device unauthorized")

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
	args := []string{}
	if c.Device != nil {
		args = append(args, "-s", c.Device.Serial)
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

func Parse(s string) []string {
	temp := strings.TrimPrefix(s, "List of devices attached\r\n")
	s = strings.TrimSuffix(temp, "\r\n\r\n")
	strdevices := strings.Split(s, "\r\n")

	// fmt.Printf("All Devices (len: %v) --> \n%v\n", len(strdevices), strings.Join(strdevices, "\n"))
	for _, v := range strdevices {

		// https://regex101.com/r/7YFfra/1
		// https://regex101.com/r/7YFfra/2

		r := regexp.
			MustCompile(
				`(?P<host>(?:\d{1,3}\.){3}\d{1,3}|` +
					`(?P<name>\w+))+` +
					`[-|:]?(?P<port>\d+)+` +
					`[^\r]+(?P<state>offline|bootloader|device)[\s]+` +
					`product:(?P<product>\w+)\s` +
					`model:(?P<model>\w+)\s` +
					`device:(?P<device>\w+)\s` +
					`transport_id:(?P<tid>\d)`)

		params := r.FindAllStringSubmatch(v, -1)
		devinfo := make(map[string]string, 0)

		for _, match := range params {
			// fmt.Printf("\nParams  #%v; val => %v", k, match)
			for ind, subName := range r.SubexpNames() {
				if subName != "" {
					devinfo[subName] = match[ind]
				}
			}
		}

	}
	return []string{s}
}

// const (
// 	state  string = "get-state"
// 	target string = "-t"
// )

// const (
// 	DEV_ID    = "tid"
// 	NAME      = "name"
// 	HOST      = "host"
// 	PORT      = "port"
// 	DEV_MODEL = "device"
// 	STATE     = "state"
// )

// func (d *Device) String() string {
// 	return fmt.Sprintf("\n# ADev: %v # < | > #Transport ID# [%v] < | >	HOSTNAME: < %v:%v >	<|>	STATE [%v]	>>> was <%v>	",
// 		d.devinfo[DEV_MODEL], d.devinfo[DEV_ID], d.devinfo[HOST], d.devinfo[PORT], d.devinfo[STATE], d.devinfo[STATE])
// }

// // Exec sh on remote Device
// func (dev *Device) sh(args ...string) ([]byte, error) {
// 	if len(args) < 1 {
// 		return nil, errors.New("Shell: 1 subcommand required")
// 	}
// 	shellArgs := strings.Join(args, " ")
// 	dev.attachAdb()
// 	res, err := dev.run(shell, shellArgs)
// 	return res, err
// }

// // Screenshot to PWD
// func (dev *Device) Capture(name string) string {
// 	dev.Screencap(name)
// 	fpath := dev.PullScreen(name)
// 	return fpath
// }

// func (dev *Device) Screencap(scrname string) ([]byte, error) {
// 	if len(scrname) < 1 {
// 		return nil, errors.New("Screencap: filename required")
// 	}

// 	res, err := dev.sh(screencap, sharedFolder+scrname+screenExt)
// 	return res, err
// }

// // made by screencap from sharedfolder
// func (dev *Device) PullScreen(scrname string) string {
// 	filename := scrname + screenExt
// 	dev.Pull(sharedFolder + filename)
// 	return filename
// }
