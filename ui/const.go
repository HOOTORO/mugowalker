package ui

import 	"github.com/charmbracelet/bubbles/list"
type Mode int

const (
    Select Mode = iota+1
    Strinput

)

var (

    truemainmenu = []list.Item{
        item{title: "Device", desc: "Device/emulator to run bot"},
        item{title: "Tasks", desc: "Push, Dailies and many more"},
        item{title: "Settings", desc: "OCR, Game Locations, Debug etc..."},
        }

	mainmenu = []string{
		"Run all",
		"Do daily?",
		"Push Campain?",
		"Climb Towers?",
		"OCR Settings",
	}
	tower = []string{
        "Kings Tower",
        "Towers of Light",
        "Brutal Citadel",
        "World Tree",
        "Forsaken Necropolis",
    }
)
