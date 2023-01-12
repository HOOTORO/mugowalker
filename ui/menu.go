package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

)

// Select list w/ help and stuff
func Run(tops interface{}) error {

    status := fmt.Sprintf("AFK Worker v0.1_alpha\n####### Active setup ###########\n%s", tops)

	m := fancymodel{list: list.New(truemainmenu, list.NewDefaultDelegate(), 10, 0)}
	m.header = headerStyle.Render(status)+"\n\n"
	m.list.Title = "Choose..."

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
		return err
	}
	return nil
}

// Simple select list #2
func SelectList(t interface{}, l []string) (choice string) {

	p := tea.NewProgram(selectModel{title:fmt.Sprint(t),choices: l}, tea.WithAltScreen())

	// Run returns the selectModel as a tea.Model.
	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local selectModel and print the choice.
	if m, ok := m.(selectModel); ok && m.choice != "" {
		//        fmt.Printf("\n---\nYou chose %s!\n", m.choice)
		return m.choice
	}
	return ""
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
