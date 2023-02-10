package ui

import (
	"errors"
	a "worker/adb"
	"worker/bot"
	c "worker/cfg"
	"worker/emulator"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type (
	vmStatusMsg    int
	connectionMsg  int
	loglevelMsg    int
	prevousMenuMsg int
	appOnlineMsg   int
)

func checkVM() tea.Msg {
	taskstr, e := c.Tasklist(bluexe)
	noxprc, e1 := c.Tasklist(emulator.Nox.String())
	if e != nil && e1 != nil {
		return errMsg{e}
	}
	if len(taskstr) > 0 {

		return vmStatusMsg(taskstr[0].Pid)
	} else {
		return vmStatusMsg(noxprc[0].Pid)
	}
}

func prevousMenu(m menuModel) tea.Cmd {
	prevousState := len(m.parents) - 1
	return func() tea.Msg {
		return prevousMenuMsg(prevousState)
	}
}

func runAfk(m *menuModel) tea.Msg {
	if m.state.adbconn > 0 {
		if runner.IsAppRunnin(m.conf.userSettings.AndroidGameID) > 0 {
			return appOnlineMsg(1)
		} else {
			runner.StartApp(m.conf.userSettings.AndroidGameID)
			return appOnlineMsg(1)
		}
	} else {
		return errMsg{err: errors.New("Device offline")}
	}
}

func initAfk(m *menuModel) tea.Cmd {
	return func() tea.Msg {
		log.Debugf("INIT AFK MAFAKA : %+v", c.Cyan(m))
		if m.state.adbconn > 0 {
			if runner.IsAppRunnin(m.conf.userSettings.AndroidGameID) > 0 {
				return appOnlineMsg(1)
			} else {
				runner.StartApp(m.conf.userSettings.AndroidGameID)
				return appOnlineMsg(1)
			}
		}
		return appOnlineMsg(0)
	}
}

func adbConnect(serial string) tea.Msg {
	dev, e := a.Connect(serial)
	if e != nil {
		log.Errorf("\nConn err: %v", e)
		return errMsg{e}
	}

	return connectionMsg(dev.DevState)
}

func connectDevice(m menuModel) tea.Cmd {
	runner = bot.New(m.UserOutput)
	d := runner.DiscoverDevices()
	dev := &a.Device{DevState: 0}
	for _, v := range d {
		// descu := c.F("State: %s, T_Id: %v, WMsize: %v", v.DevState, v.TransportId, v.Resolution)
		if v.Serial == m.conf.userSettings.Connection {
			dev = v
		}
	}
	return func() tea.Msg {
		if dev.DevState != 0 {
			m.conf.userSettings.Connection = dev.Serial
			runner.Connect(dev)
			m.state.adbconn = 1
		} else {
			runner.Device = dev
		}
		log.Debugf("Device --> %+v", c.Cyan(runner))
		return connectionMsg(runner.DevState)
	}
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
