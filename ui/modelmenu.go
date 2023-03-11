package ui

import (
	"fmt"
	"io"
	c "worker/cfg"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type selectMsg struct {
	ChosenItem string
}
type backMsg struct {
	Level uint
}

type modelSelect struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m modelSelect) Init() tea.Cmd {
	return nil
}

func (m modelSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l.Tracef(c.MgCy("<| LIST UPD. INC. |> %+v"), msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
				l.Tracef("Chosen -> %v", c.Red(m.choice))
				// cmd := inputChosen(m.choice)
				cmd := selectedItem(m.choice)
				return m, cmd
			}
			// return m.Update(msg)
		case "backspace":
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelSelect) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? Thats's cool.")
	}
	return "\n" + m.list.View()
}

func initSelectModel(li []list.Item) modelSelect {

	m := list.New(li, itemDelegate{}, defaultWidth, listHeight)
	m.Title = "Choose your destiny..."
	m.SetShowStatusBar(false)
	m.SetFilteringEnabled(false)
	m.Styles.Title = titleStyle
	m.Styles.PaginationStyle = paginationStyle
	m.Styles.HelpStyle = helpStyle

	return modelSelect{list: m}
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
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
