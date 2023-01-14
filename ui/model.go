package ui

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/indent"

	"github.com/charmbracelet/bubbles/spinner"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type menuModel struct {
	mode      Mode
	header    string
	status    string
	devstatus bool

	menulist list.Model
	parents  []list.Model
	choice   string

	textInput textinput.Model

	focusIndex int
	manyInputs []textinput.Model
	cursorMode textinput.CursorMode

	response string
	spinme   spinner.Model
	quitting bool
	err      error
	cursor   int
	opts     map[string]string
}

type (
	errMsg error
)

type item struct {
	title, desc string
	children    interface{}
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.String() }
func (i item) FilterValue() string { return i.title }

func (i item) String() string {
	elems := i.desc
	switch children := i.children.(type) {
	case []list.Item:
		for _, v := range children {
			elems += "<" + v.FilterValue() + sep
		}
		if elems == "" {
			return i.desc
		}
		return elems
	case textinput.Model:
		return children.Placeholder

	default:
		return elems
	}
}

func (i item) NextLevel(m menuModel) []list.Item {
	switch c := i.children.(type) {
	case []list.Item:
		return c
	case func(m menuModel) []list.Item:
		return c(m)

	}
	return nil
}

func InitialMenuModel(userOptions map[string]string) menuModel {
	m := menuModel{
		mode:       selectList,
		header:     "Worker Setup",
		menulist:   list.New(toplevelmenu, list.NewDefaultDelegate(), 19, 0),
		parents:    nil,
		choice:     "",
		textInput:  initTextModel("...", false, ""),
		focusIndex: 0,
		manyInputs: make([]textinput.Model, 0),
		cursorMode: textinput.CursorBlink,
		quitting:   false,
		err:        nil,
		cursor:     0,
		opts:       userOptions,
	}
	// m.showStatus()
	return m
}

// //////////////////////////
// /////// General //////////
// init / update / view ///
// ////////////////////////
func (m menuModel) Init() tea.Cmd {
	log.Warnf("Init model:  \n%s", m)
	return textinput.Blink
}

// ////////////////////////
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// always exit keys
	log.Debugf(mag("UPDATE INCOMING -> %+v [%T]"), msg, msg)
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "ctrl+c" || k == "esc" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	log.Debugf(cyan("PREUPD mod -> %v"), m)
	switch m.mode {
	case selectList:
		return updateList(msg, m)
	case inputMessage:
		return updateInput(msg, m)
	case multiInput:
		return updateFormInput(msg, m)
	case runExec:
		return updateExec(msg, m)
	}
	return m, tea.Quit
}

///////////////////////////////////

func (m menuModel) View() string {
	var s string
	log.Warnf(cyan("SEND TO VIEW (model) -> %v"), m)
	if m.quitting {
		return "\n  See you later, Space Cowboy!\n\n"
	}
	switch m.mode {
	case selectList:
		s = listView(m)
	case inputMessage:
		s = inputView(m)
	case multiInput:
		s = inputFormView(m)
	case runExec:
		s = execView(m)
	}
	// r := strings.NewReplacer("[", "", "\n", "", "#", "", " ", "")
	// log.Trace("RESULTING STRING --> %50s", r.Replace(s))
	return indent.String("\n\n"+s+"\n\n", 3)
}

//////////////////////////////////

///////////////////////
//// helper func  ////
/////////////////////

func initTextModel(placeholder string, focus bool, prom string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CursorStyle = cursorStyle
	ti.CharLimit = 156
	if focus {
		ti.Focus()
		ti.PromptStyle = focusedStyle
		ti.TextStyle = focusedStyle
	}
	// ti.Width = 20
	ti.PromptStyle.Underline(true)
	ti.Prompt = prom + sep
	return ti

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
}

func (m *menuModel) showStatus() {
	var b strings.Builder
	m.status = ""
	t := fmt.Sprintf("Device --> [%v] \nProfile --> [Game: %v, User: %v]\nConnection status: ",
		m.opts[connection], m.opts[game], m.opts[account])
	b.WriteString(t)
	if m.devstatus {
		b.WriteString(green("Online"))
	} else {
		b.WriteString(red("Offline"))
	}
	m.status = statusStyle.Render(b.String())
}

func (m menuModel) isSet(property string) bool {
	if m.opts[property] != "" {
		return true
	}
	return false
}

func (m menuModel) String() string {
	log.Tracef("[ options ]\n[ %v ]\n[ from yaml ]", m.opts)
	return fmt.Sprintf(green("\n[Mode : %s ][DevStatus : %v][Quitting : %v]\n\t[Choice : %v]"), m.mode, m.devstatus, m.quitting, m.choice)
}
