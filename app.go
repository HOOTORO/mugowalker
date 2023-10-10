package main

import (
	"context"
	"fmt"
	"mugowalker/backend"
	s "mugowalker/backend/settings"
	"mugowalker/backend/taskmanager"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	b   *backend.Config
}

// NewApp creates a new App application struct
func NewApp(b *backend.Config) *App {
	return &App{b: b}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	tm := taskmanager.NewTaskManager(a.b, a.FireEvent)
	tm.InitDevice(a.b.DevicePath)

	rt.EventsOn(ctx, s.CFG_SET, func(optionalData ...interface{}) {
		data := fmt.Sprintf("[CONFIG UPD]: %v", optionalData)
		a.FireEvent(s.MSG, data)
		tm.UpdateConfig(optionalData...)
	})
	rt.EventsOn(ctx, s.CFG_PILOT, func(optionalData ...interface{}) {
		tm.UpdatePilot(optionalData...)
	})
	rt.EventsOn(ctx, s.TASK, func(optionalData ...interface{}) {
		rt.LogWarning(a.ctx, fmt.Sprintf("[TASK]: %v", optionalData))
		tm.RunTask(fmt.Sprintf("%v", optionalData...))
	})
	rt.EventsOn(ctx, s.ADB_CONN, func(optionalData ...interface{}) {
		rt.LogWarning(a.ctx, fmt.Sprintf("[CONNECTION]: %v", optionalData))
		result := tm.InitDevice(fmt.Sprintf("%v", optionalData...))
		if result {
			a.FireEvent(s.ADB_STS, "success")
		} else {
			a.FireEvent(s.ADB_STS, "failed")
		}
	})
}

func (a *App) OnDomReady(ctx context.Context) {
	a.FireEvent(s.MSG, "DOM Loaded")
}

func (a *App) FireEvent(event string, message string) {
	rt.EventsEmit(a.ctx, event, message)
}
