package bot

import (
	"time"
	"worker/adb"
	"worker/fshelp"
	"worker/img"
	"worker/navi"

	log "github.com/sirupsen/logrus"
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
	log.Infof("Connected to device: %v", dev)
	return &AfkBot{Device: dev, state: startPlace}
}

func (bot *AfkBot) Walk(e *navi.Place) {
	log.Debugf("Let's take a walk from >>%v<< to >>%v<<", bot.state.Name, e.Name)
	for bot.state.Depth != 1 {
		bot.Back()
		time.Sleep(3 * time.Second)
		bot.state = bot.state.Parent
	}
	for i := 0; i < e.Depth; i++ {
		nextStep := e.Nparent(i + 1)
		log.Debugf("Make a #%d step to --> %v", i+1, nextStep.Name)
		//bot.Tap(nextStep.Entry.X, nextStep.Entry.Y)
		bot.ClickArea(nextStep.Entry.X, nextStep.Entry.Y, nextStep.Name)
		time.Sleep(7 * time.Second)
		bot.state = nextStep

	}
}

func (b *AfkBot) ClickArea(x, y int, name string) {
	fpa := "/sdcard/" + name + ".png"
	r, err := b.Screencap(fpa)
	log.Debugf("take a pic. %v, %v", r, err)
	r, err = b.Pull(fpa)
	log.Debugf("pull a pic. %v, %v", r, err)
	scrn := fshelp.OpenImg(name + ".png")

	clickplace := img.Concat(scrn, x-30, y-30, x+30, y+30)
	fshelp.SaveAsPng(name+"sample.png", clickplace)

	b.Tap(x, y)
}
