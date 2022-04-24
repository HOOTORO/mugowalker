package bot

import (
	"log"
	"time"
	"worker/adb"
	"worker/navi"
)

type Bot interface {
	New(*adb.Device)
}

type DayWalker interface {
	Walk(*navi.Place)
}

type NightStalker interface {
	Hunt()
}

type DayWalkerNightStalker interface {
	DayWalker
	NightStalker
}

type AfkBot struct {
	state *navi.Place
	*adb.Device
}

func New(dev *adb.Device, startPlace *navi.Place) (ab *AfkBot) {
	err := dev.Connect()
	if err != nil {
		log.Panicf("AfkBOT: can't connect, check adress.")
	}
	return &AfkBot{Device: dev, state: startPlace}
}

func (bot *AfkBot) Walk(e *navi.Place) {

	road := []navi.TPoint{e.Entry}
	for _, v := range road {
		if v.X != 0 {
			bot.Tap(v.X, v.Y)
		} else {
			bot.Back()
		}
		time.Sleep(5 * time.Second)
	}
}
