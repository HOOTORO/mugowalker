package ui

import "fmt"

type UserData interface{}
type uModel struct {
    header string
    items []string
    choice string
    promt string
    footer string
}

func Model() *uModel {
	return nil
}

func (m *uModel) Update(ud UserData) *uModel {
    switch ud := ud.(type) {
    case uModel:
        if ud.choice !=  "" {
			m.choice = ud.choice
		}
        if ud.header !=  "" {
            m.header = ud.header
		}
        if ud.promt !=  "" {
            m.promt = ud.promt
		}
        if len(ud.items) > 0 {
            m.items = ud.items
		}
        if ud.footer !=  "" {
            m.footer = ud.footer
		}
	}
    return m
}

func (m *uModel) View() string{
    res := 
}