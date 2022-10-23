package bot

import "worker/adb"

type Daywalker struct {
	Tasks []Task
	*Status
	*adb.Device
}

const taskfile = "../vscode/afkarena/worker/bot/mission/task.yaml"

type Scenario struct {
	Tasks    []Task // filepath
	Path     string
	Pattern  string //
	Duration int
}

/*
	Representing complex action resulting location change

entry - start location key
exit - finish location
actions - base action key and custom properties
repeat - only if entry = exit
*/

type Task struct {
	Entry      string                  `yaml:"entry"`
	Exit       string                  `yaml:"exit"`
	Properties []map[string]Properties `yaml:"actions"`
	Repeat     int                     `yaml:"repeat,omitempty"`
}
type Location struct {
	Label    string            `yaml:"label,omitempty"`
	hits     int               `yaml:"hits,omitempty"`
	Keywords []string          `yaml:"keywords"`
	Actions  map[string]Action `yaml:"actions"`
}

type Action struct {
	*adb.Point
	Properties
	BaseDelay int
}

type Properties struct {
	Order  int  `yaml:"order,omitempty"`
	Check  bool `yaml:"check,omitempty"`
	Delay  int  `yaml:"delay,omitempty"`
	Repeat int  `yaml:"repeat,omitempty"`
}

type Status struct {
	loc  Location
	last Action
}
