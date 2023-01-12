package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func SelectWithTopinfo(tops interface{}) error {
	toph1 := tophedStyle.Render("AFK Worker v0.1_alpha\n####### Active setup ###########\n")
	status := fmt.Sprintf("%v\n%s", toph1, tops)

	m := InitialMenuModel() // menuModel{menulist: list.New(toplevelmenu, list.NewDefaultDelegate(), 15, 0)}
	m.header = headerStyle.Render(status) + "\n\n"
	m.menulist.Title = "Choose..."
	m.menulist.SetSize(100, 20)
	m.menulist.SetShowHelp(true)
	m.menulist.SetShowPagination(true)
	m.menulist.SetShowTitle(false)
	m.menulist.Styles.TitleBar.Height(0).Border(lipgloss.ThickBorder()).BorderForeground(hotPink)
	// m.menulist.Styles.HelpStyle.Inline(true).Height(10)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if m, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
		return err
	} else {
		if m, ok := m.(menuModel); ok && m.choice != "" {
			fmt.Printf("\n---\nHas been chosen! %s!\n", m.choice)
			//		return m.choice
		}
	}
	return nil
}

// func MultiStrInput() tea.Model {
// 	um, err := tea.NewProgram(initialUserInfoModel()).Run()
// 	if err != nil {
// 		fmt.Printf("could not start program: %s\n", err)
// 		os.Exit(1)
// 	}

// 	return um
// }

// func SoloStrInput(pl string) {
// 	p := tea.NewProgram(strInputModel(pl))
// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }
