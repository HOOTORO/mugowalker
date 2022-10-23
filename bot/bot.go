package bot

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"worker/adb"
	"worker/imaginer"

	"github.com/fatih/color"
)

var locs map[string]Location

// Instance of bot
func New(d *adb.Device) *Daywalker {
	return &Daywalker{
		Device: d,
		Tasks:  make([]Task, 0, 10),
		Status: &Status{},
	}
}

// OCRed Text
func (w *Daywalker) Peek() string {
	// TODO: Generate random filname
	filename := "p.png"
	w.Screencap(filename)
	e := w.Pull(filename, ".")
	if e != nil {
		fmt.Printf("\nerr:%v\nduring run:%v ", e, "w.Peek()")
	}
	abspath, _ := filepath.Abs(filename)

	text := imaginer.Text(abspath)
	return text
}

func (d *Daywalker) IsLocation() (r bool) {
	color.Cyan("Check location... %v", d.loc.Label)
	color.HiYellow("Keywords => %v", d.loc.Keywords)
	for retry := 1; !r; retry++ {
		screentext := d.Peek()

		words := OCRFields(screentext)

		if retry > 20 {
			color.HiRed(fmt.Sprintf("It seems we are not in [%v].\nWELOST! ABORT MISSION", d.loc))
			return false
		}

		color.HiCyan("OCRed Words => %v", words)
		i := 0
		for _, v := range d.loc.Keywords {
			for _, ocrtxt := range words {
				if strings.Contains(ocrtxt, v) {
					// TODO hits count move to config
					i++
					if i >= d.loc.hits {
						color.HiGreen("Location confirmed! <%v>", v)
						return true
					}
				}
			}
		}
	}
	return
}

func Keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func (d *Daywalker) AllowedAction(n string) bool {
	_, ok := locs[n]
	return ok
}

func (d *Daywalker) SetLocation(s string) {
	d.loc = locs[s]
	(d.loc).Label = s
}

func (d *Daywalker) Action(s string, props Properties) error {
	action, ok := d.loc.Actions[s]
	if !ok {
		return errors.New(fmt.Sprintf("NO Action<%v> in context<%v>!", s, d.loc.Label))
	}
	d.last = action
	action.run(d)
	return nil
}

func (a Action) run(d *Daywalker) {
	d.Tap(a.X, a.Y)
	if a.Delay > 0 {
		delay := time.Duration(a.Delay + a.BaseDelay)
		time.Sleep(delay * time.Second)
	}

	if a.Check {
		d.IsLocation()
	}
}

func OCRFields(s string) []string {
	res := strings.Fields(s)
	var filtered []string
	for _, v := range res {
		if len(v) > 2 {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
