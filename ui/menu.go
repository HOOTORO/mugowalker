package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle            = lipgloss.NewStyle().Margin(3,3)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func Run(tops interface{}) error {

    status := fmt.Sprintf("AFK Worker v0.1_alpha\n####### Active setup ###########\n%s", tops)

	m := fancymodel{list: list.New(truemainmenu, list.NewDefaultDelegate(), 0, 0)}
    m.header = status
	m.list.Title = "Choose..."

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
		return err
	}
	return nil
}

func SelectList(l []string) (choice string) {
	p := tea.NewProgram(selectModel{choices: l})

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

func TxtIn() tea.Model {
	um, err := tea.NewProgram(initialUserInfoModel()).Run()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	return um
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
