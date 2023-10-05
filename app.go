package main

import (
	"context"
	"fmt"
	"mugowalker/backend"
	// "mugowalker/backend/adb"
	// "mugowalker/backend/afk"
	// "mugowalker/backend/bot"
	// "mugowalker/backend/pilot"
	// "mugowalker/backend/settings"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context
	back *backend.Config
	// configuration *settings.Settings
	// pilot *pilot.Pilot
	// bot   *afk.Daywalker
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// emitter := func(s1, s2 string) { a.sendMessage(s1, s2) }
	// a.back = backend.NewConfig()
	// a.back.WailsInit(emitter)

	// a.pilot = pilot.Default()
	// a.configuration.Log(settings.WARN, fmt.Sprintf("1:%v", a.pilot))
	// gw := afk.New(a.pilot)
	// bb := bot.New(emitter, a.back.Settings)
	// a.configuration.Log(settings.WARN, fmt.Sprintf("2:%v", bb))
	// a.bot = afk.NewArenaBot(bb, gw)
	// backend.sendMessage("init", fmt.Sprintf("%v", a.dw))
	// a := bot.ScanText()
	// _ = a
	// b := activities.BoardsQuests(a.Result())
	// log.Warn(c.Red(b))
	// a.dw.outFn("APP STARTED")
	rt.EventsOn(ctx, "config", func(optionalData ...interface{}) {
		rt.LogWarning(a.ctx, fmt.Sprintf("[CONFIG UPD]: %v", optionalData))
	})

	rt.EventsOn(ctx, "task", func(optionalData ...interface{}) {
		rt.LogWarning(a.ctx, fmt.Sprintf("[TASK ARRIVed]: %v", optionalData))
		a.back.AdbConnect(fmt.Sprintf("%v", optionalData))
	})
}

func (a *App) sendMessage(event string, message string) {
	rt.EventsEmit(a.ctx, event, message)
}

// func (a *App) AdbConnect(str string) (isConnected bool) {
// 	dev, err := adb.Connect(str)

// 	if err == nil {
// 		a.bot.Connect(dev)
// 		isConnected = true

// 	}
// 	return

// }
