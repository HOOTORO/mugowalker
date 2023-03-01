package ui

import (
	"fmt"
	"strings"
	c "worker/cfg"
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
func listView(m appmenu) string {
	return menulistStl.Render(m.list.View())
}

// func nlistView(m menuModel) string {
// 	return m.cnct.View()
// }

// /////////////////////////////////
// /////// VIEW Input /////////////
// ///////////////////////////////
func inputFormView(m appmenu) string {
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
func (m *appmenu) statuStr() string {

	// con, emu = c.Red("Offline"), c.Red("Shutdown")

	return statusStl.Render(m.IsAdbAvailible())
}

func (m *appmenu) runningTasksPanel() string {
	log.Tracef("Upd status spanel....%v:%v", m.state.adbconn, m.state.vmPid)
	// var s, rt string
	s := m.statuStr()
	rt := fmt.Sprintf("\n"+
		m.state.spinme.View()+" Runing task %s...\n\n", taskName.Render(m.choice))

	for _, res := range m.state.taskmsgs {
		if res.Task == "" {
			rt += "...............................................\n"
		} else {
			rt += fmt.Sprintf("[%s] %s %s\n", res.Task, c.Cyan("|>"), res.Message)
		}
	}

	rt += indent.String(helpStyle.Render("\nPress 'alt+s' to hide/show this panel\n'Ctrl + <- ↑ ↓ ->' to change menu sizes"), 1)
	rt = runnunTaskStyle.Render(rt)
	return lipgloss.JoinVertical(lipgloss.Top, s, rt)
}

func (m *appmenu) runninVMs() string {
	if rf < refresh {
		return emuStatus.Render(vmstatus)
	}
	ems := c.Deserialize(emulator.AllVendors)
	var vms []string
	for _, v := range ems {
		ps, e := c.Tasklist(v)
		if e == nil && ps != nil {
			vms = append(vms, c.F("%v:%v", ps[0].Name, ps[0].Pid))
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

func (m *appmenu) IsAdbAvailible() string {
	r := c.F("User		|> %v <|\n", c.Cyan(m.conf.userSettings.Account))

	connectionstatus := c.Red("Disconnected")
	if m.state.adbconn > 0 {
		connectionstatus = c.Green("Connected")
		statusStl.BorderForeground(brightGreen)
	}
	r += c.F("Device |> %v | %v <|\n", m.conf.userSettings.Connection, connectionstatus)

	g := c.Red("Off")
	if m.state.gameStatus > 0 {
		g = c.Green("On")
	}
	r += c.F("Game |> %v | %v <|\n", m.conf.userSettings.AndroidGameID, g)
	return r
}
