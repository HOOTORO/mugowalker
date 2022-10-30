package main

import (
	"fmt"

	"worker/adb"
	"worker/bot"
	"worker/game"

	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		name = "Bluestacks"
		host = "127.0.0.1"
		port = "5615"
	)
	// TODO: scaling  adb shell wm size returns resolution
	log.SetLevel(log.InfoLevel)

	dev, e := adb.Connect(host, port)
	if e != nil {
		fmt.Printf("\ndev:%v\nerr:%v", dev, e)
	}

	gamecfg := "C:/Users/maruk/vscode/afkarena/worker/bot/cfg/config.yaml"

	namaewa := "Devitool"
	bt := bot.New(dev, namaewa)
	gm := game.New(gamecfg, "afkarena", bt)

	// err := gm.Push()
	err := gm.Daily()
	if err != nil {
		log.Fatalf("MISSION GOES ERRRRRRRRRRRRRRRRRRRRRRRRR%v", err.Error())
	}
}
