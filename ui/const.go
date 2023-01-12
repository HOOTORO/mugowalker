package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	Select Mode = iota + 1
	Strinput
)

var (
	truemainmenu = []list.Item{
		item{title: "Device", desc: "Device/emulator to run bot"},
		item{title: "Tasks", desc: "Push, Dailies and many more"},
		item{title: "Settings", desc: "OCR, Game Locations, Debug etc..."},
	}

	mainmenu = []string{
		"Run all",
		"Do daily?",
		"Push Campain?",
		"Climb Towers?",
		"OCR Settings",
	}
	tower = []string{
		"Kings Tower",
		"Towers of Light",
		"Brutal Citadel",
		"World Tree",
		"Forsaken Necropolis",
	}
)

var (
	docStyle            = lipgloss.NewStyle().Margin(3, 3).Height(6).MaxHeight(40).AlignHorizontal(lipgloss.Center)
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
	titleStyle        = lipgloss.NewStyle().Padding(0).PaddingLeft(3).Width(40).Height(10).MarginBottom(5).Foreground(lipgloss.Color("99"))
//	titleBarStyle     = lipgloss.NewStyle().Padding(0).MarginBottom(3).Width(40).Height(10).Background(lipgloss.Color("#AEAEAE"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(5)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(5).Foreground(lipgloss.Color("170"))
	statusStyle       = list.DefaultStyles().StatusBar.PaddingLeft(10)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(30).MarginBottom(30).AlignVertical(lipgloss.Top)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

var (
	hotStyle      = lipgloss.NewStyle().Foreground(hotPink)
	commonStyle = lipgloss.NewStyle().Foreground(darkGray)
	dS = lipgloss.NewStyle().Margin(4,10).Height(4)
	headerStyle = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderBackground(lipgloss.Color("99")).MaxWidth(200).Width(70).AlignVertical(lipgloss.Center)
)