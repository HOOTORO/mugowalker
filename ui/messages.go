package ui

import (
	"strings"
	"time"
	"worker/adb"
	"worker/cfg"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type (
	vmStatusMsg   int
	connectionMsg int
)

func checkVM() tea.Msg {
	taskstr, e := cfg.Tasklist(bluestacksexe)
	if e != nil {
		return errMsg{e}
	}
	r := strings.Fields(taskstr)
	return vmStatusMsg(cfg.ToInt(r[1]))
}

func adbConnect(serial string) tea.Msg {
	dev, e := adb.Connect(serial)
	if e != nil {
		log.Errorf("\nConn err: %v", e)
		return errMsg{e}
	}

	return connectionMsg(dev.DevState)
}

// taskinfo is send when a pretend process completes.
type taskinfo struct {
	Task     string
	Message  string
	Duration time.Time
}

// A command that waits for the activity on a channel.
func activityListener(strch chan taskinfo) tea.Cmd {
	return func() tea.Msg {
		return taskinfo(<-strch)
	}
}
