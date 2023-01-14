package ui

import (
	"fmt"
	"io"

	"worker/cfg"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	log                           *logrus.Logger
	red, green, cyan, yellow, mag func(...interface{}) string
)

func init() {
	log = cfg.Logger()
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	mag = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()
}

const (
	listHeight   = 20
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
