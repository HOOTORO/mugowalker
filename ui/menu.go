package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

// usersettings v2
var options = []string{"Loglevel", "Application Id", "VM Name", "Game", "Account", "Connect"}

type Option uint

const (
	LogLvl Option = iota + 1
	AppId
	VmName
	GameName
	AccountName
	ConnectStr
	// TessParams
)

func (o Option) String() string {
	return options[o-1]
}

func availMenuItems() []list.Item {
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
			title:    "My Settings",
			desc:     "Imagick, Tesseract and other",
			children: mySettings,
		},
		item{
			title:    "My Tasks",
			desc:     "Push, Daily etc...",
			children: myTasks,
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
			desc:  "LIGHTBEARERs home",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "Brutal Citadel",
			desc:  "Maulers trainning center",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "World Tree",
			desc:  "Wilders birthplace",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "Forsaken Necropolis",
			desc:  "Dead man's belongs here",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "Infernal Fortress",
			desc:  "Dead man's belongs here",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return runTask(m)
				}
			},
		},
		item{
			title: "Celestial Sanctum",
			desc:  "Dead man's belongs here",
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
		items = append(items, item{
			title:    "Availible devices",
			desc:     "'adb devices -l'",
			children: devices,
		})
		items = append(items, item{
			title:    "Emulator",
			desc:     "Setup bluestacks settings",
			children: emulatorSettings,
		})
		// items = append(items, getDevices()...)
		return items
	}
	mySettings = func(m menuModel) (out []list.Item) {
		out = append(out, item{
			title:    "Log Level",
			desc:     f("Current lvl |> %v", cyan(log.GetLevel().String())),
			children: loglevel,
		})
		out = append(out, item{
			title:    "Tesseract",
			desc:     "Parameters for OCR Engine",
			children: tessArgs,
		})
		out = append(out, item{
			title:    "Imagick",
			desc:     "Optimizing image before OCR",
			children: imagickArgs,
		})
		out = append(out, item{})
		return
	}

	myTasks = func(m menuModel) (out []list.Item) {
		out = append(out, item{
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
			})
		out = append(out, availTowers()...)
		return
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
	devices = func(m menuModel) []list.Item {
		var items []list.Item
		items = append(items, avalibleConnections(&m)...)
		items = append(items, item{
			title: "ADB Connect",
			desc:  "Connect via TCP/IP to emulator or remote device",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return initialMIModel(m)
				}
				// return inputModels(m.cursorMode, "HOST/IP", m.conf.userSettings[ConnectStr], "PORT", "5555")
			},
		})
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
						return setLoglevel(m.choice)
					}

				}})
			} else {
				items = append(items, item{
					title: lvl.String(),
					desc:  " ↑ Current level ",
				})
			}
		}
		return items
	}
)

// Settings inputs
var (
	imagickArgs = func(m *menuModel) (out []textinput.Model) {
		var pairs []string
		for k, v := range m.conf.magic {
			pairs = append(pairs, k, v)
		}
		out = inputModels(m.cursorMode, pairs...)
		return
	}
	tessArgs = func(m *menuModel) (out []textinput.Model) {
		var pairs []string
		for k, v := range m.conf.ocr {
			pairs = append(pairs, k, v)
		}
		out = inputModels(m.cursorMode, pairs...)
		return
	}
	settings = func(m *menuModel) []textinput.Model {
		var pairs []string
		for k, v := range m.conf.userSettings {
			pairs = append(pairs, k.String(), v)
		}
		return inputModels(m.cursorMode, pairs...)
	}

	blueArgs = func(m *menuModel) []textinput.Model {
		return inputModels(m.cursorMode, VmName.String(), m.conf.userSettings[VmName], AppId.String(), m.conf.userSettings[AppId])
	}
)

//	input field[1]name, field[1]placeholder... fiend[n]name, field[n+1]placeholder
//
// threat odds as  fieldnames set in prompt, evens as placeholder text
//
//	fields should be even sized,
func inputModels(cursorMode textinput.CursorMode, fields ...string) []textinput.Model {
	log.Warnf("inmodels: %v", fields)
	inputs := make([]textinput.Model, 0)
	for i, v := range fields {
		if i%2 == 0 {
			inputs = append(inputs, initTextModel(cursorMode, "", i == 0, v))
		} else {
			inputs[len(inputs)-1].Placeholder = v
		}
	}
	return inputs
}
