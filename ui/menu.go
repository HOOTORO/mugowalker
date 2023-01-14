package ui

import (
	"fmt"
	"os"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
)

func RunMainMenu(options map[string]string) error {
	m := InitialMenuModel(options)
	m.header = headerStyle.Render(header)
	m.menulist.Title = "Choose..."
	m.menulist.SetSize(110, 23)
	m.menulist.SetShowHelp(true)
	m.menulist.SetShowPagination(true)
	m.menulist.SetShowTitle(true)
	m.menulist.SetShowStatusBar(true)
	m.menulist.Styles.Title = hotStyle
	m.menulist.Styles.TitleBar = tbStyle

	p := tea.NewProgram(m, tea.WithAltScreen())

	if m, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
		return err
	} else {
		if m, ok := m.(menuModel); ok {
			a := fmt.Sprintf("\nmodel state\n---\n%+v\n", m)
			// strings.Trim()

			re := regexp.MustCompile("[(\\w*:\\w*)]")
			// res := re.ReplaceAllString(a, "\n")
			_ = re.ReplaceAllString(a, "\n")

			// r := strings.NewReplacer(":{", "\n\t") //, "}", "\n")
			fmt.Print(a)
			// fmt.Print(r.Replace(res))
		}
	}
	return nil
}
