package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var f = fmt.Sprintf

// keymapping
const (
	connection   = "connection"
	account      = "account"
	game         = "game"
	taskconfigs  = "taskconfigs"
	imagick      = "imagick"
	tesseract    = "tesseract"
	blueInstance = "bluestance"
	bluePackage  = "bluepackage"
	bluexe       = "HD-Player"
)

func availMenuItems() []list.Item {
	toplevelmenu = append(toplevelmenu, availTowers()...)
	log.Debugf("Menu items: %v", toplevelmenu)
	return toplevelmenu
}

var (
	toplevelmenu = []list.Item{
		item{
			title: "Launch Bluestacks",
			desc:  "check args in settings before!",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runBluestacks(m)
				}
			},
		},
		item{
			title: "Kill Blue",
			desc:  "Kills bluestacks by pid (if not nil)",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return kill(m.bluestcksPid)
				}
			},
		},
		item{
			title: "Connect to",
			desc:  "serial/ip set in 'Device'",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return adbConnect(m.usersettings[connection])
				}
			},
		},
		item{
			title:    "Availible devices",
			desc:     "'adb devices -l'",
			children: devices,
		},
		item{
			title:    "Settings",
			desc:     "Imagick, Tesseract and other",
			children: settings,
		},
		item{
			title: "Do daily?",
			desc:  "Do quest till 100pts",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "Push Campain?",
			desc:  "if you cant",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "Kings Tower",
			desc:  "Not yours",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
	}

	tasks = []list.Item{
		item{title: "Run all",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
	}
	towers = []list.Item{
		item{
			title: "Towers of Light",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "Brutal Citadel",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "World Tree",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "Forsaken Necropolis",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
	}
)

// func
var (
	settings = func(m menuModel) []textinput.Model {
		var items []textinput.Model
		for k, v := range m.usersettings {
			items = append(items, initTextModel(v, false, k))
		}
		if len(items) > 0 {
			items[0].Focus()
			items[0].PromptStyle = focusedStyle
			items[0].TextStyle = focusedStyle
		}
		return items
	}

	devices = func(m menuModel) []list.Item {
		var items []list.Item
		items = append(items, getDevices()...)
		return items
	}

	availTowers = func() []list.Item {
		var items []list.Item
		switch time.Now().UTC().Weekday() {
		case time.Monday:
			items = append(items, towers[0])
		case time.Tuesday:
			items = append(items, towers[1])
		case time.Wednesday:
			items = append(items, towers[2])
		case time.Thursday:
			items = append(items, towers[3])
		case time.Friday:
			items = append(items, towers[0])
			items = append(items, towers[1])
		case time.Saturday:
			items = append(items, towers[2])
			items = append(items, towers[3])
		case time.Sunday:
			items = append(items, towers...)
		}
		return items
	}
)
