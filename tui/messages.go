package tui

import tea "github.com/charmbracelet/bubbletea"

type inputMsg struct {
	m string
	s state
}
type inputDoneMsg string

func inputChosen(m string, st state) tea.Cmd {
	return func() tea.Msg {
		return inputMsg{m: m, s: st}
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
