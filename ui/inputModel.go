package ui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/textinput"
)

func hostValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16 {
		return fmt.Errorf("host is too long")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, ".", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func portValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

const (
	host = iota
	port
)

type multiIputModel struct {
	parent  *menuModel
	inputs  []textinput.Model
	focused int
	err     error
}

type inputField struct {
	fieldname          string
	placeholder, promt string
	charlim, width     int
	focus              bool
}

func (m multiIputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m multiIputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debugf(mag("(mUPDA) incoming. -> %+v [%T]"), msg, msg)
	var cmds = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m.parent.Update(adbConnect(f("%v:%v", m.inputs[0].Value(), m.inputs[1].Value())))
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
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m multiIputModel) View() string {
	return topInputStyle.Render(fmt.Sprintf(
		` Please, enter
 %s
 %s
 %s
 %s
 %s
`,
		inputStyle.Width(30).Render("HOST"),
		m.inputs[host].View(),
		inputStyle.Width(6).Render("PORT"),
		m.inputs[port].View(),
		continueStyle.Render("Continue ->"),
	) + "\n")
}

// nextInput focuses the next input field
func (m *multiIputModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *multiIputModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func initMultiModel(p *menuModel, fields []inputField) multiIputModel {
	inpt := make([]textinput.Model, len(fields))
	for i, v := range fields {
		inpt[i] = newInput(v)

	}
	return multiIputModel{
		parent:  p,
		inputs:  inpt,
		focused: 0,
		err:     nil,
	}
}

func newInput(field inputField) textinput.Model {
	model := textinput.New()
	if field.focus {
		model.Focus()
	}
	model.Placeholder = field.placeholder
	model.CharLimit = field.charlim
	model.Width = field.width
	model.Prompt = field.promt
	log.Debugf("field: %v\n model: %v", field, model)
	return model
}
