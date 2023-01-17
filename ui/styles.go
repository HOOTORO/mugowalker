package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Margin(0, 0, 0, 6)

	// titleBar
	hotStyle = lipgloss.NewStyle().Foreground(hotPink)

	tbStyle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(hotPink)

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
			MarginLeft(30)

	statusStyle = lipgloss.NewStyle().
			MarginLeft(1).
			Border(lipgloss.RoundedBorder()).
			Bold(true).
			Width(50).
			Align(lipgloss.Left)

	runnunTaskStyle = statusStyle.Copy().MarginTop(2).UnsetBorderStyle().Width(70)

	//	happyClr = colorful.FastHappyColor()
	helpStyle = docStyle.Copy().
			MarginBottom(3).
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
	topInputStyle       = lipgloss.NewStyle().Margin(10, 0, 10, 10).Width(60)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
)
