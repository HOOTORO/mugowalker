package bot

import (
	"errors"

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
		e := d.Run(task)
		if e != nil {
			return e
		}
	}
	return nil
}

// run user scenario([s] - path to scenario yaml)
func (d *Daywalker) Run(t Task) error {
	if !d.IsLocation() {
		return errors.New(
			color.HiRedString("LOCATION MISMATCH[%v], Please check entries", d.loc))
	}
	for k, v := range t.Properties {
		for actionName, prop := range v {
			color.HiGreen("GO ACTION #%v [%v]", k, actionName)
			d.Action(actionName, prop)
		}
	}
	return nil
}
