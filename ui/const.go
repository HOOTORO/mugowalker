package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const (
	header = "AFK Worker v0.1_alpha\n####### Active setup ###########"
)

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
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	purple   = lipgloss.Color("99")
	sep      = "> "
)

type Mode int

const (
	selectList Mode = iota + 1
	inputMessage
	runExec
)

type Status int

const (
	Device Status = iota + 1
	Settings
	Software
)

var (
	args   = [...]string{imagick, tesseract, bluestacks}
	device = [...]string{connection, game, account}
	soft   = [...]string{adbp, magick, tesserexe, bluestacksexe}
)

func (s Status) Opts() []string {
	switch s {
	case Device:
		return device[:]
	case Settings:
		return args[:]
	case Software:
		return soft[:]
	default:
		return nil
	}
}

var (
	settings = func(m menuModel) []list.Item {
		var items []list.Item
		for k, v := range m.opts {
			items = append(items, item{title: k, children: initTextModel(v)})
		}
		return items
	}
	settingsV2 = func(m menuModel) []textinput.Model {
		var items []textinput.Model
		for _, v := range m.opts {
			items = append(items, initTextModel(v))
		}
		return items
	}
	devices = func(m menuModel) []list.Item {
		return getDevices()
	}
	toplevelmenu = []list.Item{
		item{title: "Device", desc: "Connection settings: ", children: devsmenu},
		item{title: "Tasks", children: tasks},
		item{title: "Settings", children: settings},
	}

	devsmenu = []list.Item{
		item{title: "Availible devices", desc: "via 'adb devices'", children: devices},
		item{title: "Direct connect", desc: "from settings: ", children: func(m menuModel) {
			Connect(m.opts[connection])
		}},
		item{title: "Run Bluestacks VM", desc: "Using args", children: func(m menuModel) { runBluestacks(m) }},
	}

	tasks = []list.Item{
		item{title: "Run all", children: "Do everything by a little"},
		item{title: "Do daily?", children: "Only dailies till 100 pt"},
		item{title: "Push Campain?", children: "Strike through CAMPAIN"},
		item{title: "Climb Towers?", children: towers},
	}
	towers = []list.Item{
		item{title: "Kings Tower"},
		item{title: "Towers of Light"},
		item{title: "Brutal Citadel"},
		item{title: "World Tree"},
		item{title: "Forsaken Necropolis"},
	}
)

// Menu Styles
var (
	docStyle = lipgloss.NewStyle().Margin(1, 5)
	// titleBar
	hotStyle = lipgloss.NewStyle().Foreground(hotPink)
	tbStyle  = lipgloss.NewStyle().Height(1).MarginTop(1).
			Border(lipgloss.ThickBorder()).BorderForeground(hotPink)
	// Header
	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderBackground(purple).
			AlignVertical(lipgloss.Center).
			Bold(true).MarginBottom(2)

	redProps = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).
			Width(50).AlignHorizontal(lipgloss.Left)
	greenProps = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).
			Width(50).AlignHorizontal(lipgloss.Left)
)

var (
	// Output Style selectlist
	// MultiText Input
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
	helpStyle         = docStyle.Copy().MarginBottom(3)
)
