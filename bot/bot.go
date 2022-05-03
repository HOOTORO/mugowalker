package bot

import (
	"worker/adb"
	"worker/datman"
	"worker/navi"
	// "worker/navi/cam"

	log "github.com/sirupsen/logrus"
)

type Bot interface {
	New(*adb.Device)
	Walk(*navi.Location)
	Work(job interface{})
}

//TODO: to complex. extract bot from esperia
type AfkBot struct {
	*adb.Device
	datman.DataManager
	*navi.Navigator
}

func New(dev *adb.Device, startLocation *navi.Location) (ab *AfkBot) {
	err := dev.Connect()
	if err != nil {
		log.Panicf("AfkBOT: can't connect, check adress.")
	}
	log.Infof("Connected to device: %v", dev)
	dman := datman.NewFM(dev.Name)
	nav := &navi.Navigator{Liveloc: startLocation}
	return &AfkBot{Device: dev, DataManager: dman, Navigator: nav}
}

func (dw *AfkBot) Walk(dst *navi.Location) {

	// ss := cam.Capture(bot.Device, e.Name)
	// ss.Area(bot.DataManager, e.Entry.X, e.Entry.Y, 60)
	// bot.Walk(bot.Device, e)
	if dw.Liveloc == nil {
		dw.Liveloc = dst.Nparent(1)
	}

	log.Debugf("Let's take a walk from >>%v<< to >>%v<<", dw.Liveloc.Name, dst.Name)

	for dw.Liveloc != dst {
		dw.Step(dw.Device, dst.Nparent(dw.Liveloc.Depth+1))
		log.Debugf("Made a #%d step to --> %v", dw.Liveloc.Depth-1, dw.Liveloc.Name)

	}

	// for n.curentPLace.Depth != 1 {
	// 	w.GoForward(target.Entry.X, target.Entry.Y)
	// 	time.Sleep(3 * time.Second)
	// 	n.curentPLace = n.curentPLace.Parent
	// }

	// for i := 0; i < target.Depth; i++ {
	// 	nextStep := target.Nparent(i + 1)

	// 	w.GoForward(nextStep.Entry.X, nextStep.Entry.Y)
	// 	time.Sleep(7 * time.Second)
	// 	n.curentPLace = nextStep
	// }
}
