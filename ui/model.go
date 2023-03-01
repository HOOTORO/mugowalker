/*
Package ui
User Interface
*/
package ui

import (
	"errors"

	c "worker/cfg"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

var log = c.Logger()

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

// Main appmenu
type appmenu struct {
	state *state
	conf  *config

	showmore bool

	list        list.Model
	parents     []list.Model
	choice      string
	inputChosen bool

	focusIndex int
	manyInputs []textinput.Model
	cursorMode textinput.CursorMode

	input inputDialog

	winx, winy int
}

func (m appmenu) String() string {
	log.Tracef("[ options ]\n[ %v ]\n[ from yaml ]", m.conf.userSettings)
	return c.F(c.Green("\n"+
		"\t|> [DevStatus : %v]\t\t [Choice : %v]\n"+
		"\t[Is input chosen? : %v]\n"+
		"\t|> [BluePid : %v]\n"+
		"\t|> userSettings --> %+v\n"+
		"\t|> Magick --> %+v\n"+
		"\t|> Tess -> %v\n"+
		"\t|> List -> %v"),
		m.state.adbconn, m.choice, m.inputChosen, m.state.vmPid, m.conf.userSettings, m.conf.magic, m.conf.ocr, m.list)
}

// UserOutput to running tasks panel
func (m appmenu) UserOutput(task, info string) {
	m.state.taskch <- notify(c.F("%v", task), info)
}

type state struct {
	spinme     spinner.Model
	vmPid      int
	adbconn    int
	gameStatus int
	taskch     chan taskinfo
	taskstate  chan int
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
func (m appmenu) Init() tea.Cmd {
	log.Warnf(c.Red("\nInit model: %+v \n"), m)
	return tea.Batch(
		// textinput.Blink,
		// spinner.Tick,
		connectDevice(m),
		initAfk(&m),
		activityListener(m.state.taskch), // wait for activity
	)
}

// ////////////////////////
func (m appmenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Debugf(c.Mgt("(UPD) MSG INC. -> %+v [%T]"), msg, msg)
	switch k := msg.(type) {
	case tea.KeyMsg:
		// return KeyPressed(k.String(), msg, m)
		// KeyPressed(k.String(), msg, m)
		var cmd tea.Cmd
		if k.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if k.String() == "alt+s" {
			m.showmore = !m.showmore
		}

		if k.String() == "ctrl+up" {
			var cmd tea.Cmd
			m.UserOutput(c.F("%vx%v", m.winx, m.winy), " <| MenuList Size")
			m.winy++
			m.list.SetSize(m.winx, m.winy)
			return m, cmd
		}
		if k.String() == "ctrl+down" {
			var cmd tea.Cmd
			m.UserOutput(c.F("%vx%v", m.winx, m.winy), " <| MenuList Size")
			m.winy--
			m.list.SetSize(m.winx, m.winy)
			return m, cmd
		}
		if k.String() == "ctrl+left" {
			var cmd tea.Cmd
			m.UserOutput(c.F("%vx%v", m.winx, m.winy), " <| MenuList Size")
			m.winx--
			m.list.SetSize(m.winx, m.winy)

			return m, cmd
		}
		if k.String() == "ctrl+right" {

			m.UserOutput(c.F("%vx%v", m.winx, m.winy), " <| MenuList Size")
			m.winx++
			m.list.SetSize(m.winx, m.winy)
			return m, cmd
		}

	case taskinfo:
		k.Message = c.Shorterer(k.Message, 57)
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
		m.list.Title += c.F("\nShow output level |> %v", c.Cyan(log.GetLevel().String()))
		m.MenuEntry(msg)
		m.UserOutput("LOG", c.F("LVL UPDATED to -> %v", logrus.Level(k)))

		return m, cmd

	case inputDialog:
		m.state.view = inputView
		m.input = k
		return m.input.Update(msg)

	case appOnlineMsg:
		m.state.gameStatus = int(k)
		m.MenuEntry(msg)
		return m, cmd

	case prevousMenuMsg:
		prevous := m.parents[int(k)].Items()
		m.list.SetItems(prevous)
		m.list.ResetSelected()
		m.parents = m.parents[:k]
		return m, cmd
	case func(appmenu) tea.Msg:
		return m.Update(k(m))
	}

	log.Debugf(c.Ylw("↓ VIEW INC ↓ \n%v"), m)

	if m.inputChosen {
		return updateInput(msg, m)
	}
	return updateMenu(msg, m)
}

///////////////////////////////////

func (m appmenu) View() string {
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

// func (m model) View() string {
// 	if m.choice != "" {
// 		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
// 	}
// 	if m.quitting {
// 		return quitTextStyle.Render("Not hungry? That’s cool.")
// 	}
// 	return "\n" + m.list.View()
// }

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
	ti.Prompt = c.F("%10s	%v ", prom, sep)
	return ti
}

type item struct {
	title string
	child interface{}
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

// func (i item) String() string {
// 	elems := i.desc
// 	switch children := i.children.(type) {
// 	case []list.Item:
// 		for _, v := range children {
// 			elems += "<" + v.FilterValue() + sep
// 		}
// 		if elems == "" {
// 			return i.desc
// 		}
// 		return elems
// 	case textinput.Model:
// 		return children.Placeholder

// 	default:
// 		return elems
// 	}
// }

func (i item) Sub(m appmenu) []list.Item {
	switch c := i.child.(type) {
	case []list.Item:
		return c
	case func(m appmenu) []list.Item:
		return c(m)

	}
	return nil
}

func initialMenuModel(tess, magick map[string]string, options *AppUser) appmenu {
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
		taskstate:  make(chan int),
		taskmsgs:   make([]taskinfo, showLastTasks),
		spinme:     spinner.New(),
		view:       selectView,
	}
	state.spinme.Spinner = spinner.Moon
	state.spinme.Style = spinnerStyle

	delegate := list.NewDefaultDelegate()
	l := list.New(mainmenu(), delegate, x/2, y)
	log.Debugf("ENTRY LIST --> %v", c.MgCy(l))

	m := appmenu{
		list:       l, //list.New(availMenuItems(), delegate, x/2, y),
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

func (m *appmenu) MenuEntry(msg tea.Msg) {
	m.list.SetItems(mainmenu())
	m.parents = m.parents[:0]
	m.list.ResetSelected()
	m.list.Update(msg)
}

func (m *appmenu) NewList(msg tea.Msg, items []list.Item) {
	m.parents = append(m.parents, m.list)
	m.list.SetItems(items)
	m.list.ResetSelected()

	m.list.Update(msg)

}
