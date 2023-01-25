package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	a "worker/adb"
	"worker/afk"
	"worker/cfg"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

var mgt = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()

func getDevices() []list.Item {
	var devs []list.Item
	d, e := a.Devices()
	if e != nil {
		return devs
	}
	for _, v := range d {
		descu := fmt.Sprintf("State: %s, T_Id: %v, WMsize: %v", v.DevState, v.TransportId, v.Resolution)
		devs = append(devs, item{title: v.Serial, desc: descu, children: func(m *menuModel) tea.Cmd {
			return func() tea.Msg {
				return adbConnect(m.usersettings[connection])
			}
		}})
	}
	return devs
}

func runTask(m *menuModel) bool {
	cf := DtoCfg(m.usersettings)
	m.menulist.Styles.HelpStyle = noStyle
	m.spinme.Style = noStyle
	fn := func(s, d string) {
		// m.taskch <- notify(f("%v", s), d)
		m.taskch <- notify(f("%s", s), d)
	}
	log.Warnf(yellow("\nCHOSEN RUNTASK >>> %v <<<"), m.choice)

	afk.Push(cfg.PushCampain, cf, fn)
	return true
}

func runBluestacks(m *menuModel) bool {
	p := DtoCfg(m.usersettings)
	// pid, e := cfg.StartProc(bluestacksexe, strings.Fields(m.opts[bluestacks])...)
	cmd := cfg.RunProc(bluexe, p.Bluestacks.Args()...)

	m.bluestcksPid = cmd.Process.Pid
	m.statuStr()
	log.Warnf("\nwait in another gourutine %v", cmd.Process.Pid)

	go func() {
		e := cmd.Wait()
		if e != nil {
			m.taskch <- notify(bluexe, f("|> error: %v, pid: %v", e, m.bluestcksPid))
		}
		m.taskch <- notify(bluexe, f("|> finished, pid: %v", m.bluestcksPid))
		m.bluestcksPid = 0
	}()
	return checkBlueStacks(m)
}

func userSettings(c *cfg.Profile) map[Option]string {
	dto := make(map[Option]string, 0)
	dto[AppId] = c.Bluestacks.Package
	dto[VmName] = c.Bluestacks.Instance
	return dto
}

func CfgDto(conf *cfg.Profile) map[string]string {
	dto := make(map[string]string, 0)
	dto[connection] = conf.DeviceSerial
	dto[account] = conf.User.Account
	dto[game] = conf.User.Game
	dto[imagick] = strings.Join(conf.Imagick, " ")
	dto[tesseract] = strings.Join(conf.Tesseract, " ")
	dto[blueInstance] = conf.Bluestacks.Instance
	dto[bluePackage] = conf.Bluestacks.Package
	// dto[] = conf.
	// dto[] = conf.
	return dto
}

func DtoCfg(m map[string]string) *cfg.Profile {
	res := cfg.ActiveUser()

	for k, v := range m {
		switch k {
		case connection:
			res.DeviceSerial = v
		case account:
			res.User.Account = v
		case game:
			res.User.Game = v
		case imagick:
			res.Imagick = strings.Split(v, " ")
		case tesseract:
			res.Tesseract = strings.Split(v, " ")
		case blueInstance:
			res.Bluestacks.Instance = v
		case bluePackage:
			res.Bluestacks.Package = v
		}
	}
	return res
}

/////////////////////////////
//////// helper func ///////
///////////////////////////

func notify(ev, desc string) taskinfo {
	return taskinfo{Task: ev, Message: desc, Duration: time.Now()}
}

func updateDto(v map[string]string) {
	o := DtoCfg(v)
	cfg.Save(cfg.UserFile(o.User.Account+".yaml"), o)
}

func checkBlueStacks(m *menuModel) bool {
	if m.bluestcksPid != 0 {

		return cfg.ProcessInfo(m.bluestcksPid)
	}
	return false
}
func kill(pid int) bool {
	p, e := os.FindProcess(pid)
	if e == nil {
		e = p.Kill()
	}
	return e == nil
}

// func loglevel(m *menuModel) {

// }
