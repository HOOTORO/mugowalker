package afk

import (
	a "worker/adb"
	"worker/bot"
	"worker/cfg"
)

func Push(miss cfg.Mission) {
	if miss == cfg.PushCampain {

	}
}

func Nightstalker(cf *cfg.Profile, out func(a string, b string)) *Daywalker {
	dev, e := a.Connect(cf.DeviceSerial)
	if e != nil {
		log.Errorf("\ndeverr:%v", e)
	}
	gm := New(cf.User)
	d := bot.New(dev, out)

	return NewArenaBot(d, gm)
}
