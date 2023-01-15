package ui

import (
	"fmt"
	"os"

	"worker/cfg"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = cfg.Logger()
}

func RunMainMenu(options map[string]string) error {
	log.Debug("entered UI")
	m := InitialMenuModel(options)
	m.header = headerStyle.Render(header)
	m.menulist.Title = "Choose..."
	m.menulist.SetSize(50, 20)
	m.menulist.SetShowHelp(true)
	m.menulist.SetShowPagination(true)
	m.menulist.SetShowTitle(true)
	m.menulist.SetShowStatusBar(true)
	m.menulist.Styles.Title = hotStyle
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
