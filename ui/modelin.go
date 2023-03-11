package ui

import (
	c "worker/cfg"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type modelIn struct {
	inputs  []userField
	focused int
	err     error
}

func (m modelIn) Ufield(n string) *userField {
	for _, v := range m.inputs {
		if n == v.name {
			return &v
		}
	}
	return nil
}

func initialModel(fields ...string) modelIn {

	var inputs []userField = make([]userField, len(fields))

	for i, v := range fields {
		inputs[i] = uField(v, 20)

	}
	if len(inputs) > 0 {
		inputs[0].in.Focus()
	}

	return modelIn{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m modelIn) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelIn) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l.Tracef(c.RFW("<| INPUT UPD. INC. |> %+v"), msg)
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			c.UpdateUserInfo(sessionUser)
			if m.focused == len(m.inputs)-1 {
				return m, inputDone(m.inputs[m.focused].in.Value())
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].in.Blur()
		}
		m.inputs[m.focused].in.Focus()

		// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i].in, cmds[i] = m.inputs[i].in.Update(msg)
		switch m.inputs[i].name {
		case AccField:
			sessionUser.AccountName = m.inputs[i].in.Value()
		case ConnectionField:
			sessionUser.Connection = m.inputs[i].in.Value()
		case LoglvlField:
			sessionUser.Loglvl = m.inputs[i].in.Value()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m modelIn) View() string {
	var res string

	res = underLine.Render("Required data\n")
	for _, v := range m.inputs {
		res = lipgloss.JoinVertical(lipgloss.Left, res, inputStyle.Width(v.in.Width).Render(v.name), v.in.View())
	}

	res = lipgloss.JoinVertical(lipgloss.Left, rndBorder.Render(res), continueStyle.Render("Enter to continue...\n"))
	return res
}

// nextInput focuses the next input field
func (m *modelIn) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *modelIn) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

type userField struct {
	name string
	in   textinput.Model
}

func uField(name string, l int) userField {
	uf := userField{name: name, in: textinput.New()}
	uf.in.CharLimit = l
	uf.in.Width = l
	uf.in.Prompt = ""
	return uf
}
