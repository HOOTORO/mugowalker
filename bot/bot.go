package bot

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"worker/adb"
	"worker/imaginer"

	"github.com/fatih/color"
)

type Daywalker struct {
	// Represents bot
	Location Location
	job      []Mission
	locs     map[string]Location
	actions  map[string]adb.Point
	// supaloc  map[string]SupaLocation
	*adb.Device
}

// Instance of bot
func New(d *adb.Device) *Daywalker {
	return &Daywalker{
		Device:  d,
		job:     make([]Mission, 0),
		locs:    make(map[string]Location),
		actions: make(map[string]adb.Point),
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
	color.HiWhite("Recognized text: %v", text)
	return text
}

// run user scenario([s] - path to scenario yaml)
func (d *Daywalker) Mission(m string) error {
	color.HiMagenta("I'M ON A MISSION!")
	mission := d.Load(m)
	// d.save(&d.Location)
	// return errors.New("GO CHECK NEW CONF")
	firsttask := mission[0]
	if firsttask.Entry == "" {
		return errors.New("MISSION FAILED! No entry point. First Task must have enty point")
	}
	d.Location = d.locs[firsttask.Entry]
	if !d.IsFamiliarPlace(firsttask.Entry) {
		return errors.New(
			fmt.Sprintf("UNKNOWN LOCATION [%v], Please check:%v", firsttask.Entry, loccnf),
		)
	}

	if d.checkLocation() {
		for _, task := range mission {
			for k, v := range task.Plan {
				color.HiGreen("GOTO -> [%v] | CHECK?: %v", k, v)
				d.Tap(d.actions[k].X, d.actions[k].Y)
				if v {
					d.checkLocation()
				}
			}
		}
	}

	return nil
}

func (d *Daywalker) checkLocation() (r bool) {
	color.Cyan("Check location...")

	for retry := 1; !r; retry++ {
		screentext := d.Peek()
		if retry > 20 {
			color.HiRed(fmt.Sprintf("Is seems we are not in [%v]. ABORT MISSION", d.Location))
			return false
		}
		for _, v := range d.Location.Keywords {
			if strings.Contains(screentext, v) {
				color.HiGreen("Alraight we done here, go next")
				return true
			}
		}
	}
	return
}

func (d *Daywalker) IsFamiliarPlace(n string) bool {
	_, ok := d.locs[n]
	return ok
}

// func (w *Daywalker) Fight() {
// 	for {
// 		txt := w.Peek()
// 		if strings.Contains(txt, "Formations") {
// 			color.Red("\n >> FIGHT! <<\n")
// 			w.Tap(["Battle"].X, cfg["Battle"].Y)
// 		}
// 		if strings.Contains(txt, "Continue") {
// 			w.Tap(cfg["Retry"].X, cfg["Retry"].Y)
// 			color.Blue("\nPress Continue...\n")
// 		}
// 		time.Sleep(1 * time.Second)
// 	}
// }
