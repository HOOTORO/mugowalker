package bot

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"worker/adb"
	"worker/ocr"
)

// Instance of bot
func New(d *adb.Device, ch string) *Daywalker {
	return &Daywalker{
		Character: ch,
		Device:    d,
		Tasks:     make([]Task, 0, 10),
		Status:    &Status{},
	}
}

// OCRed Text
// TODO: maybe add args to peek like peek(data interface)
// smth like this should be w.Peek(Location) \n w.Peek(Stage)
func (w *Daywalker) Peek() string {
	// TODO: Generate random filname
	filename := "p.png"
	w.Screencap(filename)
	e := w.Pull(filename, ".")
	if e != nil {
		fmt.Printf("\nerr:%v\nduring run:%v ", e, "Peek()")
	}
	abspath, _ := filepath.Abs(filename)

	text := ocr.Text(abspath)
	return text
}

func (d *Daywalker) SetLocation(l Location) {
	d.loc = l
}

func (d *Daywalker) Action(s string) error {
	action, ok := d.loc.Actions[s]
	if !ok {
		return errors.New(fmt.Sprintf("NO Action<%v> in context<%v>!", s, d.loc.Name))
	}
	d.last = action
	action.run(d)
	// d.SetLocation(actionD)
	return nil
}

func (a Action) run(d *Daywalker) {
	d.Tap(a.X, a.Y)
	if a.Delay > 0 {
		delay := time.Duration(a.Delay)
		time.Sleep(delay * time.Second)
	}
}

// run user scenario([s] - path to scenario yaml)
func (d *Daywalker) Do(t Task) (e error) {
	for k, actionName := range t.NamedActions {
		color.HiGreen("GO ACTION #%v [%v]", k, actionName)
		e = d.Action(actionName)

	}
	return
}

func (d *Daywalker) WhereIs(locs map[string]Location) Location {
	current := ocr.OCRFields(d.Peek())
	color.HiYellow("### Where we? ###\n ## %v ## \n", current)
	maxhits, locName := 0, ""
	for name, v := range locs {
		hits := ocr.KeywordHits(v.Keywords, current)
		if hits > maxhits {
			maxhits = hits
			locName = name
		}

	}
	if locName != "" {
		color.HiBlue("### %v ###\n", locName)
	}
	d.SetLocation(locs[locName])
	return locs[locName]
}
