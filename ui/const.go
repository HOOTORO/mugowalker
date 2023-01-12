package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type Level int

const (
	Top Level = iota + 1
	First
	Second
)

var (
	toplevelmenu = []list.Item{
		item{title: "Device", children: devsmenu},
		item{title: "Game", children: gamemenu},
		item{title: "Tasks", children: tasks},
		item{title: "Settings", children: settings},
	}

	devsmenu = []list.Item{
		item{title: "Connect to discovered"},
		item{title: "Set IP directly"},
	}

	gamemenu = []list.Item{
		item{title: "Set Game", children: strInputModel("Name")},
		item{title: "Set Account", children: strInputModel("Username")},
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
	docStyle = lipgloss.NewStyle().Margin(3, 3).Height(6).MaxHeight(40) //.AlignHorizontal(lipgloss.Center)
	// MultiText Input
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	brHelpStyle         = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

var (
	titleStyle = lipgloss.NewStyle().Padding(0).PaddingLeft(3).Width(40).Height(10).MarginBottom(5).Foreground(lipgloss.Color("99"))
	//	titleBarStyle     = lipgloss.NewStyle().Padding(0).MarginBottom(3).Width(40).Height(10).Background(lipgloss.Color("#AEAEAE"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(5)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(5).Foreground(lipgloss.Color("170"))
	statusStyle       = list.DefaultStyles().StatusBar.PaddingLeft(10)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).Height(20) // PaddingBottom(30).MarginBottom(30).AlignVertical(lipgloss.Top)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

var (
	hotStyle    = lipgloss.NewStyle().Foreground(hotPink)
	commonStyle = lipgloss.NewStyle().Foreground(darkGray)
	dS          = lipgloss.NewStyle().Margin(4, 10).Height(4)
	headerStyle = lipgloss.NewStyle().MaxWidth(200).Width(73).AlignVertical(lipgloss.Center) //.Border(lipgloss.ThickBorder()).BorderBackground(lipgloss.Color("99"))
)
