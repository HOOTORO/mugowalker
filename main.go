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
	afkbot := bot.New(blueStacks)

	some := esperia.ClownRealm
	// log.Printf("some: %v", some.Nparent(3))
	log.Printf("some: %v", some)
	// afkbot.Screencap("so")
	// afkbot.Pull("so")
	// return
	afkbot.TransferTo(some)
	_ = some

}
