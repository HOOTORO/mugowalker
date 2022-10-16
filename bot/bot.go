package bot

import (
	"time"

	"worker/adb"
	"worker/datman"
	"worker/navi"

	log "github.com/sirupsen/logrus"
)

type Bot interface {
	New(*adb.Device)
	Walk(*navi.Location)
	Work(job interface{})
}

// TODO: to complex. extract bot from esperia
type AfkBot struct {
	*adb.Device
	datman.DataManager
	*navi.Navigator
	cnt, maxtry int
}

func New(dev *adb.Device, startLocation *navi.Location) (ab *AfkBot) {
	dman := datman.NewFM("dev")
	nav := &navi.Navigator{Liveloc: startLocation}
	return &AfkBot{Device: dev, DataManager: dman, Navigator: nav, maxtry: 3}
}

func (dw *AfkBot) Walk(dst *navi.Location) {
	log.Debugf("Let's take a walk from >>%v<< to >>%v<<", dw.Liveloc.Name, dst.Name)
	for dw.Liveloc != dst {
		dw.Step(dst.Nparent(dw.Liveloc.Depth + 1))
		log.Debugf("Made a #%d step to --> %v", dw.Liveloc.Depth, dw.Liveloc.Name)
	}
}

func (n *AfkBot) Step(target *navi.Location) {
	log.Debugf("Little step for bot (>>%v>> to <<%v>>", n.Liveloc.Name, target.Name)
	n.GoForward(target.Entry.X, target.Entry.Y)
	// give oit time to load
	time.Sleep(5 * time.Second)

	screen := n.OpenPng(n.Capture(target.Name))
	target.Etalons, _ = n.LocEtalons(target.Name)

	if len(target.Etalons) == 0 {
		n.Candidate(target, screen)
		log.Panicf("No etalons images for Location >>%v<<", target.Name)
	}
	// n.SaveImg("1.png", target.Etalons[0])
	// n.SaveImg("2.png", screen)
	n.cnt = 0
	for !target.IsLocation(screen) {
		time.Sleep(5 * time.Second)
		screen = n.OpenPng(n.Capture(target.Name))
		n.cnt++
		if n.cnt > n.maxtry {
			n.Unknownplace(target, screen, target.Entry)
			panic("WE FAILED MASTA")
		}
	}
	n.Liveloc = target
}

// func (b *AfkBot) InitEtalons(uimap map[string]*navi.Location) {

// 	for _, v := range uimap {
// 		for v.Depth > 2 {
// 			b.EtalonStep(v.Nparent(1))
// 			log.Debugf("# Let's take a walk from >>%v<< to >>%v<<", b.Liveloc.Name, v.Name)

// 			for b.Liveloc != v {
// 				b.EtalonStep(v.Nparent(b.Liveloc.Depth + 1))
// 				log.Debugf("### Made a #%d step to --> %v", b.Liveloc.Depth-1, b.Liveloc.Name)

// 			}
// 			for i := 1; i < v.Depth; i++ {
// 				b.GoBack()
// 				time.Sleep(3 * time.Second)
// 			}
// 		}
// 	}

// }

// func (bt *AfkBot) EtalonStep(t *navi.Location) {
// 	log.Debugf(">>>>>>>>>>>> ETALON STEP for bot (>>%v>> to <<%v>>", bt.Liveloc.Name, t.Name)
// 	captscr := t.Capture(bt)
// 	bt.Liveloc.Etalons = append(bt.Liveloc.Etalons, captscr)

// 	// bt.Liveloc.Areas[t.Name] = captscr.Area(t.Entry.X, t.Entry.Y, 60)
// 	bt.Candidate(bt.Liveloc, captscr)
// 	bt.GoForward(t.Entry.X, t.Entry.Y)
// 	time.Sleep(10 * time.Second)
// 	bt.Liveloc = t
// }
