package ui

import (
	"strings"

	a "worker/adb"
	"worker/cfg"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type (
	vmStatusMsg   int
	connectionMsg int
	loglevelMsg   int
)

func checkVM() tea.Msg {
	taskstr, e := cfg.Tasklist(bluexe)
	if e != nil {
		return errMsg{e}
	}
	r := strings.Fields(taskstr)
	return vmStatusMsg(cfg.ToInt(r[1]))
}

func adbConnect(serial string) tea.Msg {
	dev, e := a.Connect(serial)
	if e != nil {
		log.Errorf("\nConn err: %v", e)
		return errMsg{e}
	}

	return connectionMsg(dev.DevState)
}
func setLoglevel(lvl string) tea.Msg {
	obj, e := logrus.ParseLevel(lvl)
	if e != nil {
		log.Error("Wrong LogLevel String")
		return errMsg{e}
	}
	log.SetLevel(obj)
	return loglevelMsg(obj)
}

// taskinfo is send when a pretend process completes.
type taskinfo struct {
	Task    string
	Message string
	// Duration time.Time
}

// A command that waits for the activity on a channel.
func activityListener(strch chan taskinfo) tea.Cmd {
	return func() tea.Msg {
		return taskinfo(<-strch)
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
// Validator functions to ensure valid input
