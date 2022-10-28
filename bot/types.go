package bot

import "worker/adb"

type Daywalker struct {
	Character string
	Tasks     []Task
	*Status
	*adb.Device
}

const taskfile = "../vscode/afkarena/worker/bot/mission/task.yaml"

type Scenario struct {
	Character string
	Tasks     []Task // filepath
	Path      string
	Pattern   string //
	Duration  int
}
type Screen interface {
	Label() string
	Hits() int
	Keywords() []string
	Actions() map[string]Action
}

/*
	Representing complex action resulting location change

entry - start location key
exit - finish location
actions - base action key and custom properties
repeat - only if entry = exit
*/

type Task struct {
	Entry        string   `yaml:"entry"`
	Exit         string   `yaml:"exit"`
	NamedActions []string `yaml:"actions"`
	Repeat       int      `yaml:"repeat,omitempty"`
}

type Location struct {
	Name     string            `yaml:"name,omitempty"`
	hits     int               `yaml:"hits,omitempty"`
	Keywords []string          `yaml:"keywords"`
	Actions  map[string]Action `yaml:"actions"`
}

type Action struct {
	*adb.Point
	Destination string `yaml:"destination,omitempty"`
	Check       bool   `yaml:"check,omitempty"`
	Delay       int    `yaml:"delay,omitempty"`
	Repeat      int    `yaml:"repeat,omitempty"`
}

type Status struct {
	loc  Location
	last Action
}
