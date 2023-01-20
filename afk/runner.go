package afk

import (
	a "worker/adb"
	"worker/afk/activities"
	"worker/bot"
	"worker/cfg"
)

var notifyUI func(a string, b string)

func Push(miss cfg.Mission, userprofile *cfg.Profile, out func(a string, b string)) {
	notifyUI = out
	if miss == cfg.PushCampain {
		ns := Nightstalker(userprofile)
		activities.Push(ns, out)

	}
}

func Nightstalker(cf *cfg.Profile) *Daywalker {
	dev, e := a.Connect(cf.DeviceSerial)
	if e != nil {
		log.Errorf("\ndeverr:%v", e)
	}
	gm := New(cf.User)
	d := bot.New(dev, notifyUI)

	return NewArenaBot(d, gm)
}
