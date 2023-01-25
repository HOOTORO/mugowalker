package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
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

// usersettings v2
var options = []string{"Application Id", "VM Name"}

type Option uint

const (
	AppId Option = iota + 1
	VmName
)

func (o Option) String() string {
	return options[o-1]
}

func availMenuItems() []list.Item {
	toplevelmenu = append(toplevelmenu, availTowers()...)
	log.Debugf("Menu items: %v", toplevelmenu)
	return toplevelmenu
}

var (
	toplevelmenu = []list.Item{
		item{
			title:    "My Device",
			desc:     "Setup platform where to run autotasks",
			children: deviceSetup,
		},
		item{
			title:    "Loglevel",
			desc:     "Change app log level",
			children: loglevel,
		},

		// item{
		// 	title: "Connect to",
		// 	desc:  "serial/ip set in 'Device'",
		// 	children: func(m *menuModel) tea.Cmd {
		// 		return func() tea.Msg {
		// 			return adbConnect(m.usersettings[connection])
		// 		}
		// 	},
		// },
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
	deviceSetup = func(m menuModel) []list.Item {
		var items []list.Item
		items = append(items, item{title: "ADB Connect", desc: "Connect via TCP/IP to emulator or remote device",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return adbConnect(m.usersettings[connection])
				}
			}})

		items = append(items, item{
			title:    "Emulator",
			desc:     "Setup bluestacks settings",
			children: emulatorSettings(m),
		})
		items = append(items, getDevices()...)

		return items
	}

	emulatorSettings = func(m menuModel) []list.Item {
		var items []list.Item
		items = append(items, item{
			title: "Launch Bluestacks",
			desc:  "With a given args",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runBluestacks(m)
				}
			},
		})
		items = append(items, item{
			title:    "Change arguments",
			desc:     "VM name, Launch Application",
			children: blueArgs,
		})
		return items
	}
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

	blueArgs = func(m menuModel) []textinput.Model {
		var input []textinput.Model
		input = append(input, initTextModel(m.usersettingsv2[VmName], true, VmName.String()))
		input = append(input, initTextModel(m.usersettingsv2[AppId], false, AppId.String()))
		return input
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
	loglevel = func(m menuModel) []list.Item {
		var items []list.Item
		current := log.GetLevel()
		all := logrus.AllLevels

		for _, lvl := range all {
			if lvl != current {
				items = append(items, item{title: lvl.String(), children: func(m *menuModel) tea.Cmd {
					return func() tea.Msg {
						log.SetLevel(lvl)
						NotifyUI("LogLvl", "Changed to >"+lvl.String())
						return log.GetLevel()
					}

				}})
			}
		}
		return items
	}
)
