package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"worker/cfg"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc string
}

type (
	errMsg error
)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type fancymodel struct {
	header string
	list   list.Model
}

func (m fancymodel) Init() tea.Cmd {
	return nil
}

func (m fancymodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m fancymodel) View() string {
	return m.header+docStyle.Render(m.list.View())
}

type MenuModel struct {
	mode   Mode
	list   list.Model
	header string
	items  []string
	choice string
	promt  string
	footer string
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

type multiInputModel struct {
	header     string
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

type selectModel struct {
	choices []string
	cursor  int
	choice  string
}

type altSelectModel struct {
	list list.Model
}

type inputModel struct {
	textInput textinput.Model
	err       error
}

func initialModel() inputModel {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return inputModel{
		textInput: ti,
		err:       nil,
	}
}
func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

type UserData interface{}

func initialUserInfoModel() multiInputModel {
	m := multiInputModel{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Game"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Account"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Connection str"
			//                    t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}
func (m multiInputModel) Init() tea.Cmd {
	return textinput.Blink
}
func (m multiInputModel) Data() []textinput.Model {
	return m.inputs
}
func (m multiInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

// *
func (m multiInputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m multiInputModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
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
func (m selectModel) Init() tea.Cmd {
	return nil
}

func Init() *MenuModel {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = m.choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m *MenuModel) View() string {
	res := ListDesc(m.items, m.header, "Default[%v] --> ", "0")
	res += fmt.Sprintf("\n\n\n%v", m.footer)
	return res
}
func (m selectModel) View() string {
	s := strings.Builder{}
	s.WriteString("What kind of Bubble Tea would you like to order?\n\n")

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

func (m *MenuModel) Menu() string {
MainMenu:
	m.header = "AFK Bot\n What bot should do?"
	m.items = mainmenu
	choice := UserListInput(m.items, m.header, "Exit")

	switch choice {

	case 4:
	Towers:
		choice = UserListInput(tower, "Which one?", "Back")
		switch {
		case choice > 0:
			return yellow("Climbing... %v", tower[choice-1])
		case choice == 0:
			goto MainMenu
		default:
			goto Towers
//			return red("DATS WRONG TOWAH MAFAKA!")
		}
//		time.Sleep(3 * time.Second)
	case 5:
	Nine:
		choice = UserListInput(cfg.Env.Imagick, "Current setup", "Back")
		switch {
		case choice > 0:
			cfg.Env.Imagick[choice-1] = (cfg.Env.Imagick[choice-1])
			green("dosomething")
			time.Sleep(2 * time.Second)
			goto Nine
		default:
			goto MainMenu
		}
	case 0:
		os.Exit(0)
	default:
		red("DATS WRONG NUMBA MAFAKA!")
		time.Sleep(2 * time.Second)
		goto MainMenu
	}
	return ""
}
