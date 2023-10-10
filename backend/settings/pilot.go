package settings

import (
	"fmt"
)

// Pilot profile
type Pilot struct {
	DevicePath string `json:"dev"`
	Account    string `json:"account"`
	GameId     string `json:"game"`
	Online     bool   `json:"online"`
}

func (up *Pilot) String() string {
	return fmt.Sprintf("\n\t[A] |> %v [D] |> %v [O] |> %v", up.Account, up.DevicePath, up.Online)
}

// New user profile
func New(accname, game string) *Pilot {
	return &Pilot{Account: accname, GameId: game, Online: false}
}

func DefaultPilot() *Pilot {
	return &Pilot{DevicePath: "localhost:5555", Account: "Bad Dev", GameId: "com.lilithgames.hgame.gp.id", Online: false}
}
