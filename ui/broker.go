package ui

import (
	"fmt"
	"os"

	"worker/afk"
	"worker/bot"
	"worker/cfg"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

// var options = []string{"Loglevel", "Application Id", "VM Name", "Game", "Account", "Connect"}

// type Option uint

// const (
// 	LogLvl Option = iota + 1
// 	AppId
// 	VmName
// 	GameName
// 	AccountName
// 	ConnectStr
// 	// TessParams
// )

// func (o Option) String() string {
// 	return options[o-1]
// }

var runner *bot.BasicBot

var mgt = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()

func avalibleConnections(m *menuModel) []list.Item {
	var devs []list.Item
	fn := func(s, d string) {
		m.userstate.taskch <- notify(f("%s", s), d)
	}
	runner = bot.New(fn)
	d := runner.DiscoverDevices()
	for _, v := range d {
		descu := fmt.Sprintf("State: %s, T_Id: %v, WMsize: %v", v.DevState, v.TransportId, v.Resolution)
		devs = append(devs, item{title: v.Serial, desc: descu, children: func(m *menuModel) tea.Cmd {
			return func() tea.Msg {
				runner.Connect(v)
				return connectionMsg(runner.DevState)
			}
		}})
	}
	return devs
}

func runTask(m *menuModel) bool {
	cf := DtoCfg(m.conf.userSettings)
	m.menulist.Styles.HelpStyle = noStyle
	m.userstate.spinme.Style = noStyle

	log.Warnf(yellow("\nCHOSEN RUNTASK >>> %v <<<"), m.choice)
	ns := afk.Nightstalker(runner, cf)
	ns.Run(m.choice)
	return true
}

func runBluestacks(m *menuModel) bool {
	p := DtoCfg(m.conf.userSettings)
	// pid, e := cfg.StartProc(bluestacksexe, strings.Fields(m.opts[bluestacks])...)
	cmd := cfg.RunProc(bluexe, p.Bluestacks.Args()...)

	m.userstate.bluestcksPid = cmd.Process.Pid
	m.statuStr()
	log.Warnf("\nwait in another gourutine %v", cmd.Process.Pid)

	go func() {
		e := cmd.Wait()
		if e != nil {
			m.userstate.taskch <- notify(bluexe, f("|> error: %v, pid: %v", e, m.userstate.bluestcksPid))
		}
		m.userstate.taskch <- notify(bluexe, f("|> finished, pid: %v", m.userstate.bluestcksPid))
		m.userstate.bluestcksPid = 0
	}()
	return checkBlueStacks(m)
}

func userSettings(c *cfg.Profile) map[Option]string {
	dto := make(map[Option]string, 0)
	dto[AppId] = c.Bluestacks.Package
	dto[VmName] = c.Bluestacks.Instance
	dto[LogLvl] = c.Loglevel
	dto[GameName] = c.User.Game
	dto[AccountName] = c.User.Account
	dto[ConnectStr] = c.DeviceSerial
	return dto
}

func ocrSettings(c *cfg.Profile, e cfg.Executable) map[string]string {
	dto := make(map[string]string, 0)
	var currentkey string
	for i, v := range c.ImagickCfg() {
		if i%2 == 0 {
			currentkey = v
		} else {
			dto[currentkey] = v
		}
	}
	return dto
}

func CfgDto(conf *cfg.Profile) map[string]string {
	dto := make(map[string]string, 0)
	dto[connection] = conf.DeviceSerial
	dto[account] = conf.User.Account
	dto[game] = conf.User.Game
	dto[blueInstance] = conf.Bluestacks.Instance
	dto[bluePackage] = conf.Bluestacks.Package
	return dto
}

func DtoCfg(m map[Option]string) *cfg.Profile {
	res := cfg.ActiveUser()

	for k, v := range m {
		switch k {
		case ConnectStr:
			res.DeviceSerial = v
		case AccountName:
			res.User.Account = v
		case GameName:
			res.User.Game = v
		case VmName:
			res.Bluestacks.Instance = v
		case AppId:
			res.Bluestacks.Package = v
		}
	}
	return res
}

/////////////////////////////
//////// helper func ///////
///////////////////////////

func updateDto(v map[Option]string) {
	o := DtoCfg(v)
	cfg.Save(cfg.UserFile(o.User.Account+".yaml"), o)
}

func checkBlueStacks(m *menuModel) bool {
	if m.userstate.bluestcksPid != 0 {

		return cfg.ProcessInfo(m.userstate.bluestcksPid)
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
