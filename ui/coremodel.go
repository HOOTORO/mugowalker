package ui

import (
	c "worker/cfg"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var sessionUser *AppUser

type state uint

const (
	inputView state = iota
	selectView
)

type coreModel struct {
	current    state
	input      modelIn
	menuList   modelSelect
	parents    []*modelSelect
	showStatus bool
	quitting   bool
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perfowwrm an initial command return nil.
func (c coreModel) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (cm coreModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// l.Tracef(co.RFW("<| CORE UPD. INC. |> %+v"), msg)
	var cmd tea.Cmd
	// var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cm.menuList.list.SetWidth(msg.Width)
		return cm, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// global exit
		case "ctrl+c", "q":
			cm.quitting = true
			return cm, tea.Quit
		case "backspace":
			l.Tracef("Return to parent -> %+v", c.Green(cm.parents))
			if len(cm.parents) > 1 {
				cm.parents = cm.parents[:len(cm.parents)-1]
				cm.menuList.list.SetItems(cm.parents[len(cm.parents)-1].list.Items())
			}

		}

		switch cm.current {
		case inputView:
			in, cmd := cm.input.Update(msg)
			if in, ok := in.(modelIn); ok {
				cm.input = in
			}
			return cm, cmd

		case selectView:
			selc, cmd := cm.menuList.Update(msg)
			if selc, ok := selc.(modelSelect); ok {
				cm.menuList = selc

			}
			return cm, cmd
		}
	case inputMsg:
		switch msg {
		case "Settings":
			cm.current = selectView
			cm.menuList = initSelectModel(settings)
		case "LogLevel":
			cm.current = selectView
			cm.menuList = initSelectModel(loglvls)
		case "Device":
			cm.current = selectView
			cm.menuList = initSelectModel(device)
		case "AFK Arena":
			cm.current = selectView
			cm.menuList = initSelectModel(afkarena)
		case "Parse Image":
			cm.current = selectView
			cm.menuList = initSelectModel(imagework)
		case "Statistics":
			cm.current = selectView
			cm.menuList = initSelectModel(statistics)
		case "manual connect via ADB":
			cm.current = inputView
			// c.menu.parent = device
			cm.input = initialModel("HOST", "PORT")
		case "Account":
			cm.current = inputView
			// c.menu.parent = settings
			cm.input = initialModel("Account")
		case "App (AFK Global/Test)":
			cm.current = selectView
			cm.menuList = initSelectModel(afkgameids)
		case "Imagick":
			cm.current = inputView
			cm.input = initialModel("-colorspace",
				"-alpha",
				"-threshold",
				"-edge",
				"-negate",
				"-black-threshold")
			// c.current = userInput
			// return c.Update(msg)

		}
	case selectMsg:
		if nextItems := FindItem(msg.ChosenItem); nextItems != nil {

		}

	case inputDoneMsg:
		cm.current = selectView
		cm.menuList.choice = ""
	}
	return cm, cmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (c coreModel) View() string {
	var v string
	switch c.current {
	case inputView:
		v += c.input.View()
	case selectView:
		v += c.menuList.View()
	}

	if c.showStatus {
		v = lipgloss.JoinHorizontal(lipgloss.Top, rndBorder.Render(v), statusStyle.Render(sessionUser.String()))
	}

	res := "See you, Space Cowboy..."
	rndBorder.Render(res)
	return v
}

func initCore(in modelIn, menu modelSelect) coreModel {
	prnts := make([]*modelSelect, 0)
	prnts = append(prnts, &menu)
	return coreModel{
		showStatus: true,
		input:      in,
		parents:    prnts,
		menuList:   menu,
		current:    selectView,
		quitting:   false,
	}
}
