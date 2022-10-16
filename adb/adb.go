package adb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type EmulatorManager interface {
	AndroidDevice(string, string, string) *Device
	Screencap(string) string
	ShareFolder() string
	Adb(string) ([]byte, error)
}

type adbd struct {
	*exec.Cmd
}

type Device struct {
	*Connection
	devinfo map[string]string
}
type Connection struct {
	*adbd
	status bool
}

const (
	sharedFolder = "/mnt/windows/BstSharedFolder/"
	screenExt    = ".png"
)

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

const (
	state  string = "get-state"
	target string = "-t"
)

const (
	DEV_ID    = "tid"
	NAME      = "name"
	HOST      = "host"
	PORT      = "port"
	DEV_MODEL = "device"
	STATE     = "state"
)

// var gadb *adbd

func AndroidDevice(name, host, port string) (dev *Device, e error) {
	a, _ := getAdb()
	// // TODO: Rework this. f devices() should ret []*Device
	// conn := &Connection{host: host, port: port, status: false}
	// dev := &Device{Name: name, Connection: conn}
	for _, v := range a.devices() {
		if v.devinfo[HOST] == host && v.devinfo[PORT] == port {

			v.connect()

			dev = v
		}
	}
	return nil, errors.New("Device not found")
}

func (dev *Device) param(k, v string) {
	dev.devinfo[k] = v
}

func (d *Device) String() string {
	return fmt.Sprintf("\n# ADev: %v # < | > #Transport ID# [%v] < | >	HOSTNAME: < %v:%v >	<|>	STATE [%v]	>>> was <%v>	",
		d.devinfo[DEV_MODEL], d.devinfo[DEV_ID], d.devinfo[HOST], d.devinfo[PORT], d.devinfo[STATE], d.devinfo[STATE])
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

		return &adbd{exec.Command(adb)}, nil
	}
	//}
}

func (d *Device) connect() {
	if d.devinfo[STATE] == "device" {
		// dest := d.devinfo[HOST] + ":" + d.devinfo[PORT]
		// res, err := d.cmd.run(connect, dest)
		// if err != nil || string(res)[:5] == "canno" {
		// 	d.status = false
		// 	return errors.New("Connection to host failed: " + dest)
		// }

		if d.adbd == nil {
			d.attachAdb()
		}
		if !d.status {
			d.status = true
		}

	}
}

func (d *Device) attachAdb() {
	d.Connection = &Connection{
		adbd:   &adbd{exec.Command(adb, target, d.devinfo[DEV_ID])},
		status: true,
	}
}

func (d *Device) state() string {
	state, _ := d.run(target, d.devinfo[DEV_ID], state)
	fmt.Printf("\nDev<%v> state was >>> %v | Current >>> %v <<<", d.devinfo[DEV_MODEL], d.devinfo[STATE], state)
	d.devinfo[STATE] = string(state)
	return string(state)
}

// Exec sh on remote Device
func (dev *Device) sh(args ...string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Shell: 1 subcommand required")
	}
	shellArgs := strings.Join(args, " ")
	dev.attachAdb()
	res, err := dev.run(shell, shellArgs)
	return res, err
}

/*
	Run adb, first argument must be a adb subcommand

"connect", "localhost:1111"

"shell", "input", "tap", "100", "200"

"screencap", "- p ", "/sdcard/ff.png"

"pull", "/sdcard/ff.png"
*/
func (ad *adbd) run(args ...string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("Adb: 1 subcommand required")
	}

	ad.Args = append(ad.Args, args...)

	var stdoutBuf, stderrBuf bytes.Buffer

	ad.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	ad.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	// stdout, err := ad.StdoutPipe()
	err := ad.Run()
	if err != nil {
		log.Fatalf("\nRun() failed with %s\n", err)
	}

	outStr, errStr := stdoutBuf.String(), stderrBuf.String()
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	// ?OFF DOCad

	/* 	stdout, err := ad.StdoutPipe()
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	if err := ad.Start(); err != nil {
	   		log.Fatal(err)
	   	}

	   	if err := ad.Wait(); err != nil {
	   		log.Fatal(err)
	   	}
	   	fmt.Printf("%s is %d years old\n", stdout, err) */

	fmt.Printf("\nRun:>>> <%v>\n	Output >>>> %s,\n	 errr >> %v", strings.Join(ad.Args, " "), outStr, errStr)

	return []byte(stdoutBuf.Bytes()), err
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

	// fmt.Printf("All Devices (len: %v) --> \n%v\n", len(strdevices), strings.Join(strdevices, "\n"))
	for _, v := range strdevices {
		// fmt.Printf("\nDev # %v -->>> %v <<< \n", k, v)

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
					// fmt.Printf("\n	<%v>:  	#>>> %v <<<", subName, match[ind])
				}
			}
		}

		onedev := &Device{devinfo: devinfo}
		onedev.attachAdb()
		devices = append(devices, onedev)
		fmt.Printf("%v", devices)
	}
	return
}

// func (d *Device) attachAdb(adb *adbd){
// 	if

// }
func (dev *Device) GoForward(x, y int) {
	xPos := strconv.Itoa(x)
	yPos := strconv.Itoa(y)
	dev.Input(tap, xPos, yPos)
}

func (dev *Device) GoBack() {
	dev.Input(back)
}

func (c *Connection) Alive() bool {
	state, _ := c.run(state)
	return string(state) == "device"
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

	res, err := dev.sh(screencap, sharedFolder+scrname+screenExt)
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
	res, err := dev.run(target, dev.devinfo[DEV_ID], pull, fname)
	return res, err
}

func (dev *Device) Input(args ...string) error {
	if len(args) < 2 {
		return errors.New("Input: min 2 args required, input source/command and args")
	}
	shellArgs := strings.Join(args, " ")
	_, err := dev.sh(input, shellArgs)
	return err
}
