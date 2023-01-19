package afk

import (
	"fmt"
	"strings"
	"time"

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
var Fnotify func(string, string)
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

func (dw *Daywalker) CurrentLocation() ArenaLocation {
	if dw.lastLoc != nil {
		return ArenaLoc(dw.lastLoc.Key)
	}
	return 0
}
func (dw *Daywalker) TempScreenshot(name string) string {
	imgf := f("%v_%v.png", dw.fprefix, name)
	dw.lastscreenshot = cfg.TempFile(imgf)
	pt := dw.Screenshot(cfg.TempFile(imgf))
	return pt
}
func ScreenAction(r []ocr.AltoResult, act string) (x, y int) {
	for _, line := range r {
		if strings.Contains(line.Linechars, act) {
			x, y = line.X, line.Y
		}
	}
	return
}

func (dw *Daywalker) AltoRun(str string, fn func(string, string)) {
	Fnotify = fn
	taks := dw.Reactivalto(str)
	Fnotify(alto, red("TAPTARGET \n %s", taks))
	altos := dw.ScanText()
	// Fnotify(alto, f("%+v", altos))
	where := bot.GuessLocByKeywords(altos, dw.Locations)

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

func (dw *Daywalker) UpAll() {
	daily, err := dw.Daily()
	if err != nil {
		log.Errorf("Daily errr")
	}
	log.Infof("Daily done? %v", daily)
	for _, v := range dw.Tasks() {
		if availiableToday(v.Avail) {
			dw.React(&v)
		}
	}
}

func (dw *Daywalker) React(r *cfg.ReactiveTask) error {
	dw.Reactive = true
	cnt := 0
	for dw.Reactive && r.Limit >= cnt {
		txt := dw.ScanText()
		loc := bot.GuessLocByKeywords(txt, dw.Locations)
		grid, off := r.React(loc)
		before, ok := IsAction(r.Before(loc))

		if ok {
			dw.RunBefore(before)
		}

		dw.Tap(grid.X, grid.Y, off)

		after, ok := IsAction(r.After(loc))

		if ok {
			dw.RunAfter(after)
		}
		if r.Limit > 0 && r.Criteria == loc {
			cnt++
		}
	}
	return nil
}

func (dw *Daywalker) Daily() (bool, error) {
	Fnotify("|>", red("\n--> DAILY <-- \nUndone:   %08b", dw.ActiveDailies()))
	ignoredDailies := []DailyQuest{QCamp, QKT}
	for _, daily := range dw.ActiveDailies() {
		Fnotify("daily", red("--> RUN # [%s]", daily))
		if slices.Contains(ignoredDailies, daily) {
			Fnotify("daily", cyan("--> IGNORING # [%s]", daily))
			dw.MarkDone(daily)
			continue
		}
		task := dw.DailyTask(daily)
		e := dw.React(task)
		if e == nil {
			dw.MarkDone(daily)
			dw.ZeroPosition()
		}
	}
	return true, nil
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

func (dw *Daywalker) RunBefore(action Action) {
	switch action {
	case UpdProgress:
		loc := dw.CurrentLocation()
		scr := dw.Screenshot(cfg.UserFile(f("%v", loc)))
		t := ocr.TextExtract(scr)
		dw.UpdateProgress(loc, t)
		//        default:
		//            ords := cfg.StrToGrid()
		//            dw.TapGO(ords.X, ords.Y, ords.Offset)
	case RepeatX:
		task := dw.Task(dw.CurrentLocation())
		point, off := task.React(dw.CurrentLocation().String())
		for i := 0; i < 5; i++ {
			dw.Tap(point.X, point.Y, off)
			time.Sleep(time.Second)
			dw.Tap(1, 18, off)
		}

	}
}

func (dw *Daywalker) RunAfter(action Action) {
	switch action {
	case Gshot:
		time.Sleep(3 * time.Second)
		loc := dw.CurrentLocation()
		pr := dw.User.GetProgress(loc.Id())
		dw.Screenshot(cfg.UserFile(f("%v_heroinfo", pr.Level)))
		dw.Tap(1, 18, 1)
		time.Sleep(time.Second)
		dw.Tap(3, 17, 1)
	case Deactivate:
		dw.Reactive = false
	}
}
