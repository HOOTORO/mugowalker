package ui

import (
	"fmt"
	"io"
	"os"
	"worker/adb"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listHeight = 20
    defaultWidth = 200
)



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

type mainModel struct {
	list     list.Model
    dev		string
	choice   string
	quitting bool
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			i, ok := m.list.SelectedItem().(menuItem)
			if ok {
				m.choice = string(i)
                if i.FilterValue() == "Devices"{
                    m.list = list.New(getDevices(), itemDelegate{}, defaultWidth, listHeight)
				} else {
			return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func getDevices() []list.Item {
	var devs []list.Item
    d, e := adb.Devices()
    if e != nil{
        devs = append(devs, menuItem("No devices found, try to connect"))
        return devs
	}
    for _, v := range d {
        devs = append(devs, menuItem(v.Serial))
    }
    return devs
}

func (m mainModel) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

func SimpleMenu(title interface{}, li []string)  {
	var items []list.Item
	for _, v := range li{
		items = append(items, menuItem(v))
	}
//		menuItem("Devices"),
//		menuItem("OCR Setup"),
//		menuItem("Tasks"),
//		menuItem("Game Locations"),





	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = fmt.Sprintf("Active setup: \n%s",title)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
//	l.Styles.TitleBar = titleBarStyle
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.Styles.StatusBar = statusStyle
	m := mainModel{list: l}
	if t, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err, "model", t)
		os.Exit(1)
	}


}
