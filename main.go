package main

import (
	"fmt"

	"worker/adb"
	"worker/bot"

	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		name = "Bluestacks"
		host = "127.0.0.1"
		port = "5615"
	)
	// TODO: scaling  adb shell wm size returns resolution
	log.SetLevel(log.TraceLevel)

	dev, e := adb.Connect(host, port)
	if e != nil {
		fmt.Printf("\ndev:%v\nerr:%v", dev, e)
	}

	// dev.Tap("100", "100")
	// dev.Screencap("/sdcard/gg.png")
	// e = dev.Pull("/sdcard/gg.png", ".")
	// if e != nil {
	// 	fmt.Printf("\nerr:%v\nduring run:%v", e.Error(), "pull")
	// }

	runner := &bot.Daywalker{Job: "peek", Device: dev}

	runner.Peek()
}
