package bot

import (
	"afk/worker/adb"
	"afk/worker/esperia"
	"log"
	"time"
)

type Bot interface {
	New(*adb.Device)
	Daily() string
	Arena() string
	Fight() string
	Place() []interface{}
}

type Walker interface {
	// Walk() (int, int)
	Walk() (esperia.TPoint, interface{})
}

func (bot *AfkBot) WalkIN(w Walker) interface{} {
	point, camp := w.Walk()
	bot.tap(point.X, point.Y)
	// x, y := w.WannaIn()
	// bot.WalkIn(x, y)
	time.Sleep(5 * time.Second)
	return camp
}

type AfkBot struct {
	state string
	dev   *adb.Device
}

func New(dev *adb.Device) (ab *AfkBot) {
	err := dev.Connect()
	if err != nil {
		log.Panicf("AfkBOT: can't connect, check adress.")
	}
	return &AfkBot{dev: dev}
}

func (ab *AfkBot) tap(x, y int) {
	ab.dev.Tap(x, y)
}

func (ab *AfkBot) Daily() string {
	panic("not implemented") // TODO: Implement
}

func (ab *AfkBot) Arena() string {
	panic("not implemented") // TODO: Implement
}

func (ab *AfkBot) Fight() string {
	panic("not implemented") // TODO: Implement

}

func (ab *AfkBot) Place() []interface{} {
	panic("not implemented") // TODO: Implement
}
