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
	log.SetLevel(log.InfoLevel)

	dev, e := adb.Connect(host, port)
	if e != nil {
		fmt.Printf("\ndev:%v\nerr:%v", dev, e)
	}

	mission := "C:/Users/maruk/vscode/afkarena/worker/bot/mission/task.yaml"

	runner := bot.New(dev)

	// err := runner.Mission(mission)
	// if err != nil {
	// 	log.Fatalf("MISSION GOES ERRRRRRRRRRRRRRRRRRRRRRRRR%v", err.Error())
	// }

	scn := &bot.Scenario{Path: mission, Pattern: "if"}
	err := runner.Snecnario(scn)
	if err != nil {
		log.Fatalf("MISSION GOES ERRRRRRRRRRRRRRRRRRRRRRRRR%v", err.Error())
	}
}
