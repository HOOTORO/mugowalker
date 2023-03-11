package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	mainmenu = []list.Item{
		item("Device"),
		item("AFK Arena"),
		item("Parse Images"),
		item("Settings"),
		item("Statistics"),
	}

	device = []list.Item{
		item("manual connect via ADB"),
		item("last used"),
	}
	afkarena = []list.Item{
		item("Do daily"),
		item("Push Campain"),
		item("Push King Tower"),
	}
	imagework = []list.Item{
		item("Folder path"),
		item("Destination"),
	}
	settings = []list.Item{
		item("Account"),
		item("App (AFK Global/Test)"),
		item("Imagick"),
		item("Magick"),
		item("LogLevel"),
	}

	statistics = []list.Item{
		item("Summons"),
		item("Rankings"),
		item("Boss/Team DMG"),
	}
	loglvls = []list.Item{
		item("Trace"),
		item("Info"),
		item("Debug"),
		item("Error"),
		item("Fatal"),
		item("Panic"),
	}
	afkgameids = []list.Item{
		item("com.lilithgames.hgame.gp"),
		item("com.lilithgames.hgame.gp.id"),
	}
)

type menuItem struct {
	name   string
	parent *menuItem
	value  tea.Msg
}

var items = [][]list.Item{mainmenu, device, afkarena, imagework, settings, statistics, loglvls, afkgameids}

func FindItem(s string) []list.Item {
	for _, v := range items {
		for _, i := range v {
			if i.FilterValue() == s {
				return v
			}
		}

	}
	return nil
}

var (
	mm = &menuItem{
		name:   "mainmenu",
		parent: nil,
		value: func(m string) tea.Cmd {
			return func() tea.Msg {
				return mainmenu
			}
		},
	}
	devs = &menuItem{
		name:   "Device",
		parent: mm,
		value: func(m string) tea.Cmd {
			return func() tea.Msg {
				return device
			}
		},
	}
	aa = &menuItem{
		name:   "AFK Arena",
		parent: mm,
		value: func(m string) tea.Cmd {
			return func() tea.Msg {
				return afkarena
			}
		},
	}
	im = &menuItem{
		name:   "Parse Images",
		parent: mm,
		value: func(m string) tea.Cmd {
			return func() tea.Msg {
				return imagework
			}
		},
	}
	st = &menuItem{
		name:   "Settings",
		parent: mm,
		value: func(m string) tea.Cmd {
			return func() tea.Msg {
				return settings
			}
		},
	}
	sts = &menuItem{
		name:   "Statistics",
		parent: mm,
		value: func(m string) tea.Cmd {
			return func() tea.Msg {
				return statistics
			}
		},
	}
	lgs = &menuItem{
		name:   "LogLevel",
		parent: st,
		value: func(m string) tea.Cmd {
			return func() tea.Msg {
				return loglvls
			}
		},
	}
	afid = &menuItem{
		name:   "App (AFK Global/Test)",
		parent: nil,
		value: func(m string) tea.Cmd {
			return func() tea.Msg {
				return afkgameids
			}
		},
	}
)
