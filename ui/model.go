package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

type menuModel struct {
	mode   Mode
	header string

	menulist list.Model
	parents  []list.Model
	choice   string

	textInput textinput.Model

	focusIndex int
	manyInputs []textinput.Model
	cursorMode textinput.CursorMode

	quitting bool
	err      error
	cursor   int
}

func InitialMenuModel() menuModel {
	m := menuModel{
		mode:       SelectList,
		header:     "Worker Setup",
		menulist:   list.New(toplevelmenu, list.NewDefaultDelegate(), 19, 0),
		parents:    nil,
		choice:     "",
		textInput:  initTextModel("..."),
		focusIndex: 0,
		manyInputs: make([]textinput.Model, 0),
		cursorMode: textinput.CursorBlink,
		quitting:   false,
		err:        nil,
		cursor:     0,
	}

	// 	manyInputs: make([]textinput.Model, 3),
	// }

	// var t textinput.Model
	// for i := range m.manyInputs {
	// 	t = textinput.New()
	// 	t.CursorStyle = cursorStyle
	// 	t.CharLimit = 32

	// 	switch i {
	// 	case 0:
	// 		t.Placeholder = "Game"
	// 		t.Focus()
	// 		t.PromptStyle = focusedStyle
	// 		t.TextStyle = focusedStyle
	// 	case 1:
	// 		t.Placeholder = "Account"
	// 		t.CharLimit = 64
	// 	case 2:
	// 		t.Placeholder = "Connection str"
	// 		//                    t.EchoMode = textinput.EchoPassword
	// 		t.EchoCharacter = '•'
	// 	}

	// 	m.manyInputs[i] = t

	return m
}

type (
	errMsg error
)

type item struct {
	title    string //, desc
	children interface{}
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.String() }
func (i item) FilterValue() string { return i.title }

func (i item) String() string {
	elems := ""
	switch children := i.children.(type) {
	case []list.Item:
		for _, v := range children {
			elems += "<" + v.FilterValue() + sep
		}
		return elems
	case textinput.Model:
		return children.Placeholder
	}
	return fmt.Sprintf("%s", i.children)
}

func (i item) NextLevel() []list.Item {
	chld, ok := i.children.([]list.Item)
	if ok {
		return chld
	}
	return nil
}

// //////////////////////////
// /////// General //////////
// init / update / view ///
// ////////////////////////
func (m menuModel) Init() tea.Cmd {
	return textinput.Blink
}

// ////////////////////////
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// always exit keys
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	switch m.mode {
	case SelectList:
		return updateList(msg, m)
	case InputMessage:
		return updateInput(msg, m)
	}
	return m, tea.Quit
}

///////////////////////////////////

func (m menuModel) View() string {
	var s string
	if m.quitting {
		return "\n  See you later, Space Cowboy!\n\n"
	}
	switch m.mode {
	case SelectList:
		s = listView(m)
	case InputMessage:
		s = inputView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

//////////////////////////////////

/////////////////////////////
//// UPD. SelectList ///////
///////////////////////////

func updateList(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "spacebar":
			// Send the choice on the channel and exit.
			m.choice = m.menulist.SelectedItem().FilterValue()

			if itm, ok := m.menulist.SelectedItem().(item); ok {
				switch chld := itm.children.(type) {
				case textinput.Model:
					m.textInput = chld
					m.mode = InputMessage
					m.textInput, cmd = m.textInput.Update(msg)
					return m, cmd
				case []list.Item:
					// go deeper in menu
					m.parents = append(m.parents, m.menulist)
					m.menulist.SetItems(itm.NextLevel())
				}
			}

		case "backspace":
			// go up to top using chain parents
			if len(m.parents) > 0 {
				m.menulist.SetItems(m.parents[len(m.parents)-1].Items())
				m.parents = m.parents[:len(m.parents)-1]
			}
		}
		// May be... some day
		// case tea.WindowSizeMsg:
		// 	h, v := docStyle.GetFrameSize()
		// 	m.menulist.SetSize(msg.Width-h, msg.Height-v)
	}
	m.menulist, cmd = m.menulist.Update(msg)
	return m, cmd
}

//	else {
//		// itm, ok = m.list.SelectedItem()
//		return m, tea.Quit
//
// //////////////////////////

// ///////////////////////////
// ///// UPD. Input /////////
// /////////////////////////
func updateInput(msg tea.Msg, m menuModel) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			m.mode = SelectList
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
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
			if s == "enter" && m.focusIndex == len(m.manyInputs) {
				return m, tea.Quit
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
	cmd := m.updatemanyInputs(msg)

	return m, cmd
}

// func (m multiInputModel) Data() []textinput.Model {
// 	return m.manyInputs
// }

// *
func (m menuModel) updatemanyInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.manyInputs))

	// Only text manyInputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.manyInputs {
		m.manyInputs[i], cmds[i] = m.manyInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// /////////////////////////////////
// /////// VIEW Input /////////////
// ///////////////////////////////

// /////////////////////////
// func (m menuModel) View() string {
// //////////////////////////
func listView(m menuModel) string {
	return docStyle.Render(m.header + m.menulist.View())
}

func inputView(m menuModel) string {
	return fmt.Sprintf(
		"Please, enter %v\n\n%s\n\n%s",
		m.textInput.Placeholder,
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

	return b.String()
}

//func (m *MenuModel) Update(ud UserData) *MenuModel {
//	switch ud := ud.(type) {
//	case MenuModel:
//		if ud.choice != "" {
//			m.choice = ud.choice
//		}
//		if ud.header != "" {
//			m.header = ud.header
//		}
//		if ud.promt != "" {
//			m.promt = ud.promt
//		}
//		if len(ud.items) > 0 {
//			m.items = ud.items
//		}
//		if ud.footer != "" {
//			m.footer = ud.footer
//		}
//	}
//	return m
//}

///////////////////////
//// helper func  ////
/////////////////////

func initTextModel(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return ti
}
