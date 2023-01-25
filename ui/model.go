package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
)

const showLastTasks = 7

type menuModel struct {
	bluestcksPid     int
	connectionStatus int

	statusInfo string
	header     string
	showmore   bool

	menulist    list.Model
	parents     []list.Model
	choice      string
	inputChosen bool

	focusIndex int
	manyInputs []textinput.Model
	cursorMode textinput.CursorMode

	spinme         spinner.Model
	quitting       bool
	usersettings   map[string]string
	usersettingsv2 map[Option]string
	taskch         chan taskinfo
	taskmsgs       []taskinfo
	winx, winy     int
}

func (m menuModel) String() string {
	log.Tracef("[ options ]\n[ %v ]\n[ from yaml ]", m.usersettings)
	return f(green("\n\t|> [DevStatus : %v]\t[Quitting : %v]\n\t|> [Choice : %v]\t[BluePid : %v]\n\t|> userSettings --> %+v\n\t|> userSettingsv2 --> %+v"), m.connectionStatus, m.quitting, m.choice, m.bluestcksPid, m.usersettings, m.usersettingsv2)
}

// //////////////////////////
// /////// General //////////
// init / update / view ///
// ////////////////////////
func (m menuModel) Init() tea.Cmd {
	log.Warnf(red("\nInit model: %+v \n"), m)
	return tea.Batch(
		textinput.Blink,
		// spinner.Tick,
		checkVM,
		activityListener(m.taskch), // wait for activity
	)
}

// ////////////////////////
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch k := msg.(type) {
	case tea.KeyMsg:
		log.Debugf(mag("(UPD) KEY INC. -> %+v [%T]"), msg, msg)
		// always exit keysl
		if k.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		if k.String() == "alt+s" {
			m.showmore = !m.showmore
		}

		if k.String() == "ctrl+up" {
			var cmd tea.Cmd
			m.taskch <- notify("MenuList Size |>", f("%vx%v", m.winx, m.winy))
			m.winy++
			m.menulist.SetSize(m.winx, m.winy)
			return m, cmd
		}
		if k.String() == "ctrl+down" {
			var cmd tea.Cmd
			m.taskch <- notify("MenuList Size |>", f("%vx%v", m.winx, m.winy))
			m.winy--
			m.menulist.SetSize(m.winx, m.winy)
			return m, cmd
		}
		if k.String() == "ctrl+left" {
			var cmd tea.Cmd
			m.taskch <- notify("MenuList Size |>", f("%vx%v", m.winx, m.winy))
			m.winx--
			m.menulist.SetSize(m.winx, m.winy)

			return m, cmd
		}
		if k.String() == "ctrl+right" {
			var cmd tea.Cmd
			m.taskch <- notify("MenuList Size |>", f("%vx%v", m.winx, m.winy))
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
		m.taskmsgs = append(m.taskmsgs[1:], k)
		return m, activityListener(m.taskch)

	case spinner.TickMsg:
		m.spinme, cmd = m.spinme.Update(msg)
		return m, cmd

	case vmStatusMsg:
		m.bluestcksPid = int(k)
		return m, cmd

	case connectionMsg:
		m.connectionStatus = int(k)
		return m, cmd
	case tea.WindowSizeMsg:
		// w, h := menulistStyle.GetFrameSize()
		m.winx = k.Width
		m.winy = k.Height
		m.taskch <- notify("WinSize |>", f("%vx%v", k.Width, k.Height))
		m.menulist.SetSize(k.Width, k.Height)
	}

	log.Debugf(yellow("\nVIEW INC -> %v\n%v"), m)

	if m.inputChosen {
		return updateInput(msg, m)
	}
	return updateMenu(msg, m)
}

///////////////////////////////////

func (m menuModel) View() string {
	var srt, res string

	if m.quitting {
		return quitStyle.Render("\n  See you later, Space Cowboy!\n\n")
	}

	if m.inputChosen {
		res = inputFormView(m)
	} else {
		res = listView(m)
	}

	if m.showmore {
		srt = m.runningTasksPanel()
		res = lipgloss.JoinHorizontal(lipgloss.Top, res, srt)
	}

	return indent.String("\n\n"+res+"\n\n", 2)
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
	ti.Width = 30
	ti.PromptStyle.Underline(true)
	ti.Prompt = prom + sep
	return ti
}

func (m *menuModel) isSet(property string) bool {
	return m.usersettings[property] != ""

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

type menuItem string

func (i menuItem) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}
