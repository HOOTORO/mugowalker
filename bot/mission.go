package bot

import (
	"github.com/fatih/color"
)

func (d *Daywalker) Mission(t string) error {
	color.HiMagenta("I'M ON A MISSION!")
	tasks := d.Load(t)
	e := d.RunTasks(tasks)
	return e
}

func (d *Daywalker) RunTasks(ts []Task) error {
	for _, task := range ts {
		d.SetLocation(task.Entry)
		e := d.Do(task)
		if e != nil {
			return e
		}
	}
	return nil
}

// run user scenario([s] - path to scenario yaml)
func (d *Daywalker) Do(t Task) (e error) {
	for k, v := range t.Properties {
		for actionName, prop := range v {
			color.HiGreen("GO ACTION #%v [%v]", k, actionName)
			e = d.Action(actionName, prop)
		}
	}
	return
}
