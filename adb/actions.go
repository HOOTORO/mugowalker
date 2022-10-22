/*
ADB Commands const
*/
package adb

import (
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
	shell, input, tap, swipe = "shell", "input", "tap", "swipe"
	screencap, screenrecord  = "screencap", "screenrecord"
)

// "adb shell input tap x,y"
func (d *Device) Tap(x, y string) {
	e := d.Command(input, tap, x, y).Run()
	if e != nil {
		fmt.Printf("\nerr:%v\nduring run:%v", e, "tap")
	}
	time.Sleep(3 * time.Second)
}

// adb shell input swipe <x1> <y1> <x2> <y2> [duration(ms)]
func (dev *Device) Swipe(x, y, x1, y1, td int) {
	xPos := strconv.Itoa(x)
	yPos := strconv.Itoa(y)
	x1Pos := strconv.Itoa(x1)
	y1Pos := strconv.Itoa(y1)
	duration := strconv.Itoa(td)
	e := dev.Command(swipe, xPos, yPos, x1Pos, y1Pos, duration).Run()
	if e != nil {
		fmt.Printf("\nerr:%vduring run:%v", e, "swipe")
	}
}

// "screencap -p /sdcard/ff.png"
func (d *Device) Screencap(f string) {
	// -p for png
	e := d.Command(screencap, remotedir+f).Run()
	if e != nil {
		fmt.Printf("\nrun: %v err: %v", "scr", e.Error())
	}
}
