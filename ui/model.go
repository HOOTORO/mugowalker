package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/sirupsen/logrus"
)

type sessionState uint

const (
	selectView sessionState = iota + 1
	inputView
)

type menuModel struct {
	userstate *state
	conf      *sessionConfiguration
	session   sessionState
	showmore  bool

	menulist    list.Model
	parents     []list.Model
	choice      string
	inputChosen bool

	focusIndex int
	manyInputs []textinput.Model
	cursorMode textinput.CursorMode

	cnct multiIputModel

	winx, winy int
}

func (m menuModel) String() string {
	log.Tracef("[ options ]\n[ %v ]\n[ from yaml ]", m.conf.userSettings)
	return f(green("\n"+
		"\t|> [DevStatus : %v]\t\t [Choice : %v]\n"+
		"\t[Is input chosen? : %v]\n"+
		"\t|> [BluePid : %v]\n"+
		"\t|> userSettings --> %+v\n"+
		"\t|> Magick --> %+v\n"+
		"\t|> Tess -> %v\n"+
		"\t|> ManyInput -> %+v"),
		m.userstate.connectionStatus, m.choice, m.inputChosen, m.userstate.bluestcksPid, m.conf.userSettings, m.conf.magic, m.conf.ocr, m.manyInputs)
}

// SendTaskInfo to running tasks panel
func (m *menuModel) SendTaskInfo(task, info string) {

	m.userstate.taskch <- notify(f("%v |>", task), info)
}

type state struct {
	spinme           spinner.Model
	bluestcksPid     int
	connectionStatus int
	gameStatus       int
	taskch           chan taskinfo
	taskmsgs         []taskinfo
}

type sessionConfiguration struct {
	userSettings map[Option]string
	ocr, magic   map[string]string
}

// //////////////////////////
// /////// General //////////
// init / update / view ///
// ////////////////////////
func (m menuModel) Init() tea.Cmd {
	log.Warnf(red("\nInit model: %+v \n"), m)
	return tea.Batch(
		// textinput.Blink,
		// spinner.Tick,
		checkVM,
		activityListener(m.userstate.taskch), // wait for activity
	)
}

// ////////////////////////
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Debugf(mag("(UPD) MSG INC. -> %+v [%T]"), msg, msg)
	switch k := msg.(type) {
	case tea.KeyMsg:
		// always exit keysl
		if k.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if k.String() == "alt+s" {
			m.showmore = !m.showmore
		}

		if k.String() == "ctrl+up" {
			var cmd tea.Cmd
			m.userstate.taskch <- notify("MenuList Size |>", f("%vx%v", m.winx, m.winy))
			m.winy++
			m.menulist.SetSize(m.winx, m.winy)
			return m, cmd
		}
		if k.String() == "ctrl+down" {
			var cmd tea.Cmd
			m.userstate.taskch <- notify("MenuList Size |>", f("%vx%v", m.winx, m.winy))
			m.winy--
			m.menulist.SetSize(m.winx, m.winy)
			return m, cmd
		}
		if k.String() == "ctrl+left" {
			var cmd tea.Cmd
			m.userstate.taskch <- notify("MenuList Size |>", f("%vx%v", m.winx, m.winy))
			m.winx--
			m.menulist.SetSize(m.winx, m.winy)

			return m, cmd
		}
		if k.String() == "ctrl+right" {
			var cmd tea.Cmd
			m.userstate.taskch <- notify("MenuList Size |>", f("%vx%v", m.winx, m.winy))
			m.winx++
			m.menulist.SetSize(m.winx, m.winy)
			return m, cmd
		}
		// if k.String() == "ctrl+up" {
		// 	var cmd tea.Cmd
		// 	// w, h := menulistStyle.GetFrameSize()

		// 	m.taskch <- notify("WinSize |>", f("%vx%v", m.winx, m.winy))
		// 	m.winy += 1
		// 	m.menulist.SetSize(m.winx, m.winy)
		// 	// rt := taskinfo{Task: "eureka", Message: "Some sh! happened, ctrl-l pressed"}
		// 	// m.taskch <- rt
		// 	return m, cmd
		// }

	case taskinfo:
		k.Message = shorterer(k.Message)
		m.userstate.taskmsgs = append(m.userstate.taskmsgs[1:], k)
		return m, activityListener(m.userstate.taskch)

	case spinner.TickMsg:
		m.userstate.spinme, cmd = m.userstate.spinme.Update(msg)
		return m, cmd

	case vmStatusMsg:
		m.userstate.bluestcksPid = int(k)
		return m, cmd

	case connectionMsg:
		m.userstate.connectionStatus = int(k)
		m.session = selectView
		m.menulist.SetItems(mySettings(m))
		m.menulist.Update(msg)
		return m, cmd

	case loglevelMsg:
		m.conf.userSettings[LogLvl] = log.GetLevel().String()
		m.menulist.SetItems(mySettings(m))
		m.menulist.Update(msg)

		m.SendTaskInfo("LOG", f("LVL UPDATED to -> %v", logrus.Level(k)))

		return m, cmd

	case multiIputModel:
		m.session = inputView
		m.cnct = k
		return m.cnct.Update(msg)
	}

	log.Debugf(yellow("↓ VIEW INC ↓ \n%v"), m)

	if m.session == inputView {

	}

	if m.inputChosen {
		return updateInput(msg, m)
	}
	return updateMenu(msg, m)
}

///////////////////////////////////

func (m menuModel) View() string {
	var srt, res string
	if m.session == inputView {
		res = nlistView(m)
	}
	if m.inputChosen {
		res = inputFormView(m)
	} else {
		res = listView(m)
	}

	if m.showmore && !m.inputChosen {
		srt = m.runningTasksPanel()
		res = lipgloss.JoinHorizontal(0, res, srt)
	}

	return indent.String("\n\n"+res+"\n\n", 2) + "\n\n\n\n\r"
}

//////////////////////////////////

///////////////////////
//// helper func  ////
/////////////////////

func initTextModel(ci textinput.CursorMode, placeholder string, focus bool, prom string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetCursorMode(ci)
	ti.CursorStyle = cursorStyle
	ti.CharLimit = 0
	if focus {
		ti.Focus()
		ti.PromptStyle = focusedStyle
		ti.TextStyle = focusedStyle
	}
	// ti.Width = 30
	ti.PromptStyle.Bold(true).AlignHorizontal(1)
	ti.Prompt = f("%10s	%v ", prom, sep)
	return ti
}

func (m *menuModel) isSet(property Option) bool {
	return m.conf.userSettings[property] != ""

}
func shorterer(str string) string {
	if len(str) > 60 {
		return str[:57] + "..."
	}
	return str
}

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

func InitialMenuModel(tess, magick map[string]string, options map[Option]string) menuModel {
	const showLastTasks, x, y = 7, 100, 20

	sessionConf := &sessionConfiguration{
		userSettings: options,
		ocr:          tess,
		magic:        magick,
	}

	state := &state{
		bluestcksPid:     0,
		connectionStatus: 0,
		gameStatus:       0,
		taskch:           make(chan taskinfo),
		taskmsgs:         make([]taskinfo, showLastTasks),
		spinme:           spinner.New(),
	}
	state.spinme.Spinner = spinner.Moon
	state.spinme.Style = spinnerStyle

	m := menuModel{
		menulist:   list.New(availMenuItems(), list.NewDefaultDelegate(), x/2, y),
		parents:    nil,
		session:    selectView,
		choice:     "",
		manyInputs: make([]textinput.Model, 0),
		cursorMode: textinput.CursorStatic,

		showmore: true,

		userstate: state,
		conf:      sessionConf,

		winx: x,
		winy: y,
	}

	return m
}

func notify(ev, desc string) taskinfo {
	return taskinfo{Task: ev, Message: desc} //, Duration: time.Now()}
}
