package bot

import (
	"github.com/fatih/color"
)

// run user scenario([s] - path to scenario yaml)
func (d *Daywalker) Do(t Task) (e error) {
	for k, actionName := range t.NamedActions {
		color.HiGreen("GO ACTION #%v [%v]", k, actionName)
		e = d.Action(actionName)

	}
	return
}
