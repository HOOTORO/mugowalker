package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/list"
)

/////////////////////////////
//// UPD. SelectList ///////
///////////////////////////

func updateList(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
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
				case textinput.Model:
					m.textInput = chld
					m.mode = inputMessage
					// m.textInput, cmd = m.textInput.Update(msg)
					return m, cmd

				case func(m menuModel) []textinput.Model:
					m.manyInputs = chld(m)
					m.mode = multiInput
					// m.textInput, cmd = m.textInput.Update(msg)
					return m, cmd
				case []list.Item, func(m menuModel) []list.Item:
					// go deeper in menu
					m.mode = selectList
					m.parents = append(m.parents, m.menulist)
					m.menulist.SetItems(itm.NextLevel(m))
				case func(m *menuModel):
					chld(&m)
					m.mode = runExec
					// m.updateStatus()
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
		// 	h, v := docStyle.GetFrameSize()
		// 	m.menulist.SetSize(msg.Width-h, msg.Height-v)
	}

	m.menulist, cmd = m.menulist.Update(msg)
	return m, cmd
}

// ///////////////////////////
// ///// UPD. Input /////////
// /////////////////////////
func updateInput(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			m.mode = selectList
			m.opts[m.choice] = m.textInput.Value()
			updateDto(m.opts)
			m.updateStatus()
			return m, cmd
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// ///////// UPD Multiline input /////////
func updateFormInput(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
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
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if (s == "enter" && m.focusIndex == len(m.manyInputs)) || s == "esc" {
				m.mode = selectList
				updateDto(m.opts)
				m.updateStatus()
				// cmd = m.updatemanyInputs(msg)
				return m, cmd
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
	r := strings.NewReplacer(sep, "")
	for i := range m.manyInputs {
		if m.manyInputs[i].Value() != "" {
			m.opts[r.Replace(m.manyInputs[i].Prompt)] = m.manyInputs[i].Value()
		}
		m.manyInputs[i], cmds[i] = m.manyInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func updateExec(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter", "backspace", "esc":
			m.mode = selectList
			m.response = "some response"
			updateDto(m.opts)
			m.updateStatus()
			// return m, cmd
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}
	time.Sleep(3 * time.Second)
	m.mode = selectList
	m.response = "some response"
	m.updateStatus()
	m.menulist, cmd = m.menulist.Update(msg)
	return m, cmd
}
