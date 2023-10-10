package main

import (
	"embed"
	"mugowalker/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {

	back := backend.NewConfig()
	back.WailsInit()

	// Create an instance of the app structure
	app := NewApp(back)
	// Create application with options
	err := wails.Run(&options.App{
		Title:         "Mugo Walker",
		Width:         900,
		Height:        1200,
		Frameless:     false,
		DisableResize: true,
		Fullscreen:    false,
		AlwaysOnTop:   false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:   &options.RGBA{R: 11, G: 11, B: 21, A: 230},
		Logger:             nil,
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.ERROR,
		OnStartup:          app.startup,
		OnDomReady:         app.OnDomReady,
		CSSDragProperty:    "--wails-draggable",
		CSSDragValue:       "drag",
		Bind: []interface{}{
			app,
			back,
		},
		ErrorFormatter: func(err error) any { return err.Error() },
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameAccessibilityHighContrastVibrantDark,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "Mugo Walker",
				Message: "© 2023 HOOTORO",
				Icon:    icon,
			},
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
	},
	)

	if err != nil {
		println("Error:", err.Error())
	}
}
