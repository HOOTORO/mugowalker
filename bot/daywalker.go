package bot

import (
	"fmt"
	"strings"
	"time"

	"worker/adb"
	"worker/afk"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type Daywalker struct {
	id              uint32
	xmax, ymax, cnt uint8
	ActiveTask      string
	Reactive        bool
	lastLoc         *cfg.Location
	fprefix         string
	lastscreenshot  string
	maxocrtry       int
	*adb.Device
	*afk.Game
}

var log *logrus.Logger

func init() {
	log = cfg.Logger()
}

func (dw *Daywalker) String() string {
	return fmt.Sprintf("Bot status:\n   Game: %v\n ActiveTask: %v\n Last Location: %v", dw.Game, dw.ActiveTask, dw.lastLoc)
}

func (dw *Daywalker) CurrentLocation() afk.ArenaLocation {
	if dw.lastLoc != nil {
		return afk.ArenaLoc(dw.lastLoc.Key)
	}
	return 0
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
		loc := dw.MyLocation()
		grid, off := r.React(loc)
		before, ok := afk.IsAction(r.Before(loc))

		if ok {
			dw.RunBefore(before)
		}

		dw.TapGO(grid.X, grid.Y, off)

		after, ok := afk.IsAction(r.After(loc))

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
	color.HiRed("\n--> DAILY <-- \nUndone:   %08b", dw.ActiveDailies())
	ignoredDailies := []afk.DailyQuest{afk.QCamp, afk.QKT}
	for _, daily := range dw.ActiveDailies() {
		color.HiRed("--> RUN # [%s]", daily)
		if slices.Contains(ignoredDailies, daily) {
			color.HiCyan("--> IGNORING # [%s]", daily)
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

func (dw *Daywalker) Screenshot(name string) (string, error) {
	f := fmt.Sprintf("%v_%v_%v.png", dw.fprefix, dw.id, name)
	dw.Screencap(f)
	dw.lastscreenshot = cfg.ImageDir(f)
	err := dw.Pull(f, cfg.ImageDir(""))
	return cfg.ImageDir(f), err
}

func (dw *Daywalker) ZeroPosition() bool {
	log.Tracef("Returning to Ranhorny...")
	if dw.cnt > maxattempt {
		log.Errorf("Reach recursion limit, quitting....")
		return false
	}
	if dw.MyLocation() == afk.RANHORNY.String() {
		dw.cnt = 0
		return true
	} else {
		dw.TapGO(1, ygrid, 1)
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

func (dw *Daywalker) Gameshot(name string) string {
	f := fmt.Sprintf("%v_%v_%v.png", dw.fprefix, dw.id, name)
	dw.Screencap(f)
	_ = dw.Pull(f, cfg.UsrDir(""))
	return cfg.UsrDir(f)
}

func (dw *Daywalker) RunBefore(action afk.Action) {
	switch action {
	case afk.UpdProgress:
		loc := dw.CurrentLocation()
		scr := dw.Gameshot(fmt.Sprintf("%v", loc))
		t := ocr.TextExtract(scr)
		dw.UpdateProgress(loc, t)
		//        default:
		//            ords := cfg.StrToGrid()
		//            dw.TapGO(ords.X, ords.Y, ords.Offset)
	case afk.RepeatX:
		task := dw.Task(dw.CurrentLocation())
		point, off := task.React(dw.CurrentLocation().String())
		for i := 0; i < 5; i++ {
			dw.TapGO(point.X, point.Y, off)
			time.Sleep(time.Second)
			dw.TapGO(1, 18, off)
		}

	}
}

func (dw *Daywalker) RunAfter(action afk.Action) {
	switch action {
	case afk.Gshot:
		time.Sleep(3 * time.Second)
		loc := dw.CurrentLocation()
		pr := dw.User.GetProgress(loc.Id())
		dw.Gameshot(fmt.Sprintf("%v_heroinfo", pr.Level))
		dw.TapGO(1, 18, 1)
		time.Sleep(time.Second)
		dw.TapGO(3, 17, 1)
	case afk.Deactivate:
		dw.Reactive = false
	}
}
