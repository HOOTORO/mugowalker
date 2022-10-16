package main

import (
	"fmt"

	"worker/adb"

	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		name = "Bluestacks"
		host = "127.0.0.1"
		port = "5555"
	)
	// TODO: scaling  adb shell wm size returns resolution
	log.SetLevel(log.TraceLevel)
	bs, _ := adb.AndroidDevice(name, host, port)
	// bs.Conn5ect()
	if bs != nil {
		bs.Capture("EZX")
	} else {
		fmt.Printf("\nALIVE STATUS: %v ", bs)
	}

	// log.Infof("Connection status --> %b", bs.Alive())
}
