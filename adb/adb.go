package adb

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
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
	Name    string
	adbpath string
	*Connection
}

type Connection struct {
	host   string
	port   string
	status bool
}

const sharedFolder = "/mnt/windows/BstSharedFolder/"
const screenExt = ".png"

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
	back             = "keyevent 4"
	swipe            = "swipe"
)

func New(name, host, port string) *Device {
	checkExeExists(adb)
	conn := &Connection{host: host, port: port, status: false}
	return &Device{Name: name, Connection: conn}
}

func (d *Device) Connect() error {
	if !d.Alive() {
		dest := d.host + ":" + d.port
		res, err := d.Adb(connect, dest)
		if err != nil || string(res)[:5] == "canno" {
			d.status = false
			return errors.New("Connection to host failed: " + dest)
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

	// cmd = exec.Command("adb", "connect", "localhost:1111")
	// exec.Command("adb", "shell", "input", "tap", "100", "200")
	// exec.Command("adb", "screeencap", "- p ", "/sdcard/ff.png")
	//exec.Command("adb", "pull", "/sdcard/ff.png")

	log.Tracef("Adb: CMD Output --> %s", res)

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

func (d *Device) Screencap(scrname string) ([]byte, error) {
	if len(scrname) < 1 {
		return nil, errors.New("Screencap: filename required")
	}

	res, err := d.Shell(screencap, sharedFolder+scrname+screenExt)
	return res, err
}

// made by screencap from sharedfolder
func (d *Device) PullScreen(scrname string) string {
	filename := scrname + screenExt
	d.Pull(sharedFolder + filename)
	return filename
}

func (d *Device) Pull(fname string) ([]byte, error) {
	if len(fname) < 1 {
		return nil, errors.New("Pull: Filename required") //Specify path to file. Output optional, if not set - wd")
	}
	res, err := d.Adb(pull, fname)
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

func (d *Device) GoForward(x, y int) {
	xPos := strconv.Itoa(x)
	yPos := strconv.Itoa(y)
	d.Input(tap, xPos, yPos)
}

func (d *Device) GoBack() {
	d.Input(back)
}

// nargs: swipe <x1> <y1> <x2> <y2> [duration(ms)]
func (d *Device) Swipe(x, y, x1, y1, td int) {

	xPos := strconv.Itoa(x)
	yPos := strconv.Itoa(y)
	x1Pos := strconv.Itoa(x1)
	y1Pos := strconv.Itoa(y1)
	duration := strconv.Itoa(td)
	d.Input(swipe, xPos, yPos, x1Pos, y1Pos, duration)
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
