package ui

import (
	"errors"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

var (
	// ErrNoAdb cant reach device via adb
	ErrNoAdb = errors.New("no ADB, setup Device")
	// ErrAppNotRunning returns when gameStatus = 0
	ErrAppNotRunning = errors.New("app not running")
)

type sessionState uint

const (
	selectView sessionState = iota + 1
	inputView
)

type menuModel struct {
	state *state
	conf  *config

	showmore bool

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
		"\t|> ManyInput -> %v"),
		m.state.adbconn, m.choice, m.inputChosen, m.state.vmPid, m.conf.userSettings, m.conf.magic, m.conf.ocr, len(m.manyInputs))
}

// UserOutput to running tasks panel
func (m menuModel) UserOutput(task, info string) {

	m.state.taskch <- notify(f("%v", task), info)
}

type state struct {
	spinme     spinner.Model
	vmPid      int
	adbconn    int
	gameStatus int
	taskch     chan taskinfo
	taskmsgs   []taskinfo
	view       sessionState
}

type config struct {
	userSettings *AppUser
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
		// checkVM,
		activityListener(m.state.taskch), // wait for activity
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
			m.UserOutput(f("%vx%v", m.winx, m.winy), " <| MenuList Size")
			m.winy++
			m.menulist.SetSize(m.winx, m.winy)
			return m, cmd
		}
		if k.String() == "ctrl+down" {
			var cmd tea.Cmd
			m.UserOutput(f("%vx%v", m.winx, m.winy), " <| MenuList Size")
			m.winy--
			m.menulist.SetSize(m.winx, m.winy)
			return m, cmd
		}
		if k.String() == "ctrl+left" {
			var cmd tea.Cmd
			m.UserOutput(f("%vx%v", m.winx, m.winy), " <| MenuList Size")
			m.winx--
			m.menulist.SetSize(m.winx, m.winy)

			return m, cmd
		}
		if k.String() == "ctrl+right" {
			var cmd tea.Cmd
			m.UserOutput(f("%vx%v", m.winx, m.winy), " <| MenuList Size")
			m.winx++
			m.menulist.SetSize(m.winx, m.winy)
			return m, cmd
		}

	case taskinfo:
		k.Message = shorterer(k.Message)
		m.state.taskmsgs = append(m.state.taskmsgs[1:], k)
		return m, activityListener(m.state.taskch)

	case spinner.TickMsg:
		m.state.spinme, cmd = m.state.spinme.Update(msg)
		return m, cmd

	case vmStatusMsg:
		m.state.vmPid = int(k)
		return m, cmd

	case connectionMsg:
		m.state.adbconn = int(k)
		m.state.view = selectView
		m.MenuEntry(msg)
		return m, cmd

	case loglevelMsg:
		m.conf.userSettings.Loglvl = log.GetLevel().String()
		m.menulist.Title += f("\nShow output level |> %v", cyan(log.GetLevel().String()))
		m.MenuEntry(msg)
		m.UserOutput("LOG", f("LVL UPDATED to -> %v", logrus.Level(k)))

		return m, cmd

	case multiIputModel:
		m.state.view = inputView
		m.cnct = k
		return m.cnct.Update(msg)

	case appOnlineMsg:
		m.state.gameStatus = int(k)
		m.MenuEntry(msg)
		return m, cmd

	case prevousMenuMsg:
		if k >= 0 {
			prevous := m.parents[int(k)].Items()
			m.menulist.SetItems(prevous)
			m.parents = m.parents[:k]
		}
		return m, cmd
	}

	log.Debugf(yellow("↓ VIEW INC ↓ \n%v"), m)

	if m.inputChosen {
		return updateInput(msg, m)
	}
	return updateMenu(msg, m)
}

///////////////////////////////////

func (m menuModel) View() string {
	var srt, res string
	if m.inputChosen {
		m.showmore = false
		res = inputFormView(m)
	} else {
		res = listView(m)
	}

	if m.showmore {
		srt = m.runningTasksPanel()
		res = lipgloss.JoinHorizontal(0, res, srt)

	}
	return lipgloss.JoinVertical(lipgloss.Center, m.runninVMs(), res)
	// return res
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

func InitialMenuModel(tess, magick map[string]string, options *AppUser) menuModel {
	const showLastTasks, x, y = 7, 100, 20

	sessionConf := &config{
		userSettings: options,
		ocr:          tess,
		magic:        magick,
	}

	state := &state{
		vmPid:      0,
		adbconn:    0,
		gameStatus: 0,
		taskch:     make(chan taskinfo),
		taskmsgs:   make([]taskinfo, showLastTasks),
		spinme:     spinner.New(),
		view:       selectView,
	}
	state.spinme.Spinner = spinner.Moon
	state.spinme.Style = spinnerStyle

	m := menuModel{
		menulist:   list.New(availMenuItems(), list.NewDefaultDelegate(), x/2, y),
		parents:    nil,
		choice:     "",
		manyInputs: make([]textinput.Model, 0),
		cursorMode: textinput.CursorStatic,

		showmore: true,

		state: state,
		conf:  sessionConf,

		winx: x,
		winy: y,
	}

	return m
}

func notify(ev, desc string) taskinfo {
	return taskinfo{Task: ev, Message: desc}
}

func (m *menuModel) MenuEntry(msg tea.Msg) {
	m.menulist.SetItems(toplevelmenu)
	m.parents = m.parents[:0]
	m.menulist.Update(msg)
}
