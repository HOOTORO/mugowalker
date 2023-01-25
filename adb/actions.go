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
	keyevent, backbtn, home  = "keyevent", "4", "3"
	enter, backspace         = "66", "67"
)

// "adb shell input tap x,y"
func (d *Device) Tap(x, y string) error {
	e := d.Command(input, tap, x, y).Run()
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v", e, "tap")
		// keepAliveVM()
		time.Sleep(10 * time.Second)
		d.Tap(x, y)
	}
	time.Sleep(1 * time.Second)
	return e
}

// adb shell input swipe <x1> <y1> <x2> <y2> [duration(ms)]
func (d *Device) Swipe(x, y, x1, y1, td int) error {
	xPos := strconv.Itoa(x)
	yPos := strconv.Itoa(y)
	x1Pos := strconv.Itoa(x1)
	y1Pos := strconv.Itoa(y1)
	duration := strconv.Itoa(td)
	e := d.Command(swipe, xPos, yPos, x1Pos, y1Pos, duration).Run()
	if e != nil {
		log.Errorf("\nerr:%vduring run:%v", e, "swipe")
	}
	return e
}

// "screencap -p /sdcard/ff.png"
func (d *Device) Screencap(f string) {
	// -p for png
	e := d.Command(screencap, remotedir+f).Run()
	if e != nil {
		fmt.Printf("\nrun: %v err: %v", "scr", e.Error())
		// keepAliveVM()
		time.Sleep(10 * time.Second)
		d.Screencap(f)
	}
}

// adb shell input keyevent 4
func (d *Device) Back() {
	e := d.Command(input, keyevent, backbtn).Run()
	if e != nil {
		fmt.Printf("\nrun: %v err: %v", "scr", e.Error())
	}
}

func (d *Device) Home() {
	e := d.Command(input, keyevent, home).Run()
	if e != nil {
		fmt.Printf("\nrun: %v err: %v", "scr", e.Error())
	}
}
