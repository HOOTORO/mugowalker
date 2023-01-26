package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/list"
)

/////////////////////////////
//// UPD. SelectList ///////
///////////////////////////

func updateMenu(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", " ":
			if itm, ok := m.menulist.SelectedItem().(item); ok {

				m.choice = itm.FilterValue()

				switch chld := itm.children.(type) {
				// go deeper in menu
				case []list.Item, func(m menuModel) []list.Item:
					m.parents = append(m.parents, m.menulist)
					m.menulist.SetItems(itm.NextLevel(m))

				case func(m *menuModel) []textinput.Model:
					m.inputChosen = true
					m.focusIndex = -1
					m.manyInputs = chld(&m)

					return updateInput(msg, m)

					//// Run something go to updateExec
				case func(m *menuModel) tea.Cmd:
					m.taskch <- notify(itm.title, "Launched!")
					cmd = chld(&m)
					m.menulist.Update(msg)
					return m, cmd
				}
			}

		case "backspace":
			// go up to top using chain parents
			if len(m.parents) > 0 {
				m.menulist.SetItems(m.parents[len(m.parents)-1].Items())
				m.parents = m.parents[:len(m.parents)-1]
			}
		}
		// log.Debugf(red("FOCUS # %v"), m.menulist.SelectedItem().FilterValue())
		// May be... some day
		// case tea.WindowSizeMsg:
		// 	m.menulist.SetSize(msg.Width/2, msg.Height)
		// m.statusInfo =
	}
	// 	// updateDto(m.opts)
	// 	m.updateStatus()

	m.menulist, cmd = m.menulist.Update(msg)
	return m, cmd
}

// //////////////////////////////
// ///// UPD. Input ////////////
// ////////////////////////////
// / UPD Multiline input /////
func updateInput(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.manyInputs))
			for i := range m.manyInputs {
				cmds[i] = m.manyInputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

			// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down", "esc":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			log.Warnf("Focus: %v  len(%v)", (&m).focusIndex, len(m.manyInputs))
			if (s == "enter" && m.focusIndex == len(m.manyInputs)) || s == "esc" {
				m.inputChosen = false
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.manyInputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.manyInputs)
			}

			cmds := make([]tea.Cmd, len(m.manyInputs))
			for i := 0; i <= len(m.manyInputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.manyInputs[i].Focus()
					m.manyInputs[i].PromptStyle = focusedStyle
					m.manyInputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.manyInputs[i].Blur()
				m.manyInputs[i].PromptStyle = noStyle
				m.manyInputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd = m.updatemanyInputs(msg)

	return m, cmd
}

func (m menuModel) updatemanyInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.manyInputs))

	// Only text manyInputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	// r := strings.NewReplacer(sep, "")
	for i := range m.manyInputs {
		// if m.manyInputs[i].Value() != "" {
		// 	m.userSettings[r.Replace(m.manyInputs[i].Prompt)] = m.manyInputs[i].Value()
		// }
		m.manyInputs[i], cmds[i] = m.manyInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// func updateExec(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c":
// 			m.quitting = true
// 			return m, tea.Quit
// 		}
// 	}
// 	if item, ok := m.menulist.SelectedItem().(item); ok {
// 		switch fu := item.children.(type) {
// 		case func(m *menuModel) tea.Cmd:

// 			m.taskch <- notify(item.title, "Launched!")
// 			return m, fu(&m)
// 		}
// 	}
// 	// updateDto(m.opts)
// 	m.updateStatus()
// 	var cmd tea.Cmd
// 	m.menulist, cmd = m.menulist.Update(msg)
// 	return m, cmd
// }
