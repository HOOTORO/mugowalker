package ui

import (
	"fmt"
	"math/rand"
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
	mode       Mode
	header     string
	status     string
	devstatus  bool
	activeTask string
	showmore   bool

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
	taskch   chan taskinfo
	taskmsgs []taskinfo
}

func (m menuModel) String() string {
	log.Tracef("[ options ]\n[ %v ]\n[ from yaml ]", m.opts)
	return fmt.Sprintf(green("\n[Mode : %s ][DevStatus : %v][Quitting : %v]\n\t[Choice : %v]"), m.mode, m.devstatus, m.quitting, m.choice)
}

// responseMsg    struct{}
// responseMsgStr string
type taskinfo struct {
	Task     string
	Message  string
	Duration time.Duration
}

// A command that waits for the activity on a channel.
func activityListener(strch chan taskinfo) tea.Cmd {
	return func() tea.Msg {
		return taskinfo(<-strch)
	}
}

// processFinishedMsg is send when a pretend process completes.
type processFinishedMsg time.Duration

// pretendProcess simulates a long-running process.
func runPretendProcess() tea.Msg {
	pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond
	time.Sleep(pause)
	return processFinishedMsg(pause)
}

// generate activity
func listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			execTime := time.Second * time.Duration(rand.Intn(5))
			time.Sleep(execTime)
			sub <- struct{}{}
		}
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
		spinner.Tick,
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
			//			var cmd tea.Cmd
			//			m.quitting = true
			m.showmore = !m.showmore
			//			return m, runPretendProcess
		}
		if k.String() == "ctrl+k" {
			var cmd tea.Cmd
			//			m.quitting = true
			d := time.Second * time.Duration(rand.Intn(5))
			rt := taskinfo{Task: "eureka", Message: "Some sh! happened, ctrl-l pressed", Duration: d}
			m.activeTask = rt.Task
			m.taskch <- rt
			return m, cmd
		}
	case taskinfo:
		m.taskmsgs = append(m.taskmsgs[1:], k)
		return m, activityListener(m.taskch)
	case processFinishedMsg:
		d := time.Duration(k)
		res := taskinfo{Task: "Immitating task", Duration: d}
		// log.Printf("%s finished in %s", res.Task, res.Duration)
		m.taskmsgs = append(m.taskmsgs[1:], res)
		return m, runPretendProcess
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinme, cmd = m.spinme.Update(msg)
		return m, cmd
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
	log.Warnf(cyan("SEND TO VIEW (model) -> %v"), m)

	var srt, res string

	if m.showmore {
		srt = m.statusPanel()
	}
	if m.quitting {
		return quitStyle.Render("\n  See you later, Space Cowboy!\n\n")
	}
	switch m.mode {
	case selectList:
		res = listView(m)
	case inputMessage:
		res = inputView(m)
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
		textInput:  initTextModel("...", false, ""),
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
