package ui

import tea "github.com/charmbracelet/bubbletea"

type inputMsg string
type inputDoneMsg string

func inputChosen(m string) tea.Cmd {
	return func() tea.Msg {
		return inputMsg(m)
	}
}
func inputDone(m string) tea.Cmd {
	return func() tea.Msg {
		return inputDoneMsg(m)
	}
}

func selectedItem(m string) tea.Cmd {
	return func() tea.Msg {
		return selectMsg{ChosenItem: m}
	}
}
