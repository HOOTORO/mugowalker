package adb

import (
	// "io"
	"errors"
	"fmt"
	"strconv"
	"strings"

	// "os"
	"os/exec"
)

type Adb interface {
	New(string) *Adb
	Connect(string)
	Shell(string) (string, error)
	Screencap(string) string
	ShareFolder() string
	Adb(string) ([]byte, error)
}

type Device struct {
	name    string
	adbpath string
	*Connection
}

type Connection struct {
	host   string
	port   string
	status bool
}

func (c *Connection) Alive() bool {
	return c.status
}

const (
	adb       string = "adb"
	shell            = "shell"
	devices          = "devices"
	connect          = "connect"
	screencap        = "screencap -p"
	pull             = "pull"
	input            = "input"
	tap              = "tap"
)

func New(name, host, port string) *Device {
	checkExeExists(adb)
	conn := &Connection{host: host, port: port, status: false}
	return &Device{name: name, Connection: conn}
}

func (d *Device) Connect() error {
	if !d.Alive() {
		dest := d.host + ":" + d.port
		_, err := d.Adb(connect, dest)
		if err != nil {
			d.status = false
			return errors.New("Connection: FAIL")
		}
		d.status = true
	}
	return nil
}

//Run Adb, first agrgument must be a adb subcommand
func (d *Device) Adb(args ...string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Adb: 1 subcommand required")
	}
	cmd := exec.Command(adb, args...)
	res, err := cmd.CombinedOutput()

	cmd = exec.Command("adb", "connect", "localhost:1111")
	exec.Command("adb", "shell", "input", "tap", "100", "200")
	exec.Command("adb", "screeencap", "- p ", "/sdcard/ff.pmng")

	return res, err
}

func (d *Device) Shell(args ...string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Shell: 1 subcommand required")
	}
	shellArgs := strings.Join(args, " ")
	res, err := d.Adb(shell, shellArgs)
	return res, err
}

func (d *Device) Screencap(args ...string) ([]byte, error) {
	if len(args) < 1 {
		//screencap -p /sdcard/screenshot.png
		return nil, errors.New("Screencap: filename(full path) required")
	}
	shellArgs := strings.Join(args, " ")
	res, err := d.Shell(screencap, shellArgs)
	return res, err
}

func (d *Device) Pull(args ...string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Pull: Specify path to file. Output optional, if not set - wd")
	}
	shellArgs := strings.Join(args, " ")
	res, err := d.Adb(pull, shellArgs)
	return res, err
}

func (d *Device) Input(args ...string) error {
	if len(args) < 2 {
		return errors.New("Input: min 2 args required, input source/command and args")
	}
	shellArgs := strings.Join(args, " ")
	_, err := d.Shell(input, shellArgs)
	return err
}

func (d *Device) Tap(x, y int) {
	xPos := strconv.Itoa(x)
	yPos := strconv.Itoa(y)
	d.Input(tap, xPos, yPos)
}
func (d *Device) ShareFolder() (docpath, picpath string) {
	//bluestacks shared folders
	docpath = "/mnt/windows/Documents/"
	picpath = "/mnt/windows/Pictures/"
	return
}

func checkExeExists(program string) {
	// fmt.Printf("Current Env: %v", os.Environ())
	path, err := exec.LookPath(program)
	if err != nil {
		fmt.Printf("didn't find '%v' executable\n", program)
	} else {
		fmt.Printf("'%v' executable is in '%s'\n", program, path)
	}
}
