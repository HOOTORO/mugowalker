package ui

import (
	"fmt"
	"os"
	"worker/cfg"

	"github.com/sirupsen/logrus"

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

func RunMainMenu(c *cfg.Profile) error {
	log.Debug("entered UI")
	// options := CfgDto(c)
	optionsv2 := userSettings(c)
	img := ocrSettings(c, cfg.MagicExe)
	tess := ocrSettings(c, cfg.TessExe)
	m := InitialMenuModel(tess, img, optionsv2)
	// m.header = headerStyle.Render(header)
	m.menulist.Title = header
	// m.menulist.SetSize(40, 10)
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
