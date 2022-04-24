package bot

import (
	"afk/worker/adb"
	"log"
	"time"
)

type Bot interface {
	New(*adb.Device)
}

// type Walker interface {
// 	Walk() (int, int)
// }

type AfkBot struct {
	state Location
	dev   *adb.Device
}

type Location interface {
	Path(Location) []struct{ x, y int }
}

func (bot *AfkBot) DayWalker(e Location) {
	road := bot.state.Path(e)
	for _, v := range road {
		if v.x != 0 {
			bot.Walk(v.x, v.y)
		} else {
			bot.Back()
		}
		time.Sleep(5 * time.Second)
	}
}

func New(dev *adb.Device) (ab *AfkBot) {
	err := dev.Connect()
	if err != nil {
		log.Panicf("AfkBOT: can't connect, check adress.")
	}
	return &AfkBot{dev: dev}
}

func (ab *AfkBot) Walk(x, y int) {
	ab.dev.Tap(x, y)
}

func (ab *AfkBot) Back() {
	ab.dev.Back()
}
