package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	sep      = "> "
)

type Mode int

const (
	SelectList Mode = iota + 1
	InputMessage
	RunExec
)

var (
	toplevelmenu = []list.Item{
		item{title: "Device", children: devsmenu},
		item{title: "Game", children: gamemenu},
		item{title: "Tasks", children: tasks},
		item{title: "Settings", children: settings},
	}

	devsmenu = []list.Item{
		item{title: "Connect to discovered", children: getDevices()},
		item{title: "Set IP directly", children: initTextModel("host:port")},
		item{title: "Run Bluestacks VM", children: initTextModel("host:port")},
	}

	gamemenu = []list.Item{
		item{title: "Set Game", children: initTextModel("Name")},
		item{title: "Set Account", children: initTextModel("Username")},
	}

	tasks = []list.Item{
		item{title: "Run all", children: "Do everything by a little"},
		item{title: "Do daily?", children: "Only dailies till 100 pt"},
		item{title: "Push Campain?", children: "Strike through CAMPAIN"},
		item{title: "Climb Towers?", children: towers},
	}
	settings = []list.Item{
		item{title: "Imagick args", children: "see https://tesseract-ocr.github.io/tessdoc/ImproveQuality.html#image-processing \n https://imagemagick.org/Usage/transform/#vision"},
		item{title: "Tesseract args", children: "run 'tesseract --help-extra' or '--print-parameters'"},
		item{title: "Bluestacks", children: ""},
	}
	towers = []list.Item{
		item{title: "Kings Tower"},
		item{title: "Towers of Light"},
		item{title: "Brutal Citadel"},
		item{title: "World Tree"},
		item{title: "Forsaken Necropolis"},
	}
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
	// titleStyle = lipgloss.NewStyle().Padding(0).PaddingLeft(3).Width(40).Height(10).MarginBottom(5).Foreground(lipgloss.Color("99"))
	//	titleBarStyle     = lipgloss.NewStyle().Padding(0).MarginBottom(3).Width(40).Height(10).Background(lipgloss.Color("#AEAEAE"))
	///// Ui.go
	itemStyle         = lipgloss.NewStyle().PaddingLeft(5)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(5).Foreground(lipgloss.Color("170"))
	helpStyle         = docStyle.Copy().MarginBottom(3)
)

var (
	docStyle    = lipgloss.NewStyle().Margin(1, 5)
	headerStyle = lipgloss.NewStyle().MaxWidth(200).Width(73).
			AlignVertical(lipgloss.Center)

	hotStyle    = lipgloss.NewStyle().Foreground(hotPink)
	tophedStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderBackground(lipgloss.Color("99")).
			Bold(true)
	// commonStyle = lipgloss.NewStyle().Foreground(darkGray)
	// dS          = lipgloss.NewStyle().Margin(4, 10).Height(4)
)
