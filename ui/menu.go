package ui

import (
	"time"

	c "worker/cfg"
	"worker/emulator"

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

var mainmenu = func() []list.Item {
	var items []list.Item
	items = append(items, item{title: "Device", child: devices})
	items = append(items, item{title: "Settings", child: mySettings()})
	items = append(items, item{title: "AFK tasks", child: devices})
	return items
}

var connect = func(m appmenu) tea.Cmd {
	return func() tea.Msg {
		fields := make([]userField, 0)
		fields = append(fields,
			NewUserField("HOST", "127.0.0.1").WithPrompt(c.MgCy("|> ")).WithFocus(),
			NewUserField("PORT", "5555").WithPrompt(c.MgCy("|> ")),
		)
		return initMultiModel(&m, fields)
	}
}

func (o Option) String() string {
	return options[o-1]
}

func (o Option) Values() []string {
	return options
}

var (
	////////////////////////
	/////afk tasks /////////
	///////////////////////
	myTasks = func(m appmenu) (out []list.Item) {
		// desc:       "Do quest till 100pts",
		out = append(out, item{title: "Do daily?", child: botask},
			// desc:       "if you cant",
			item{title: "Push Campain?", child: botask},
			// desc:       "Not yours",
			item{title: "Kings Tower", child: botask})
		out = append(out, availTowers()...)
		return
	}
	towers = []list.Item{
		// desc:       "LIGHTBEARERs home",
		item{title: "Towers of Light", child: botask},
		// desc:       "Maulers trainning center",
		item{title: "Brutal Citadel", child: botask},
		// desc:       "Wilders birthplace",
		item{title: "World Tree", child: botask},
		// desc:       "Dead man's belongs here",
		item{title: "Forsaken Necropolis", child: botask},
		// desc:       "Dead man's belongs here",
		item{title: "Infernal Fortress", child: botask},
		// desc:       "Dead man's belongs here",
		item{title: "Celestial Sanctum", child: botask},
	}
)

// func
var (
	deviceSetup = func(s item) []list.Item {
		var items []list.Item
		// desc:       "'adb devices -l'",
		items = append(items, item{title: "Availible devices",
			child: func(a appmenu) tea.Cmd { return func() tea.Msg { return devices(a) } }})
		// desc:       "Setup bluestacks settings",
		items = append(items, item{title: "Emulator",
			child: func(a appmenu) tea.Cmd { return func() tea.Msg { return emulatorSettings(a) } }})
		// desc:       c.F("Run app: %v", m.conf.userSettings.AndroidGameID),
		items = append(items, item{title: "Start App", child: runApp})
		return items
	}
	mySettings = func() (out []list.Item) {
		// desc:       c.F("Current lvl |> %v", c.Cyan(log.GetLevel().String())),
		// action: func(a appmenu) tea.Cmd { return func() tea.Msg { return loglevel } },
		out = append(out, item{title: "Log Level", child: func(a appmenu) tea.Cmd { return func() tea.Msg { return setLoglevel } }})
		// desc:       "Parameters for OCR Engine",
		out = append(out, item{title: "Tesseract", child: func(a appmenu) tea.Cmd { return func() tea.Msg { return tessArgs(a) } }})
		// desc:   "Optimizing image before OCR",
		out = append(out, item{title: "Imagick", child: func(a appmenu) tea.Cmd { return func() tea.Msg { return imagickArgs(a) } }})
		return
	}

	emulatorSettings = func(m appmenu) []list.Item {
		var items []list.Item
		// desc:  "With a given args",
		items = append(items, item{title: "Launch Bluestacks", child: func(m appmenu) tea.Cmd { return func() tea.Msg { return runBluestacks(&m) } }})
		// desc:  "With a given args",
		items = append(items, item{title: "Launch Nox", child: func(m appmenu) tea.Cmd {
			return func() tea.Msg {
				emulator.Run(emulator.Nox, "com.lilithgames.hgame.gp.id")
				return runNox(&m)
			}
		},
		})
		// desc:   "VM name, Launch Application",
		items = append(items, item{title: "Change arguments", child: func(a appmenu) tea.Cmd { return func() tea.Msg { return blueArgs(a) } }})
		return items
	}
	runApp = func(m appmenu) tea.Cmd {
		return func() tea.Msg {
			return runAfk(&m)
		}
	}
	devices = func(m appmenu) []list.Item {
		var items []list.Item
		items = append(items, avalibleConnections(&m)...)

		return items
	}

	loglevel = func(m item) []list.Item {
		var items []list.Item
		current := log.GetLevel()
		all := logrus.AllLevels

		for _, lvl := range all {
			if lvl != current {
				items = append(items, item{title: lvl.String(), child: func(m appmenu) tea.Cmd {
					return func() tea.Msg {
						return setLoglevel(m.choice)
					}

				}})
			} else {
				items = append(items, item{
					title: lvl.String(),
					// desc:  " ↑ Current level ",
				})
			}
		}
		return items
	}
)

// Settings inputs
var (
	imagickArgs = func(m appmenu) (out []textinput.Model) {
		var pairs []string
		for k, v := range m.conf.magic {
			pairs = append(pairs, k, v)
		}
		out = inputModels(m.cursorMode, pairs...)
		return
	}
	tessArgs = func(m appmenu) (out []textinput.Model) {
		var pairs []string
		for k, v := range m.conf.ocr {
			pairs = append(pairs, k, v)
		}
		out = inputModels(m.cursorMode, pairs...)
		return
	}
	settings = func(m appmenu) []textinput.Model {
		return inputModels(m.cursorMode, AccountName.String(), m.conf.userSettings.Account)
	}

	blueArgs = func(m appmenu) []textinput.Model {
		return inputModels(m.cursorMode, VmName.String(), m.conf.userSettings.VMName, AppId.String(), m.conf.userSettings.AndroidGameID)
	}
)

// input field[1]name, field[1]placeholder... fiend[n]name, field[n+1]placeholder
//
// threat odds as  fieldnames set in prompt, evens as placeholder text
//
// fields should be even sized,
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

var (
	// helper func
	botask = func(m appmenu) tea.Cmd {
		return func() tea.Msg {
			return runBotTask(&m)
		}
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
