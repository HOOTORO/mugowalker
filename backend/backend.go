package backend

import (
	"encoding/json"
	"fmt"
	"mugowalker/backend/localstore"
	"mugowalker/backend/settings"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

const (
	filename = "conf.json"
	account  = "acc.json"
)

type Config struct {
	*settings.Settings
	*settings.Pilot
	Log logger.Logger
	*localstore.LocalStore
}

// WailsInit performs setup when Wails is ready.
func (c *Config) WailsInit() error {
	c.Log = logger.NewFileLogger(c.Settings.Logfile)
	// os.Truncate(c.Settings.Logfile, 0)
	a, err := c.Load(account, true)
	if err != nil {
		c.Pilot = settings.DefaultPilot()
	}
	if err = json.Unmarshal(a, &c.Pilot); err != nil {
		fmt.Printf("error")
	}
	c.Log.Info("\n<------------------------------>\nConfig initialized...\n<------------------------------>\n\n\n\n\n")
	return nil
}

// NewConfig returns a new instance of Config.
func NewConfig() *Config {
	c := &Config{}
	c.LocalStore = localstore.NewLocalStore()

	a, err := c.Load(filename, true)
	if err != nil {
		c.Settings = settings.Default()
	}
	if err = json.Unmarshal(a, &c.Settings); err != nil {
		fmt.Printf("error")
	}
	return c
}
func (c *Config) CurrentConfiguration() *settings.Settings {
	return c.Settings
}
func (c *Config) CurrentPilot() *settings.Pilot {
	return c.Pilot
}
