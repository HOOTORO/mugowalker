package pilot

import (
	"fmt"
	"mugowalker/backend/cfg"
)

// Pilot profile
type Pilot struct {
	DevicePath  string   `yaml:"dev"`
	Account     string   `yaml:"account"`
	Game        string   `yaml:"game"`
	TaskConfigs []string `yaml:"taskconfigs"`
}

func (up *Pilot) String() string {
	return fmt.Sprintf("[Game] |> %v [Account] |> %v", up.Game, up.Account)
}

// New user profile
func New(accname, game string) *Pilot {
	return &Pilot{Account: accname, Game: game}
}

func Default() *Pilot {
	return &Pilot{Account: "Bad Dev", Game: "AFK"}
}

// UpdateUserInfo saves to yaml into Userhome dir
func (au *Pilot) Update() {

	cfg.Save(au.Account+".yaml", au)

}
