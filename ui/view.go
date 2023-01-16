package ui

import (
	"fmt"
	"strings"
)

// /////////////////////////////////
// /////// VIEW Input /////////////
// ///////////////////////////////

// /////////////////////////
// func (m menuModel) View() string {
// //////////////////////////
func listView(m menuModel) string {
	s := strings.Join([]string{m.header, m.menulist.View()}, "\n")
	return docStyle.Render(s)
}

func inputView(m menuModel) string {
	return fmt.Sprintf(
		"Please, enter <%v> property\n\n%s\n\n%s",
		m.choice,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

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
	m.showmore = true
	s := fmt.Sprintf("%v\n\n\nRunnin task: %v", m.header, m.activeTask)
	return docStyle.Render(s)
}
