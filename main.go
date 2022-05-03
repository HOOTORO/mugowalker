package main

import (
	"worker/adb"
	"worker/bot"
	"worker/datman"
	"worker/esperia"

	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		name = "Bluestacks"
		host = "localhost"
		port = "62065"
	)
	log.SetLevel(log.DebugLevel)
	//fshelp.CreateFolder("_workdir")
	datman.SetWD("_workdir")
	blueStacks := adb.New(name, host, port)
	// blueStacks.Adb("kill-server")
	// blueStacks.Adb("start-server")
	UImap := esperia.UIMap()
	camp := UImap["Campain"]
	afkbot := bot.New(blueStacks, camp)

	some := UImap["ClowndRealm"]
	// log.Printf("some: %v", some.Nparent(3))
	log.Printf("some: %v", some)
	// afkbot.Screencap("so")
	// afkbot.Pull("so")
	// return
	afkbot.Walk(some)
	_ = some

}
