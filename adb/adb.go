package adb

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type EmulatorManager interface {
	AndroidDevice(string, string, string) *Device
	Connect(string)
	Shell(string) (string, error)
	Screencap(string) string
	ShareFolder() string
	Adb(string) ([]byte, error)
}

type adbd struct {
	*exec.Cmd
}

type Device struct {
	*adbd
	*Connection
	Name        string
	transportId string
	product     string
	model       string
}

type Connection struct {
	host   string
	port   string
	status bool
}

const (
	sharedFolder = "/mnt/windows/BstSharedFolder/"
	screenExt    = ".png"
)

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

// var gadb *adbd

func AndroidDevice(name, host, port string) (*Device, error) {
	a, _ := getAdb()
	// // TODO: Rework this. f devices() should ret []*Device
	// conn := &Connection{host: host, port: port, status: false}
	// dev := &Device{Name: name, Connection: conn}
	for _, v := range a.devices() {
		if v.host == host && v.port == port {
			return v, nil
		}
	}
	return nil, errors.New("Device not found")
}

func (d *Device) connect() error {
	if !d.Alive() {
		dest := d.host + ":" + d.port
		res, err := d.run(connect, dest)
		if err != nil || string(res)[:5] == "canno" {
			d.status = false
			return errors.New("Connection to host failed: " + dest)
		}
		d.status = true
	}
	return nil
}

func (dev *Device) Shell(args ...string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Shell: 1 subcommand required")
	}
	shellArgs := strings.Join(args, " ")
	res, err := dev.run(shell, shellArgs)
	return res, err
}

// Screenshot to PWD
func (dev *Device) Capture(name string) string {
	dev.Screencap(name)
	fpath := dev.PullScreen(name)
	return fpath
}

func (dev *Device) Screencap(scrname string) ([]byte, error) {
	if len(scrname) < 1 {
		return nil, errors.New("Screencap: filename required")
	}

	res, err := dev.Shell(screencap, sharedFolder+scrname+screenExt)
	return res, err
}

// made by screencap from sharedfolder
func (dev *Device) PullScreen(scrname string) string {
	filename := scrname + screenExt
	dev.Pull(sharedFolder + filename)
	return filename
}

func (dev *Device) Pull(fname string) ([]byte, error) {
	if len(fname) < 1 {
		return nil, errors.New("Pull: Filename required") // Specify path to file. Output optional, if not set - wd")
	}
	res, err := dev.run(pull, fname)
	return res, err
}

func (dev *Device) Input(args ...string) error {
	if len(args) < 2 {
		return errors.New("Input: min 2 args required, input source/command and args")
	}
	shellArgs := strings.Join(args, " ")
	_, err := dev.Shell(input, shellArgs)
	return err
}

func (dev *Device) GoForward(x, y int) {
	xPos := strconv.Itoa(x)
	yPos := strconv.Itoa(y)
	dev.Input(tap, xPos, yPos)
}

func (dev *Device) GoBack() {
	dev.Input(back)
}

// nargs: swipe <x1> <y1> <x2> <y2> [duration(ms)]
func (dev *Device) Swipe(x, y, x1, y1, td int) {
	xPos := strconv.Itoa(x)
	yPos := strconv.Itoa(y)
	x1Pos := strconv.Itoa(x1)
	y1Pos := strconv.Itoa(y1)
	duration := strconv.Itoa(td)
	dev.Input(swipe, xPos, yPos, x1Pos, y1Pos, duration)
}

func getAdb() (*adbd, error) {
	// fmt.Printf("Current Env: %v", os.Environ())
	// if gadb != nil {
	// 	return gadb, nil
	// } else {
	path, err := exec.LookPath(adb)
	if err != nil {
		fmt.Printf("didn't find '%v' executable\n", adb)
		return nil, errors.New("No adb for you today, my friend!")
	} else {
		fmt.Printf("'%v' executable is in '%s'\n", adb, path)

		return &adbd{exec.Command(adb, "")}, nil
	}
	//}
}

/*
	Run adb, first argument must be a adb subcommand

"connect", "localhost:1111"

"shell", "input", "tap", "100", "200"

"screencap", "- p ", "/sdcard/ff.png"

"pull", "/sdcard/ff.png"
*/
func (ad *adbd) run(args ...string) ([]byte, error) {
	drv, _ := getAdb()
	if len(args) < 1 {
		return nil, errors.New("Adb: 1 subcommand required")
	}
	drv.Args = append([]string{drv.Args[0]}, args...)
	res, err := drv.CombinedOutput()

	log.Tracef("Adb: CMD Output --> %s", res)

	return res, err
}

func (ad *adbd) devices() (devices []*Device) {
	b, e := ad.run("devices", "-l")
	if e != nil {
		log.Errorf("DevERR: %v", e.Error())
		return nil
	}

	s := strings.TrimPrefix(string(b), "List of devices attached\r\n")
	s = strings.TrimSuffix(s, "\r\n\r\n")
	strdevices := strings.Split(s, "\r\n")

	fmt.Printf("All Devices --> %v (len: %v)", strdevices, len(strdevices))
	for k, v := range strdevices {
		onedev := &Device{adbd: ad}

		r := regexp.MustCompile(`transport_id:(\d+)`)
		onedev.transportId = r.FindStringSubmatch(v)[1]
		r = regexp.MustCompile(`product:(\w+)`)
		onedev.product = r.FindStringSubmatch(v)[1]
		r = regexp.MustCompile(`(\d.\d.\d\):(\d+)`)
		onedev.host = r.FindStringSubmatch(v)[1]
		onedev.product = r.FindStringSubmatch(v)[1]
		fmt.Printf("Dev # %v -->> %v < ", k, onedev)
		// for nf, field := range onedev {
		// 	log.Tracef(" 		Param # %v -->> %v < ", nf, field)
		// }
	}
	return
}
