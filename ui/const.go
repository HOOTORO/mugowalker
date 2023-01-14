package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const (
	header = "#### AFK Worker v0.1_alpha ####\n####### Active setup ###########"
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
	sep      = " >>> "
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
	toplevelmenu = []list.Item{
		item{title: "Device", desc: "Connection settings: ", children: devsmenu},
		item{title: "Tasks", children: tasks},
		item{title: "Settings", children: settingsV2},
	}

	devsmenu = []list.Item{
		item{title: "Availible devices", desc: "via 'adb devices'", children: devices},
		item{title: "Direct connect", desc: "from settings: ", children: func(m *menuModel) { m.devstatus = runConnect(m) }},
		item{title: "Run Bluestacks VM", desc: "Using args", children: func(m *menuModel) { m.devstatus = runBluestacks(m) }},
	}

	tasks = []list.Item{
		item{title: "Run all", children: func(m *menuModel) { runTask(m) }},
		item{title: "Do daily?", children: func(m *menuModel) { runTask(m) }},
		item{title: "Push Campain?", children: func(m *menuModel) { runTask(m) }},
		item{title: "Climb Towers?", children: towers},
	}
	towers = []list.Item{
		item{title: "Kings Tower", children: func(m *menuModel) { runTask(m) }},
		item{title: "Towers of Light", children: func(m *menuModel) { runTask(m) }},
		item{title: "Brutal Citadel", children: func(m *menuModel) { runTask(m) }},
		item{title: "World Tree", children: func(m *menuModel) { runTask(m) }},
		item{title: "Forsaken Necropolis", children: func(m *menuModel) { runTask(m) }},
	}
)

// Menu Styles
var (
	docStyle = lipgloss.NewStyle().Align(lipgloss.Left).Margin(0, 0, 0, 30)
	// titleBar
	hotStyle = lipgloss.NewStyle().Foreground(hotPink)
	tbStyle  = lipgloss.NewStyle().                                   // Height(1).  MarginTop(1).
			Border(lipgloss.ThickBorder()).BorderForeground(hotPink) //.
	// MarginLeft(20)
	// Header
	headerStyle = lipgloss.NewStyle().Width(50).
			Border(lipgloss.ThickBorder()).
			BorderBackground(purple).
			Align(lipgloss.Center).
			Bold(true) //.MarginBottom(0)

	statusStyle = lipgloss.NewStyle().
			MarginLeft(5)

	execRespStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
			BorderForeground(hotPink).Foreground(lipgloss.Color("#77DE77")).
			Align(lipgloss.Center).MarginLeft(50)

	redProps = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).
			Width(50).AlignHorizontal(lipgloss.Left)
	greenProps = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).
			Width(50).AlignHorizontal(lipgloss.Left)
)

var (
	// Output Style selectlist
	// MultiText Input
	topInputStyle       = lipgloss.NewStyle().Margin(5, 15)
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
