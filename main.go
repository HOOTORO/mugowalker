package main

import (
	"os"
	"time"

	"worker/adb"
	"worker/afk/activities"
	"worker/cfg"

	// "worker/afk/activities"
	"worker/bot"
	"worker/ui"

	"github.com/sirupsen/logrus"

	"worker/afk"
	c "worker/cfg"

	"github.com/fatih/color"
)

var (
	log                        *logrus.Logger
	user                       *c.Profile
	red, green, cyan, ylw, mgt func(...interface{}) string
)

func init() {
	user = c.ActiveUser()
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	ylw = color.New(color.FgHiYellow).SprintFunc()
	mgt = color.New(color.FgHiMagenta).SprintFunc()
}

func main() {
	log = c.Logger()
	fn := func(a string, b string) {

		log.Warnf("%v |>\n %v", mgt(a), b)
	}

	if len(os.Args) > 1 && os.Args[1] == "-t" {
		log.SetLevel(logrus.TraceLevel)

		color.HiRed("%v", "TEST RUN")

		_, e := adb.Connect("127.0.0.1:5556")
		if e != nil {
			log.Fatalf(c.Red("%v"), e.Error())
		}
		gw := afk.New(&ui.AppUser{})
		bb := bot.New(fn)
		bot := afk.NewArenaBot(bb, gw)

		a := bot.ScanText()
		_ = a
		b := activities.BoardsQuests(a.Result())
		log.Warn(c.Red(b))

		return
	}

	log.Warnf(c.Red("RUN BEGIN : %v"), time.Now())

	// err := ui.RunMainMenu(user)
	err := ui.RunUI(cfg.ActiveUser())
	if err != nil {
		log.Errorf("ERROROR: %v", err)
	}
}
