package ui

import (
	"os"
	c "worker/cfg"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

var l *logrus.Logger

func init() {
	l = c.Logger()
}

func RunUI(user c.AppUser) (err error) {
	sessionUser = &AppUser{}
	if user != nil {
		sessionUser.AccountName = user.Account()
		sessionUser.Connection = user.DevicePath()
		sessionUser.Loglvl = user.Loglevel()
	}

	inputModel := initialModel("Account")

	m := initSelectModel(mainmenu)

	core := initCore(inputModel, m)
	if _, err := tea.NewProgram(core, tea.WithAltScreen()).Run(); err != nil {
		l.Println("Error running program:", err)
		os.Exit(1)
	}

	return err
}
