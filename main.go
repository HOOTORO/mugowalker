package main

import (
	"worker/adb"
	"worker/bot"
	"worker/esperia"
	"worker/fshelp"

	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		name = "Bluestacks"
		host = "localhost"
		port = "51532"
	)
	log.SetLevel(log.DebugLevel)
	//fshelp.CreateFolder("_workdir")
	fshelp.SetWD("_workdir")
	blueStacks := adb.New(name, host, port)
	// blueStacks.Adb("kill-server")
	// blueStacks.Adb("start-server")
	afkbot := bot.New(blueStacks, esperia.Campain)

	some := esperia.ClownRealm
	// log.Printf("some: %v", some.Nparent(3))
	log.Printf("some: %v", some)

	//return
	afkbot.Walk(some)
	_ = some

	// not good enough
	// bot.Walkin(mapa)
	// Walkin(&mapa.Ranhorn).
	// 	Walkin(&mapa.Guild).
	// 	Walkin(&mapa.Hellscape).
	// 	Walkin(&mapa.Cursed)

	// err := blueStacks.Connect()
	// if err != nil {
	// 	fmt.Printf("Connections erroro: %v", err)
	// }
	//blueStacks.Tap(10, 1900)
	//blueStacks.Screencap("/sdcard/self.png")
	// pullres, err := blue.Pull("/sdcard/self.png")
	// fmt.Printf("PUll Result: %v, Error: %v", string(pullres), err)
	// das, err := blue.Shell(cmd)
	// fmt.Printf("Result: %v, Error: %v", string(das), err)
	// _ = das

	// imagg := fshelp.OpenImg("muco.png")
	// cctd := img.Concat(imagg, 20, 1400, 140, 1530)
	// fshelp.SaveAsPng("pause.png", cctd)
	// blueStacks.Tap(900, 1900)
	// blueStacks.Pull("ssss.gop")
}
