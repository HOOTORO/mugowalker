package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	listHeight   = 14
	defaultWidth = 23
)
const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	purple   = lipgloss.Color("#3300ff")
	red      = lipgloss.Color("#ff0000")
	green    = lipgloss.Color("#00FF00")
)

// ///////////////////////
// Input Model Styles///
// ////////////////////
var (
	underLine = lipgloss.NewStyle().Underline(true).Bold(true).Faint(true).UnderlineSpaces(false)
	rndBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(purple).
			Padding(1).
			PaddingLeft(3)
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

// ///////////////////////
// Select List Styles///
// /////////////////////
var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

////////////////////////
////// status box /////
//////////////////////

var (
	statusStyle = lipgloss.NewStyle().
		Align(lipgloss.Top).
		Padding(0, 5, 0, 5).
		// Width(10).Height(5).
		// Margin(0, 0, 10, 5).
		Border(lipgloss.ThickBorder(), true, true, true, true).
		// BorderBackground(red).
		BorderForeground(green).
		AlignHorizontal(lipgloss.Left).
		AlignVertical(lipgloss.Top)
)
