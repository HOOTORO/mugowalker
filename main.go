package main

import (
	"fmt"
	"os"
	"strconv"
	"worker/adb"
	"worker/afk"
	"worker/bot"
	"worker/cfg"

	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		game   = "afkarena"
		player = "Devitool"
	)
	f, _ := os.OpenFile("app.log", os.O_APPEND, 0644)
	log.SetLevel(log.TraceLevel)
	log.SetOutput(f)

	devs, e := adb.Devices()
	devno := len(devs)
	if e != nil {
		log.Panicf("\ndevs:%v\nerr:%v", devs, e)
	}
	if len(devs) > 1 {
		var desc string = "Choose, which one will be used by bot\n"
		for i, dev := range devs {
			desc += fmt.Sprintf("%v: Serial-> %v,   id-> %v,    resolution-> %v\n", i, dev.Serial, dev.TransportId, dev.Resolution)
		}
		devno, _ = strconv.Atoi(cfg.UserInput(desc, "0"))
	}

	gm := afk.New(game, player)
	devno = 0
	bt := bot.New(devs[devno], gm)
//        fmt.Printf("42: [%08b]", 42)
//		bt.MarkDone(afk.Wrizz)
//	    bt.MarkDone(afk.Loot)
//	    bt.MarkDone(afk.FastReward)
	//	bt.Tower(afk.Kings)
    bt.ZeroPosition()
    bt.GridTapOff(3,18,1)
	bt.Battle()
	//    bt.Peek()
	//    bt.GridTapOff(2,7,2)
	//	bt.Daily()
}
