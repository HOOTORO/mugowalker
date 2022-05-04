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

	if dw.Liveloc == nil {
		dw.Liveloc = dst.Nparent(1)
	}

	log.Debugf("Let's take a walk from >>%v<< to >>%v<<", dw.Liveloc.Name, dst.Name)

	for dw.Liveloc != dst {
		dw.Step(dw.Device, dst.Nparent(dw.Liveloc.Depth+1))
		log.Debugf("Made a #%d step to --> %v", dw.Liveloc.Depth-1, dw.Liveloc.Name)

	}

}

func (b *AfkBot) InitEtalons(uimap map[string]*navi.Location) {

	for _, v := range uimap {
		// if b.Liveloc == nil {
		// 	b.Liveloc = v.Nparent(1)
		// }
		b.EtalonStep(b.Device, b.DataManager, v.Nparent(1))
		log.Debugf("Let's take a walk from >>%v<< to >>%v<<", b.Liveloc.Name, v.Name)

		for b.Liveloc != v {
			b.EtalonStep(b.Device, b.DataManager, v.Nparent(b.Liveloc.Depth+1))
			log.Debugf("Made a #%d step to --> %v", b.Liveloc.Depth-1, b.Liveloc.Name)

		}
		for i := 1; i < v.Depth; i++ {
			b.GoBack()
		}

	}

}
