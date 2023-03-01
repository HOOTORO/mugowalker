package ui

import "github.com/charmbracelet/bubbles/list"

type AppMenu interface {
	Entry() []list.Item
	Sub(i list.Item) []list.Item
	Parent(i list.Item) []list.Item
	ItemType(i list.Item) uint
}
