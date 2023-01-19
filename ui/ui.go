package ui

import (
	"fmt"
	"os"
	"worker/cfg"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"

	"github.com/sirupsen/logrus"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

var red, green, cyan, yellow, mag func(...interface{}) string

func init() {
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	mag = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()
}

var log *logrus.Logger

func init() {
	log = cfg.Logger()
}

func RunMainMenu(options map[string]string) error {
	log.Debug("entered UI")
	m := InitialMenuModel(options)
	// m.header = headerStyle.Render(header)
	m.menulist.Title = header
	m.menulist.SetSize(40, 30)
	m.menulist.SetShowHelp(true)
	m.menulist.SetShowPagination(true)
	m.menulist.SetShowTitle(true)
	m.menulist.SetShowStatusBar(false)
	m.menulist.Styles.Title = tStyle
	m.menulist.Styles.TitleBar = tbStyle

	log.Debugf("Run p, w/ param %s", m)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
		return err
	}

	return nil
}

func NotifyUI(task, desc string) {

}

func InitialMenuModel(userOptions map[string]string) menuModel {
	m := menuModel{
		menulist:     list.New(availMenuItems(), list.NewDefaultDelegate(), 19, 0),
		parents:      nil,
		choice:       "",
		focusIndex:   0,
		manyInputs:   make([]textinput.Model, 0),
		cursorMode:   textinput.CursorBlink,
		quitting:     false,
		usersettings: userOptions,
		taskch:       make(chan taskinfo),
		taskmsgs:     make([]taskinfo, showLastTasks),
		spinme:       spinner.New(),
		showmore:     true,
	}

	m.spinme.Spinner = spinner.Moon
	m.spinme.Style = spinnerStyle
	return m
}
