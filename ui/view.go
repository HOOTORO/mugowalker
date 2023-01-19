package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/muesli/reflow/indent"
)

// ////////////////////
// /// Menu view /////
// //////////////////
func listView(m menuModel) string {
	return menulistStyle.Render(m.menulist.View())
}

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
	_, err := fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	if err != nil {
		return ""
	}

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return topInputStyle.Render(b.String())
}

func execView(m menuModel) string {

	s := fmt.Sprintf("%v\n\n%v exec: %v", m.header, m.spinme.View(), m.choice)
	return menulistStyle.Render(s)
}

// /////////////////////
// // ui elements /////
// ///////////////////
func (m *menuModel) statuStr() string {
	var con, emu string

	con, emu = red("Offline"), red("Shutdown")

	if m.connectionStatus != 0 {
		con = green("Online")
	}
	if m.bluestcksPid != 0 {

		emu = green("Running")
	}
	if m.connectionStatus != 0 && m.bluestcksPid != 0 {
		statusStyle.BorderForeground(brightGreen)
	}

	t := f("Device --> [%v] \n"+
		"Profile --> \n\tGame |> %v\n\tUser |> %v"+
		"\nConnection status: %v \nBluestacks: %v",
		m.usersettings[connection], cyan(m.usersettings[game]),
		cyan(m.usersettings[account]), con, emu)
	return statusStyle.Render(t)
}

func (m *menuModel) runningTasksPanel() string {
	log.Debugf("Upd status spanel....%v:%v", m.connectionStatus, m.bluestcksPid)
	var s, rt string
	s = m.statuStr()
	rt = fmt.Sprintf("\n"+
		m.spinme.View()+" Runing task %s...\n\n", m.choice)

	for _, res := range m.taskmsgs {
		if res.Task == "" {
			rt += "...................................\n"
		} else {
			rt += fmt.Sprintf("[%s]|>%s<|\n", res.Task, res.Message)
		}
	}

	rt += helpStyle.Render("\nPress 'alt+s' to hide/show this panel\n")
	//		x, y := helpStyle.GetFrameSize()
	//		rt += hotStyle.Render(fmt.Sprintf("btw:\nw: %v, h:%v", x, y))
	rt = runnunTaskStyle.Render(rt)
	return indent.String(lipgloss.JoinVertical(lipgloss.Top, s, rt), 5)
}
