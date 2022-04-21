package bot

import (
	"afk/worker/adb"
	"time"

	// "afk/worker/esperia"
	"log"
)

type Bot interface {
	New(*adb.Device)
	Daily() string
	Arena() string
	Fight() string
	Place() []interface{}
}

type Walker interface {
	In() (int, int)
}

func (bot *AfkBot) Walkin(w Walker) *AfkBot {
	x, y := w.In()
	bot.dev.Tap(x, y)
	time.Sleep(5 * time.Second)
	return bot
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

func (ab *AfkBot) In(x, y int) {
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
