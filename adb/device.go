package adb

import (
	"errors"
	"fmt"
	"strings"
)

// DevState represents the last queried state of an Android device.
type DevState int

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
	Serial string
	State  DevState
	abi    string
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
	cmd := Cmd{Args: []string{"devices"}}
	if out, err := cmd.Call(); err == nil {
		return parseDevices(out)
	} else {
		return nil, err
	}
}

func parseDevices(out string) ([]*Device, error) {
	a := strings.SplitAfter(out, "List of devices attached")
	if len(a) != 2 {
		return nil, errors.New("Device list not returned")
	}
	lines := strings.Split(a[1], "\n")
	devices := make([]*Device, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		switch len(fields) {
		case 0:
			continue
		case 2:
			state := DevState(0)
			if err := fields[1]; err != "nil" { // state.Parse(fields[1]); err != nil {
				return nil, nil // err
			}
			device := &Device{
				Serial: fields[0],
				State:  state,
			}
			devices = append(devices, device)
		default:
			return nil, errors.New("Could not parse device list")
		}
	}
	return devices, nil
}

func Connect(host, port string) (*Device, error) {
	if adb == "" {
		return nil, ErrADBNotFound
	}
	cmd := Cmd{Args: []string{"connect", fmt.Sprintf("%v:%v", host, port)}}

	if out, err := cmd.Call(); err == nil {
		fmt.Printf("Connect output --> %v", out)
		dev := &Device{Serial: "", State: 1, abi: ""}
		dev.Abi()
		return dev, nil
	} else {
		return nil, err
	}
}

// String returns a string representing the device.
func (d *Device) Abi() string {
	if d.abi == "" {
		res, err := d.Command("getprop", "ro.product.cpu.abi").Call()
		if err == nil {
			d.abi = strings.TrimSpace(res)
		}
	}
	return d.abi
}

// String returns a string representing the device.
func (d *Device) String() string {
	return fmt.Sprintf("Device<%s>", d.Serial)
}
