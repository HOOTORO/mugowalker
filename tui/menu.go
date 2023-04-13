package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	mainmenu = []list.Item{
		item{"Device", selectView},
		item{"AFK Arena", selectView},
		item{"Parse Images", selectView},
		item{"Settings", selectView},
		item{"Statistics", selectView},
	}

	device = []list.Item{
		item{"manual connect via ADB", selectView},
		item{"last used", selectView},
	}
	afkarena = []list.Item{
		item{"Do daily", selectView},
		item{"Push Campain", selectView},
		item{"Push King Tower", selectView},
	}
	imagework = []list.Item{
		item{"Folder path", selectView},
		item{"Destination", selectView},
	}
	settings = []list.Item{
		item{"Account", selectView},
		item{"App (AFK Global/Test)", selectView},
		item{"Imagick", selectView},
		item{"Magick", selectView},
		item{"LogLevel", selectView},
	}

	statistics = []list.Item{
		item{"Summons", selectView},
		item{"Rankings", selectView},
		item{"Boss/Team DMG", selectView},
	}
	loglvls = []list.Item{
		item{"Trace", selectView},
		item{"Info", selectView},
		item{"Debug", selectView},
		item{"Error", selectView},
		item{"Fatal", selectView},
		item{"Panic", selectView},
	}
	afkgameids = []list.Item{
		item{"com.lilithgames.hgame.gp", selectView},
		item{"com.lilithgames.hgame.gp.id", selectView},
	}
)

type menuItem struct {
	name   string
	parent *menuItem
	value  tea.Msg
}

var items = [][]list.Item{mainmenu, device, afkarena, imagework, settings, statistics, loglvls, afkgameids}

func Finditem(s string) []list.Item {
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
