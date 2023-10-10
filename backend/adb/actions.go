/*
ADB Commands const
*/
package adb

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type AdbArgs interface {
	Args(...string) []string
}

type File struct {
	Push func(r, l string) []string
	Pull func(r, l string) []string
}

type Script struct {
	State func() string
}

// device manipulation command arguments
const (
	input, tap, swipe       = "input", "tap", "swipe"
	screencap, screenrecord = "screencap", "screenrecord"
	keyevent, backbtn, home = "keyevent", "4", "3"
	enter, backspace        = "66", "67"
	am, start, kill         = "am", "start", "force-stop"
	ps, pipe, grep          = "pidof", "|", "grep"
)

// "adb shell input tap x,y"
func (d *Device) Tap(x, y string) error {
	e := d.Command(input, tap, x, y).Run()
	if e != nil {
		time.Sleep(10 * time.Second)
		d.Tap(x, y)
	}
	time.Sleep(1 * time.Second)
	return e
}

// adb shell input swipe <x1> <y1> <x2> <y2> [duration(ms)]
func (d *Device) Swipe(x, y, x1, y1, td int) error {
	xPos := fmt.Sprint(x)
	yPos := strconv.Itoa(y)
	x1Pos := strconv.Itoa(x1)
	y1Pos := strconv.Itoa(y1)
	duration := strconv.Itoa(td)

	e := d.Command(swipe, xPos, yPos, x1Pos, y1Pos, duration).Run()
	return e
}

// "screencap -p /sdcard/ff.png"
func (d *Device) Screencap() {
	// -p for png
	name := fmt.Sprintf("S_%v_%v.png", len(d.Files)+1, time.Now().UnixMilli())
	remotePath := remotedir + name
	e := d.Command(screencap, remotePath).Run()
	if e != nil {
		time.Sleep(10 * time.Second)
		d.Screencap()
	} else {
		d.Files = append(d.Files, &RemoteFile{name: name, path: remotePath})
	}
}

// adb shell input keyevent 4
func (d *Device) Back() {
	d.Command(input, keyevent, backbtn).Run()

}

func (d *Device) Home() error {
	e := d.Command(input, keyevent, home).Run()
	return e
}

func (d *Device) PS(appname string) string {
	cmd := d.Command(ps, appname)
	var b bytes.Buffer
	cmd.Stdout = &b
	e := cmd.Run()
	if e != nil {
		return ""
	}
	return b.String()
}

// Use:
// adb shell am start -n '<appPackageName>/<appActitivityName>'
// To get <appPackageName> run :
// adb shell pm list packages
// To get <appActitivityName> lunch app and run
// adb shell dumpsys window | grep -E 'mCurrentFocus'
func (d *Device) StartApp(appname string) error {
	cmd := d.Command(am, start, "-n", appname)
	e := cmd.Run()
	return e
}

func (d *Device) KillApp(appname string) error {
	cmd := d.Command(am, kill, appname)
	e := cmd.Run()
	return e
}
