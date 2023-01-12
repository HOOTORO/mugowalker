package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

)


func SelectWithTopinfo(tops interface{}) error {

    status := fmt.Sprintf("AFK Worker v0.1_alpha\n####### Active setup ###########\n%s", tops)

	m := menuModel{list: list.New(toplevelmenu, list.NewDefaultDelegate(), 10, 0)}
	m.header = headerStyle.Render(status)+"\n\n"
	m.list.Title = "Choose..."
	m.list.SetShowHelp(true)

	p := tea.NewProgram(m, tea.WithAltScreen())


	if m, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
		return err
	}else {
		if m, ok := m.(menuModel); ok && m.choice != "" {
			fmt.Printf("\n---\nHas been chosen! %s!\n", m.choice)
			//		return m.choice
		}
	}
	return nil
}


func MultiStrInput() tea.Model {
	um, err := tea.NewProgram(initialUserInfoModel()).Run()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	return um
}

func SoloStrInput() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
