package adb

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

// DevState represents the last queried state of an Android device.
type DevState int

type Point struct {
	X int
	Y int
}

// binary: DeviceState#Offline = offline
// binary: DeviceState#Online = device
// binary: DeviceState#Unauthorized = unauthorized
const (
	Offline DevState = iota
	Online
	Unauthorized
)

// Device represents an attached Android device.
type Device struct {
	Serial   string
	DevState DevState
    TransportId int
	WmSize   Point
	abi      string
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

func Connect(host, port string) (*Device, error) {
	if adb == "" {
		return nil, ErrADBNotFound
	}
	devs, _ := Devices()
	_ = devs

	cmd := Cmd{Args: []string{"connect", fmt.Sprintf("%v:%v", host, port)}}

	if out, err := cmd.Call(); err == nil {
		dev := &Device{Serial: "", DevState: 1, abi: ""}
		dev.Resolution()
		dev.Abi()
		fmt.Printf("--> %v <--\n", out)
		return dev, nil
	} else {
		return nil, err
	}
}

func (d *Device) Abi() string {
	if d.abi == "" {
		res, err := d.Command("getprop", "ro.product.cpu.abi").Call()
		if err == nil {
			d.abi = strings.TrimSpace(res)
		}
	}
	return d.abi
}

// Resolution of connected wm
func (d *Device) Resolution() string {
	if d.WmSize == (Point{}) {
		res, err := d.Command("wm", "size").Call()
		if err == nil {
			r := regexp.MustCompile(`Physical size: (?P<x>\d+)x(?P<y>\d+)`)
			for k, v := range r.FindStringSubmatch(res) {
				switch k {
				case 1:
					(d.WmSize).Y, _ = strconv.Atoi(v)
					break
				case 2:
					(d.WmSize).X, _ = strconv.Atoi(v)
					break
				}
			}
		}
	}
	return d.WmSize.String()
}

// String returns a string representing the device.
func parseDevices(out string) ([]*Device, error) {
	a := strings.SplitAfter(out, "List of devices attached")
	if len(a) != 2 {
		return nil, errors.New("device list not returned")
	}
	lines := strings.Split(a[1], "\n")
	devices := make([]*Device, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		switch len(fields) {
		case 0:
			continue
//            TODO: since -l was added to "adb devices" remove later
		case 2:
			device := &Device{
				Serial:   fields[0],
				DevState: state(fields[1]),
			}
            device.Resolution()
            device.Abi()
			devices = append(devices, device)
		case 6:
            tid, _ := strconv.Atoi(strings.Trim(fields[5],"transport_id:"))
			device := &Device{
				Serial: fields[0],
                DevState: state(fields[1]),
                TransportId: tid,
			}
            device.Resolution()
            device.Abi()
            devices = append(devices, device)
		default:
			return nil, errors.New("could not parse device list")
		}
	}
	log.Infof("Availiable Devices:\n%v", devices)
	return devices, nil
}

func (d *Device) String() string {
	return fmt.Sprintf("Device<%s%s>[resolution:%s]", d.Serial, d.abi, d.WmSize)
}

func (p Point) String() string {
	return fmt.Sprintf("%dx%d", p.X, p.Y)
}

func state(str string) DevState {
    if str == "device"{
        return DevState(1)
    }
    return DevState(0)
}