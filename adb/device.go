package adb

import (
	"errors"
	"fmt"
	"image"
	"regexp"
	"strconv"
	"strings"

	"worker/cfg"
)

// DevState represents the last queried state of an Android device.

type DevState int

var (
	log          = cfg.Logger()
	ErrNoDevices = errors.New("attached devices not found")
)

// Point Offset:
// 0 -> full x*height
// 1 -> center point
// 2 -> xmax-1 height
type Point struct {
	image.Point
	Offset int
}

func (p Point) String() string {
	return fmt.Sprintf("%dx%d", p.X, p.Y)
}

// binary: DeviceState#Offline = offline
// binary: DeviceState#Online = device
// binary: DeviceState#Unauthorized = unauthorized
const (
	Offline DevState = iota
	Online
	Unauthorized
)

var strstates = [...]string{"Offline", "Online", "Unautorized"}

func (d DevState) String() string {
	return strstates[d]
}

// Device represents an attached Android device.
type Device struct {
	Serial      string
	DevState    DevState
	TransportId int
	Resolution  Point
	abi         string
}

func (d *Device) String() string {
	return fmt.Sprintf("Device<%s%s>[resolution:%s]", d.Serial, d.abi, d.Resolution)
}

// Command returns a new Cmd that will run the command with the specified name
// and arguments on this device.
func (d *Device) Command(path string, args ...string) *Cmd {
	return &Cmd{
		Path:   path,
		Args:   args,
		Device: d,
	}
}

// Devices returns the list of serial numbers of all the attached Android
// devices.
func Devices() ([]*Device, error) {
	if adb == "" {
		return nil, ErrADBNotFound
	}
	cmd := Cmd{Args: []string{"devices", "-l"}}
	if out, err := cmd.Call(); err == nil {
		return parseDevices(out)
	} else {
		return nil, err
	}
}

func Connect(hostport string) (*Device, error) {
	if adb == "" {
		return nil, ErrADBNotFound
	}
	//check existing connection
	devs, e := Devices()
	if e == nil {
		for _, d := range devs {
			if d.Serial == hostport {
				return d, nil
			}
		}
	}
	// serial := fmt.Sprintf("%v:%v", host, port)
	cmd := Cmd{Args: []string{"connect", hostport}}

	if out, err := cmd.Call(); err == nil && checkOut(out) {
		dev := &Device{Serial: hostport, DevState: Online}
		_ = resolution(dev)
		Abi(dev)
		log.Infof("--> %v", out)
		return dev, nil
	} else {
		return nil, errors.New(out)
	}
}

// String returns a string representing the device.
func parseDevices(out string) ([]*Device, error) {
	a := strings.SplitAfter(out, "List of devices attached")
	if len(a) != 2 {
		return nil, ErrNoDevices
	}
	lines := strings.Split(a[1], "\n")
	devices := make([]*Device, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		switch len(fields) {
		case 0, 8:
			continue
		case 6:
			tid, _ := strconv.Atoi(strings.Trim(fields[5], "transpo_id:"))
			device := &Device{
				Serial:      fields[0],
				DevState:    state(fields[1]),
				TransportId: tid,
			}
			_ = resolution(device)
			Abi(device)
			devices = append(devices, device)
		default:
			return nil, ErrNoDevices
		}
	}
	log.Infof("Availiable Devices:\n%v", devices)
	return devices, nil
}

func Abi(d *Device) string {
	if d.abi == "" {
		res, err := d.Command("getprop", "ro.product.cpu.abi").Call()
		if err == nil {
			d.abi = strings.TrimSpace(res)
		}
	}
	return d.abi
}

// resolution of connected wm
func resolution(d *Device) error {
	res, err := d.Command("wm", "size").Call()
	if err == nil {
		r := regexp.MustCompile(`Physical size: (?P<x>\d+)x(?P<y>\d+)`)
		for k, v := range r.FindStringSubmatch(res) {
			switch k {
			case 1:
				(d.Resolution).Y, err = strconv.Atoi(v)
			case 2:
				(d.Resolution).X, err = strconv.Atoi(v)
			}
		}
	}
	return err
}

func state(str string) DevState {
	if str == "device" {
		return Online
	}
	return Offline
}

func checkOut(str string) bool {
	return strings.Contains(str, "connected to") || strings.Contains(str, "already connected to")

}
