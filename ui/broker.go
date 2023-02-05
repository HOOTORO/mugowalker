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

var runner *bot.BasicBot
var mgt = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()

func userSettings(c *cfg.Profile) *AppUser {
	return &AppUser{
		Connection:    c.DeviceSerial,
		Account:       c.User.Account,
		AndroidGameID: c.Bluestacks.Package,
		VMName:        c.Bluestacks.Instance,
		Loglvl:        c.Loglevel,
	}
}

func avalibleConnections(m *menuModel) []list.Item {
	var menuItems []list.Item

	runner = bot.New(m.UserOutput)
	d := runner.DiscoverDevices()

	for _, v := range d {
		descu := fmt.Sprintf("State: %s, T_Id: %v, WMsize: %v", v.DevState, v.TransportId, v.Resolution)
		menuItems = append(menuItems, item{title: v.Serial, desc: descu, children: func(m menuModel) tea.Cmd {
			return func() tea.Msg {
				m.conf.userSettings.Connection = v.Serial
				runner.Connect(v)
				return connectionMsg(runner.DevState)
			}
		}})
	}
	return menuItems
}

func runBotTask(m *menuModel) bool {
	m.menulist.Styles.HelpStyle = noStyle
	m.state.spinme.Style = noStyle
	s := m.state
	switch {
	case s.adbconn == 0:
		log.Errorf(ErrNoAdb.Error())
		m.UserOutput(m.choice, red(ErrNoAdb.Error()))
		return false
	case s.gameStatus == 0:
		log.Errorf(red("Game is not running"))
		m.UserOutput(m.conf.userSettings.AndroidGameID, red(ErrAppNotRunning.Error()))
		return false

	}

	log.Warnf(yellow("\n	CHOSEN RUNTASK >>> %v <<<"), m.choice)
	ns := afk.Nightstalker(runner, m.conf.userSettings)
	ns.Run(m.choice)
	return true
}

func runBluestacks(m *menuModel) bool {
	// pid, e := cfg.StartProc(bluestacksexe, strings.Fields(m.opts[bluestacks])...)
	cmd := cfg.RunProc(bluexe, cfg.ActiveUser().Bluestacks.Args()...)

	m.state.vmPid = cmd.Process.Pid
	m.statuStr()
	log.Warnf("\nwait in another gourutine %v", cmd.Process.Pid)

	go func() {
		e := cmd.Wait()
		if e != nil {
			m.state.taskch <- notify(bluexe, f("|> error: %v, pid: %v", e, m.state.vmPid))
		}
		m.state.taskch <- notify(bluexe, f("|> finished, pid: %v", m.state.vmPid))
		m.state.vmPid = 0
	}()
	return checkBlueStacks(m)
}

func ocrSettings(c *cfg.Profile, e cfg.Executable) map[string]string {
	dto := make(map[string]string, 0)
	var args []string
	if e == cfg.MagicExe {
		args = c.ImagickCfg()
	} else {
		args = c.TesseractCfg()
	}
	var currentkey string
	for i, v := range args {
		if i%2 == 0 {
			currentkey = v
		} else {
			dto[currentkey] = v
		}
	}
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
	if m.state.vmPid != 0 {

		return cfg.IsProcess(m.state.vmPid)
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
