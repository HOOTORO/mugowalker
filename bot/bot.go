package bot

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

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
	action.Run(d)
	return nil
}

func (a Action) Run(d *Daywalker) {
	d.Tap(a.X, a.Y)
	if a.Delay > 0 {
		delay := time.Duration(a.Delay)
		time.Sleep(delay * time.Second)
	}
}

func Keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
