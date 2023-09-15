package main

import (
	"context"
	"mugowalker/backend/adb"
	"mugowalker/backend/afk"
	"mugowalker/backend/bot"
	"mugowalker/backend/cfg"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	dw  *afk.Daywalker
}

// type Msg struct {
// 	Component string
// 	Message   string
// 	err       error
// }

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	fn := func(s1, s2 string) { a.sendMessage(s1, s2) }
	u := cfg.ActiveUser()
	u.GameAccount = "fdsf"
	u.DeviceSerial = "fdsdfdf"
	gw := afk.New(u)
	bb := bot.New(fn)
	a.dw = afk.NewArenaBot(bb, gw)
	a.sendMessage("init", cfg.F("%v", a.dw))
	// a := bot.ScanText()
	// _ = a
	// b := activities.BoardsQuests(a.Result())
	// log.Warn(c.Red(b))
	// a.dw.outFn("APP STARTED")
}

func (a *App) sendMessage(event, message string) {
	rt.EventsEmit(a.ctx, event, message)
}

func (a *App) AdbConnect(str string) (isConnected bool) {

	dev, err := adb.Connect(str)

	if err == nil {
		a.dw.Connect(dev)
		isConnected = true
		// a.dw.outFn("dev connnected")

	}
	return

}

// func (a *App) runMF()
