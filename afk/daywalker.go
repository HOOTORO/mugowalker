package afk

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"worker/afk/activities"
	"worker/bot"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"

	"golang.org/x/exp/slices"
)

const (
	alto = "ALTO"
)

type Daywalker struct {
	*bot.BasicBot
	*Game

	cnt            uint8
	ActiveTask     string
	Reactive       bool
	lastLoc        *cfg.Location
	fprefix        string
	lastscreenshot string
	maxocrtry      int
}

func NewArenaBot(b *bot.BasicBot, g *Game) *Daywalker {
	return &Daywalker{
		BasicBot: b,
		Game:     g,
		fprefix:  time.Now().Format("2006_01"),
		lastLoc:  g.GetLocation(Campain),
		cnt:      0, maxocrtry: 2,
	}
}

var f = fmt.Sprintf
var outFn func(string, string)
var red, green, cyan, ylw, mgt func(...interface{}) string

func init() {
	log = cfg.Logger()
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	ylw = color.New(color.FgHiYellow).SprintFunc()
	mgt = color.New(color.FgHiMagenta).SprintFunc()
}

func (dw *Daywalker) String() string {
	return fmt.Sprintf("Bot status:\n   Game: %v\n ActiveTask: %v\n Last Location: %v", dw.Game, dw.ActiveTask, dw.lastLoc)
}

/////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////////
func (dw *Daywalker) Location() string {
	txt := dw.ScanText()
	return bot.GuessLocation(txt, dw.Locations)

}
func (dw *Daywalker) TempScreenshot(name string) string {
	imgf := f("%v_%v.png", dw.fprefix, name)
	dw.lastscreenshot = cfg.TempFile(imgf)
	pt := dw.Screenshot(cfg.TempFile(imgf))
	return pt
}

func (dw *Daywalker) AltoRun(str string, fn func(string, string)) {
	outFn = fn
	taks := dw.Reactivalto(str)
	outFn(alto, red("TAPTARGET \n %s", taks))
	altos := dw.ScanText()
	// Fnotify(alto, f("%+v", altos))
	where := bot.GuessLocation(altos, dw.Locations)
	// dw.Daily()
	for _, r := range taks.Taptarget {
		if strings.Contains(where, r.If) {
			for _, do := range r.Do {
				x, y := bot.TextPosition(do, altos)
				if x != 0 && y != 0 {
					dw.Tap(x, y, 1)
				}
			}

		}
	}
}

func (dw *Daywalker) React(r *cfg.ReactiveTask) error {
	dw.Reactive = true
	cnt := 0
	for dw.Reactive && r.Limit >= cnt {
		txt := dw.ScanText()
		loc := bot.GuessLocation(txt, dw.Locations)
		grid, off := r.React(loc)
		// before, ok := IsAction(r.Before(loc))

		// if ok {
		// 	dw.RunBefore(before)
		// }

		dw.Tap(grid.X, grid.Y, off)

		// after, ok := IsAction(r.After(loc))

		// if ok {
		// 	dw.RunAfter(after)
		// }
		if r.Limit > 0 && r.Criteria == loc {
			cnt++
		}
	}
	return nil
}

func (dw *Daywalker) ZeroPosition() bool {
	log.Tracef("Returning to Ranhorny...")
	if dw.cnt > 5 {
		log.Errorf("Reach recursion limit, quitting....")
		return false
	}
	if dw.Location() == RANHORNY.String() {
		dw.cnt = 0
		return true
	} else {
		dw.Tap(1, 18, 1)
		dw.cnt++
		dw.ZeroPosition()
		return false
	}
}

func availiableToday(days string) bool {
	d := strings.Split(days, "/")
	weekday := time.Now().Weekday().String()
	return slices.Contains(d, weekday[:3])
}

// Press button, search for 'button's text in ocr results
func (dw *Daywalker) Press(b activities.Button) bool {
	or := dw.ScanText()
	x, y, e := LookForButton(or, b)
	if e != nil {
		return false
	}

	dw.NotifyUI(cyan("BTN PRSD"), green(f("%vx%v %v %v", x, y, b)))
	dw.Tap(x, y, 1)
	return true
}

func LookForButton(or []ocr.AltoResult, b activities.Button) (x, y int, e error) {
	for _, r := range or {
		if strings.Contains(r.Linechars, b.String()) {
			xo, yo := b.Offset()
			return r.X + xo, r.Y + yo, nil
		}
	}
	return 0, 0, errors.New("btn not found here")
}

func GetLine(or []ocr.AltoResult, n int) []ocr.AltoResult {
	var res []ocr.AltoResult
	for _, v := range or {
		if v.LineNo == n {
			res = append(res, v)
		}
	}
	return res
}

func Lines(or []ocr.AltoResult) map[int]string {
	res := make(map[int]string, 0)
	for _, v := range or {
		if val, ok := res[v.LineNo]; ok {
			res[v.LineNo] = strings.Join([]string{val, v.Linechars}, " ")
		} else {
			res[v.LineNo] = v.Linechars
		}
	}
	return res
}
