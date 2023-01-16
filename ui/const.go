package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	header = "#### AFK Worker v0.1_alpha ####\n####### Active setup ###########"
)

var f = fmt.Sprintf

// keymapping
const (
	connection    = "connection"
	account       = "account"
	game          = "game"
	taskconfigs   = "taskconfigs"
	imagick       = "imagick"
	tesseract     = "tesseract"
	bluestacks    = "bluestacks"
	adbp          = "adb"
	magick        = "magick"
	tesserexe     = "tess"
	bluestacksexe = "HD-Player"
)

const (
	hotPink     = lipgloss.Color("#FF06B7")
	darkGray    = lipgloss.Color("#767676")
	purple      = lipgloss.Color("99")
	brightGreen = lipgloss.Color("#00FF00")
	bloodRed    = lipgloss.Color("#FF0000")
	someG       = lipgloss.Color("#00FFa0")
	someR       = lipgloss.Color("#FFa000")
	sep         = " >>> "
)

var mds = [...]string{"select", "inpuit", "multin", "exec"}

type Mode int

const (
	selectList Mode = iota + 1
	inputMessage
	multiInput
	runExec
)

func (m Mode) String() string {
	return mds[m-1]
}

type Status int

func availMenuItems() []list.Item {
	toplevelmenu = append(toplevelmenu, availTowers()...)
	log.Debugf("Menu items: %v", toplevelmenu)
	return toplevelmenu
}

var (
	settings = func(m menuModel) []list.Item {
		var items []list.Item
		for k, v := range m.opts {
			items = append(items, item{title: k, children: initTextModel(v, false, "")})
		}
		return items
	}
	settingsV2 = func(m menuModel) []textinput.Model {
		var items []textinput.Model
		for k, v := range m.opts {
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
			title: "Check",
			desc:  "prcss",
			children: func(m *menuModel) tea.Cmd {
				return func() tea.Msg {
					return test(m)
				}
			},
		},
		item{
			title: "Connect to",
			desc:  "serial/ip set in 'Device'",
			children: func(m *menuModel) bool {
				return runConnect(m)
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
			children: settingsV2,
		},
		item{
			title: "Do daily?",
			desc:  "Do quest till 100pts",
			children: func(m *menuModel) bool {
				return runTask(m)
			},
		},
		item{
			title: "Push Campain?",
			desc:  "if you cant",
			children: func(m *menuModel) bool {
				return runTask(m)
			},
		},
		item{
			title: "Kings Tower",
			desc:  "Not yours",
			children: func(m *menuModel) bool {
				return runTask(m)
			},
		},

		// item{title: "Climb Towers?", children: towers},
		// item{title: "Tasks", children: tasks},
	}

	tasks = []list.Item{
		item{title: "Run all",
			children: func(m *menuModel) bool {
				return runTask(m)
			},
		},
	}
	towers = []list.Item{
		item{
			title: "Towers of Light",
			children: func(m *menuModel) bool {
				return runTask(m)
			},
		},
		item{
			title: "Brutal Citadel",
			children: func(m *menuModel) bool {
				return runTask(m)
			},
		},
		item{
			title: "World Tree",
			children: func(m *menuModel) bool {
				return runTask(m)
			},
		},
		item{
			title: "Forsaken Necropolis",
			children: func(m *menuModel) bool {
				return runTask(m)
			},
		},
	}
)

// Menu Styles
var (
	docStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Margin(0, 0, 0, 6)
	// titleBar
	hotStyle = lipgloss.NewStyle().Foreground(hotPink)

	tbStyle = lipgloss.NewStyle().                                   // Height(1).  MarginTop(1).
		Border(lipgloss.ThickBorder()).BorderForeground(hotPink) //.
	// MarginLeft(20)
	// Header
	headerStyle = lipgloss.NewStyle().
			Width(40).
			Border(lipgloss.ThickBorder()).
			BorderBackground(purple).
			Align(lipgloss.Center).
			Bold(true) //.MarginBottom(0)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("69"))

	execRespStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(hotPink).
			Foreground(lipgloss.Color("#77DE77")).
			Align(lipgloss.Bottom).
			MarginLeft(50)

	statusStyle = lipgloss.NewStyle().
			MarginLeft(1).
			Border(lipgloss.RoundedBorder()).
			Bold(true).
			Width(40).
			Align(lipgloss.Left)
	//	s := termenv.
	//	taskStyle = lipgloss.NewStyle().Foreground()
	runnunTaskStyle = statusStyle.Copy().MarginTop(2).UnsetBorderStyle()
	//		Align(lipgloss.Center)
	//		MarginLeft(10)
	//	happyClr = colorful.FastHappyColor()
	helpStyle = docStyle.Copy().
			MarginBottom(3).
		// UnsetMarginLeft().
		MarginLeft(1).
		Align(lipgloss.Center)

	quitStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000")).
			MarginBackground(lipgloss.Color("#00FF00")).
			Margin(60)
)

var (
	redProps   = lipgloss.NewStyle().Foreground(bloodRed).Render
	greenProps = lipgloss.NewStyle().Foreground(brightGreen).Render

	// Output Style selectlist
	// MultiText Input
	topInputStyle       = lipgloss.NewStyle().Margin(10, 15)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(5)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(5).Foreground(lipgloss.Color("170"))
)
