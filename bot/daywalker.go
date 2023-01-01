package bot

import (
	"fmt"
	"time"

	"worker/adb"
	"worker/afk"
	"worker/cfg"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type Daywalker struct {
	id              uint32
	xmax, ymax, cnt uint8
	Tasks           []cfg.Task
	ActiveTask      string
	reactive        bool
	lastLoc         *cfg.Location
	fprefix         string
	lastscreenshot  string
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

func (dw *Daywalker) CurrentLocation() string {
	if dw.lastLoc != nil {
		return dw.lastLoc.Key
	}
	return ""
}

func (dw *Daywalker) UP(r *cfg.ReactiveTask) {
	dw.reactive = true
	for dw.reactive {
		grid := r.React(dw.MyLocation())
		dw.TapGO(grid.X, grid.Y, grid.Offset)
	}
}

// Run entry point to run bot
func (dw *Daywalker) Run(t cfg.Task) {
	// availible to maintain  daily quiests
	const CanDaily = afk.Wrizz | afk.Oak | afk.Arena1x1 | afk.QKT | afk.Loot | afk.FastReward | afk.QCamp
	e := dw.runTask(t)
	if e != nil {
		panic(e)
	}
}

func (dw *Daywalker) runTask(t cfg.Task) error {
	color.HiYellow("Run task --> %v", t.Name)
	for _, a := range t.Actions {
		e := dw.Do(a)
		if e != nil {
			return e
		}
	}
	return nil
}

func (dw *Daywalker) Daily() (bool, error) {
	color.HiRed("\n--> DAILY <-- \nBDstor:   [%08b] \nConstant: [%08b]", dw.ActiveDailies(), afk.Dailies)
	if !afk.HasAll(dw.ActiveDailies(), afk.Dailies) {
		//        works!
		if !afk.HasOneOf(afk.Wrizz, dw.ActiveDailies()) {
			e := dw.Do(afk.QString(afk.Wrizz)[0])
			if e == nil {
				dw.MarkDone(afk.Wrizz)
			}
		}
		if !afk.HasOneOf(afk.Oak, dw.ActiveDailies()) {
			e := dw.Do(afk.QString(afk.Oak)[0])
			if e == nil {
				dw.MarkDone(afk.Oak)
			}
		}
		if !afk.HasOneOf(afk.Arena1x1, dw.ActiveDailies()) {
			e := dw.Do(afk.QString(afk.Arena1x1)[0])
			if e == nil {
				dw.MarkDone(afk.Arena1x1)
			}
		}
		if !afk.HasOneOf(afk.QKT, dw.ActiveDailies()) {
			e := dw.Do(afk.QString(afk.QKT)[0])
			if e == nil {
				dw.MarkDone(afk.QKT)
			}
		}
		if !afk.HasOneOf(afk.FastReward, dw.ActiveDailies()) {
			e := dw.Do(afk.QString(afk.FastReward)[0])
			if e == nil {
				dw.MarkDone(afk.FastReward)
			}
		}
		if !afk.HasOneOf(afk.QCamp, dw.ActiveDailies()) {
			e := dw.Do(afk.QString(afk.QCamp)[0])
			if e == nil {
				dw.MarkDone(afk.QCamp)
			}
		}
	}
	return true, nil
}

func (dw *Daywalker) SaveStatsFormation() {
	p := dw.User.GetProgress()
	if dw.lastLoc.Key == afk.STAT {
		filename := fmt.Sprintf("%v_%v.png", time.Now(), dw.lastLoc.Key)
		bsfname := fmt.Sprintf("stats_%v-%v_%v", p.Chapter, p.Stage, filename)
		hifname := fmt.Sprintf("info_%v-%v_%v", p.Chapter, p.Stage, filename)
		_, e := dw.Screenshot(bsfname)
		_, e = dw.Screenshot(hifname)
		dw.Back()
		if e != nil {
			log.Errorf("stst err: %v", e.Error())
		}
	}
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
	if dw.MyLocation() == afk.RANHORNY {
		dw.cnt = 0
		return true
	} else {
		dw.tapGrid(1, ygrid)
		dw.cnt++
		dw.ZeroPosition()
		return false
	}
}
