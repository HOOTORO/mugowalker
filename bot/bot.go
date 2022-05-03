package bot

import (
	"worker/adb"
	"worker/datman"
	"worker/navi"
	"worker/navi/cam"

	log "github.com/sirupsen/logrus"
)

type Bot interface {
	New(*adb.Device)
	TransferTo(*navi.Place)
	Work(job interface{})
}

//TODO: to complex. extract bot from esperia
type AfkBot struct {
	*adb.Device
	datman.DataManager
	navi.Navigator
}

func New(dev *adb.Device) (ab *AfkBot) {
	err := dev.Connect()
	if err != nil {
		log.Panicf("AfkBOT: can't connect, check adress.")
	}
	log.Infof("Connected to device: %v", dev)
	dman := datman.NewFM(dev.Name)
	return &AfkBot{Device: dev, DataManager: dman}
}

func (bot *AfkBot) TransferTo(e *navi.Place) {
	ss := cam.Capture(bot.Device, e.Name)
	ss.Area(bot.DataManager, e.Entry.X, e.Entry.Y, 60)
	bot.Walk(bot.Device, e)
}
