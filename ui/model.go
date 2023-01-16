package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/muesli/reflow/indent"

	"github.com/charmbracelet/bubbles/spinner"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const showLastTasks = 10

type menuModel struct {
	bluestcksPid int
	mode         Mode
	header       string
	status       string
	devstatus    bool
	activeTask   string
	showmore     bool

	menulist list.Model
	parents  []list.Model
	choice   string

	focusIndex int
	manyInputs []textinput.Model
	cursorMode textinput.CursorMode

	response string
	spinme   spinner.Model
	quitting bool
	err      error
	cursor   int
	opts     map[string]string
	taskch   chan taskinfo
	taskmsgs []taskinfo
}

func (m menuModel) String() string {
	log.Tracef("[ options ]\n[ %v ]\n[ from yaml ]", m.opts)
	return fmt.Sprintf(green("\n[Mode : %s ] [DevStatus : %v] [Quitting : %v]\n[ ActvTask : %v] [Choice : %v]"), m.mode, m.devstatus, m.quitting, m.activeTask, m.choice)
}

// taskinfo is send when a pretend process completes.
type taskinfo struct {
	Task     string
	Message  string
	Duration time.Time
}

// A command that waits for the activity on a channel.
func activityListener(strch chan taskinfo) tea.Cmd {
	return func() tea.Msg {
		return taskinfo(<-strch)
	}
}

type (
	errMsg error
)

// //////////////////////////
// /////// General //////////
// init / update / view ///
// ////////////////////////
func (m menuModel) Init() tea.Cmd {
	log.Warnf("Init model:  \n%s", m)
	return tea.Batch(
		textinput.Blink,
		// spinner.Tick,
		// listenForActivity(m.sub), // generate activity
		activityListener(m.taskch), // wait for activity
	)
}

// ////////////////////////
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// always exit keysl
	log.Debugf(mag("UPDATE INC. -> %+v [%T]"), msg, msg)
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		if k.String() == "alt+s" {
			m.showmore = !m.showmore
		}
		if k.String() == "ctrl+k" {
			var cmd tea.Cmd
			rt := taskinfo{Task: "eureka", Message: "Some sh! happened, ctrl-l pressed", Duration: time.Now()}
			m.activeTask = rt.Task
			m.taskch <- rt
			return m, cmd
		}

	case taskinfo:
		m.mode = selectList
		m.taskmsgs = append(m.taskmsgs[1:], k)
		return m, activityListener(m.taskch)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinme, cmd = m.spinme.Update(msg)
		return m, cmd
	}

	log.Debugf(yellow("PREUPD -> %v\n%v"), m.mode, m)
	switch m.mode {
	case selectList:
		return updateList(msg, m)
	case multiInput:
		return updateInput(msg, m)
	case runExec:
		return updateExec(msg, m)
	}
	return m, tea.Quit
}

///////////////////////////////////

func (m menuModel) View() string {
	var srt, res string

	if m.showmore {
		srt = m.statusPanel()
	}
	if m.quitting {
		return quitStyle.Render("\n  See you later, Space Cowboy!\n\n")
	}
	log.Warnf(cyan("SEND TO VIEW -> %s"), m.mode)
	switch m.mode {
	case selectList:
		res = listView(m)
	case multiInput:
		res = inputFormView(m)
	case runExec:
		res = execView(m)
		// res = lipgloss.JoinVertical(lipgloss.Bottom, srt, res)
	}
	res = lipgloss.JoinHorizontal(lipgloss.Top, res, srt)

	return indent.String("\n\n"+res+"\n\n", 10)
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
}

func (m menuModel) statusPanel() string {
	var s, rt string
	m.updateStatus()
	s = m.status
	rt = fmt.Sprintf("\n"+
		m.spinme.View()+" Runing task %s...\n\n", m.activeTask)

	for _, res := range m.taskmsgs {
		if res.Task == "" {
			rt += "........................................\n"
		} else {
			rt += fmt.Sprintf("[%s]|>%s<|\n", res.Task, res.Message)
		}
	}

	rt += helpStyle.Render("\nPress 'alt+s' to hide/show this panel\n")
	//		x, y := helpStyle.GetFrameSize()
	//		rt += hotStyle.Render(fmt.Sprintf("btw:\nw: %v, h:%v", x, y))
	rt = runnunTaskStyle.Render(rt)
	return indent.String(lipgloss.JoinVertical(lipgloss.Top, s, rt), 3)
}

func (m *menuModel) updateStatus() {
	var b strings.Builder
	m.status = ""
	t := fmt.Sprintf("Device --> [%v] \nProfile --> \n\t[Game: %v]\n\t[User: %v]\nConnection status: ",
		m.opts[connection], m.opts[game], m.opts[account])
	b.WriteString(t)
	if m.devstatus {
		statusStyle.BorderForeground(brightGreen)
		b.WriteString(green("Online"))
	} else {
		statusStyle.BorderForeground(bloodRed)
		b.WriteString(red("Offline"))
	}
	b.WriteString("\n Bluestacks: ")
	if test(m) {
		statusStyle.BorderForeground(brightGreen)
		b.WriteString(green("Running"))
	} else {
		statusStyle.BorderForeground(bloodRed)
		b.WriteString(red("Shutdown"))
	}
	m.status = statusStyle.Render(b.String())
}

func (m *menuModel) isSet(property string) bool {
	if m.opts[property] != "" {
		return true
	}
	return false
}

func InitialMenuModel(userOptions map[string]string) menuModel {
	m := menuModel{
		mode:       selectList,
		showmore:   true,
		header:     "Worker Setup",
		menulist:   list.New(availMenuItems(), list.NewDefaultDelegate(), 19, 0),
		parents:    nil,
		choice:     "",
		focusIndex: 0,
		manyInputs: make([]textinput.Model, 0),
		cursorMode: textinput.CursorBlink,
		quitting:   false,
		err:        nil,
		cursor:     0,
		opts:       userOptions,
		taskch:     make(chan taskinfo),
		taskmsgs:   make([]taskinfo, showLastTasks),
		spinme:     spinner.New(),
	}
	m.updateStatus()

	m.spinme.Spinner = spinner.Moon
	m.spinme.Style = spinnerStyle
	return m
}
