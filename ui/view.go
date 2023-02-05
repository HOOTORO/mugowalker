package ui

import (
	"fmt"
	"strings"
	"worker/cfg"
	"worker/emulator"

	"github.com/charmbracelet/lipgloss"

	"github.com/muesli/reflow/indent"
)

const refresh = 10

var (
	rf       = 11
	vmstatus = ""
)

// ////////////////////
// /// Menu view /////
// //////////////////
func listView(m menuModel) string {
	return menulistStl.Render(m.menulist.View())
}

// func nlistView(m menuModel) string {
// 	return m.cnct.View()
// }

// /////////////////////////////////
// /////// VIEW Input /////////////
// ///////////////////////////////
func inputFormView(m menuModel) string {
	var b strings.Builder

	for i := range m.manyInputs {
		b.WriteString(m.manyInputs[i].View())
		if i < len(m.manyInputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.manyInputs) {
		button = &focusedButton
	}
	_, err := fmt.Fprintf(&b, "\n\n\t%s\n\n", *button)
	if err != nil {
		return ""
	}

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return topInputStyle.Render(b.String())
	// return indent.String(topInputStyle.Render(b.String()), 3)
}

// /////////////////////
// // ui elements /////
// ///////////////////
func (m *menuModel) statuStr() string {

	// con, emu = red("Offline"), red("Shutdown")

	return statusStl.Render(m.IsAdbAvailible())
}

func (m *menuModel) runningTasksPanel() string {
	log.Tracef("Upd status spanel....%v:%v", m.state.adbconn, m.state.vmPid)
	// var s, rt string
	s := m.statuStr()
	rt := fmt.Sprintf("\n"+
		m.state.spinme.View()+" Runing task %s...\n\n", taskName.Render(m.choice))

	for _, res := range m.state.taskmsgs {
		if res.Task == "" {
			rt += "...............................................\n"
		} else {
			rt += fmt.Sprintf("[%s] %s %s\n", res.Task, cyan("|>"), res.Message)
		}
	}

	rt += indent.String(helpStyle.Render("\nPress 'alt+s' to hide/show this panel\n'Ctrl + <- ↑ ↓ ->' to change menu sizes"), 1)
	rt = runnunTaskStyle.Render(rt)
	return lipgloss.JoinVertical(lipgloss.Top, s, rt)
}

func (m *menuModel) runninVMs() string {
	if rf < refresh {
		return emuStatus.Render(vmstatus)
	}
	ems := cfg.Deserialize(emulator.AllVendors)
	var vms []string
	for _, v := range ems {
		ps, e := cfg.Tasklist(v)
		if e == nil && ps != nil {
			vms = append(vms, f("%v:%v", ps[0].Name, ps[0].Pid))
		}
	}
	r := "VMs |> "
	if len(vms) == 0 {
		r = "No running VMs"
	} else {
		r += strings.Join(vms, " | ")
	}
	vmstatus = r
	rf = 0
	return emuStatus.Render(r)
}

func (m *menuModel) IsAdbAvailible() string {
	r := f("User		|> %v <|\n", cyan(m.conf.userSettings.Account))

	connectionstatus := red("Disconnected")
	if m.state.adbconn > 0 {
		connectionstatus = green("Connected")
		statusStl.BorderForeground(brightGreen)
	}
	r += f("Device |> %v | %v <|\n", m.conf.userSettings.Connection, connectionstatus)

	g := red("Off")
	if m.state.gameStatus > 0 {
		g = green("On")
	}
	r += f("Game |> %v | %v <|\n", m.conf.userSettings.AndroidGameID, g)
	return r
}
