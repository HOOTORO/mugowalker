package ui

import (
	"fmt"
	"strconv"
	"strings"
	c "worker/cfg"

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

type inputDialog struct {
	parent  *appmenu
	inputs  []userField
	focused int
	err     error
}

type userField struct {
	ID           uint
	fieldname    string
	defaultValue string
	input        textinput.Model
}

func (m inputDialog) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debugf(c.Mgt("(INPUT) incoming. -> %+v [%T]"), msg, msg)
	var cmds = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m.parent.Update(adbConnect(c.F("%v:%v", m.inputs[0].input.Value(), m.inputs[1].input.Value())))
			}
			m.nextInput()
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].input.Blur()
		}
		m.inputs[m.focused].input.Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i].input, cmds[i] = m.inputs[i].input.Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m inputDialog) View() string {
	dialog := " Welcome to the ricefields...\n"
	for i, v := range m.inputs {
		dialog += inputStyle.Width(v.input.Width).Render(v.fieldname)
		dialog += "\n"
		dialog += m.inputs[i].input.View()

	}
	dialog += continueStyle.Render("\nContinue ->\n")
	return topInputStyle.Render(dialog)
}

// nextInput focuses the next input field
func (m *inputDialog) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *inputDialog) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func initMultiModel(p *appmenu, fields []userField) inputDialog {

	return inputDialog{
		parent:  p,
		inputs:  fields,
		focused: 0,
		err:     nil,
	}
}

//	func newInput(field userField) textinput.Model {
//		model := textinput.New()
//		if field.focus {
//			model.Focus()
//		}
//		model.Placeholder = field.placeholder
//		model.CharLimit = field.charlim
//		model.Width = field.width
//		model.Prompt = field.promt
//		model.SetValue(field.defaultValue)
//		log.Debugf("field: %v\n model: %v", field, model)
//		return model
//	}
func NewUserField(name, defVal string) userField {
	model := textinput.New()
	// if focus {
	// 	model.Focus()
	// }

	// model.Placeholder = field.placeholder
	// model.CharLimit = field.charlim
	// model.Width = field.width
	// model.Prompt = field.promt
	// model.SetValue(field.defaultValue)
	// log.Debugf("field: %v\n model: %v", field, model)
	return userField{fieldname: name, defaultValue: defVal, input: model}
}

func (uf userField) WithFocus() userField {
	uf.input.Focus()
	return uf
}

func (uf userField) WithCharlimit(n int) userField {
	uf.input.CharLimit = n
	return uf
}
func (uf userField) WithWidth(n int) userField {
	uf.input.Width = n
	return uf
}
func (uf userField) WithPrompt(s string) userField {
	uf.input.Prompt = s
	return uf
}
func (uf userField) WithPlaceholder(s string) userField {
	uf.input.Placeholder = s
	return uf
}
