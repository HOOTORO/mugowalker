package main

import (
	"os"
	"time"

	"worker/adb"
	"worker/afk/activities"
	// "worker/afk/activities"
	"worker/bot"
	"worker/ui"

	"github.com/sirupsen/logrus"

	"worker/afk"
	"worker/cfg"

	"github.com/fatih/color"
)

var (
	log                        *logrus.Logger
	user                       *cfg.Profile
	red, green, cyan, ylw, mgt func(...interface{}) string
)

func init() {
	user = cfg.ActiveUser()
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	ylw = color.New(color.FgHiYellow).SprintFunc()
	mgt = color.New(color.FgHiMagenta).SprintFunc()
}

func main() {
	log = cfg.Logger()
	fn := func(a string, b string) {
		log.Warnf("%v |>\n %v", mgt(a), b)
	}

	if len(os.Args) > 1 && os.Args[1] == "-t" {
		log.SetLevel(logrus.TraceLevel)

		color.HiRed("%v", "TEST RUN")

		d, e := adb.Connect("127.0.0.1:5556")
		if e != nil {
			log.Fatalf(red("%v"), e.Error())
		}
		gw := afk.New(user.User)
		bb := bot.New(d, fn)
		bot := afk.NewArenaBot(bb, gw)

		a := bot.ScanText()
		_ = a
		b := activities.BoardsQuests(a)
		log.Warnf(red(b))

		return
	}

	log.Warnf(red("RUN BEGIN : %v"), time.Now())

	err := ui.RunMainMenu(user)
	if err != nil {
		log.Errorf("ERROROR: %v", err)
	}
}
