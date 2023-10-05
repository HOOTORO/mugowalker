package backend

import (
	"fmt"
	"mugowalker/backend/adb"
	"mugowalker/backend/afk"
	"mugowalker/backend/bot"
	"mugowalker/backend/localstore"
	"mugowalker/backend/pilot"
	"mugowalker/backend/settings"

	"github.com/wailsapp/wails"

	"gopkg.in/yaml.v3"
)

const filename = "conf.yaml"

type Config struct {
	Settings   *settings.Settings
	Pilot      *pilot.Pilot
	Bot        *afk.Daywalker
	Logger     *wails.CustomLogger
	localStore *localstore.LocalStore
}

// WailsInit performs setup when Wails is ready.
func (c *Config) WailsInit(out func(string, string)) error {
	c.Pilot = pilot.Default()
	gw := afk.New(c.Pilot)
	bb := bot.New(out, c.Settings)

	c.Bot = afk.NewArenaBot(bb, gw)
	// c.Logger.Info("Config initialized...")
	return nil
}

// NewConfig returns a new instance of Config.
func NewConfig() *Config {
	c := &Config{}
	c.localStore = localstore.NewLocalStore()

	a, err := c.localStore.Load(filename, true)
	if err != nil {
		c.Settings = settings.Default()
	}
	if err = yaml.Unmarshal(a, &c.Settings); err != nil {
		fmt.Printf("error")
	}
	return c
}

func (c *Config) AdbConnect(str string) (isConnected bool) {
	dev, err := adb.Connect(str)

	if err == nil {
		c.Bot.Connect(dev)
		isConnected = true

	}
	return

}

// func Run(ctx *context.Context) *Backend {

// 	return &Backend{
// 		runtimeCtx:    nil,
// 		configuration: configuration,
// 		user:          user,
// 		dw:            bot,
// 	}
// }
