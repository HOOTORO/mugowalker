package bot

import (
	"fmt"
	"time"

	"worker/adb"
	"worker/afk"
	"worker/cfg"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

var maxattempt = 5

type Daywalker struct {
	LastActionResult error
	Tasks            []Task
	lastLoc          *cfg.Location
	xmax, ymax, cnt  int
	*adb.Device
	*afk.Game
}

type Scenario struct {
	Character string
	Tasks     []Task // filepath
	Path      string
	Pattern   string //
	Duration  int
}

/*
	Representing complex action resulting location change

entry - start location key
exit - finish location
actions - base action key and custom properties
repeat - only if entry = exit
*/

type Task struct {
	Entry        string   `yaml:"entry"`
	Exit         string   `yaml:"exit"`
	NamedActions []string `yaml:"actions"`
	Repeat       int      `yaml:"repeat,omitempty"`
}

func (dw *Daywalker) WalkTo(locid string) {
	target := dw.GetLocation(locid)
	p := target.Position()
	dw.gridTap(p.X, p.Y)
	result := dw.isLocation(target.Key)
	if result {
		dw.lastLoc = target
	} else {
		color.HiMagenta("Retry walk to %v", target.Key)
		dw.WalkTo(locid)
	}
}

func (dw *Daywalker) Daily() (bool, error) {
    color.HiRed("\n--> DAILY <-- %v\nBDstor:   [%08b] \nConstant: [%08b]", afk.HasAll(afk.Dailies, dw.ActiveDailies()),dw.ActiveDailies(), afk.Dailies)
    if !afk.HasAll(dw.ActiveDailies(), afk.Dailies ) {
		color.HiBlueString("ALRIGHT! TIME FOR SOME ROUTINES TO BE DONE!")
		//        works!
		if !afk.HasOneOf(afk.Wrizz, dw.ToDo) {
			e := dw.Do(afk.GIBOSSES)
			if e == nil {
				dw.MarkDone(afk.Wrizz)
			}
		}
		if !afk.HasOneOf(afk.Oak, dw.ToDo) {
            e := dw.Do(afk.OAK)
            if e == nil {
                dw.MarkDone(afk.Oak)
            }
		}
//		if !afk.HasOneOf(afk.Arena1x1, dw.ToDo) {
//			dw.ZeroPosition()
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.gridTap(2, 18)
//			time.Sleep(time.Second)
//			dw.gridTap(4, 10)
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.gridTap(3, 18)
//			time.Sleep(time.Second)
//			dw.gridTap(3, 7)
//			time.Sleep(time.Second)
//			dw.gridTap(3, 18)
//			time.Sleep(time.Second)
//			dw.GridTapOff(4, 12, 0)
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.gridTap(3, 18)
//			time.Sleep(time.Second)
//			//skip dontwork
//			dw.GridTapOff(4, 15, 0)
//			time.Sleep(time.Second)
//			dw.MarkDone(afk.Arena1x1)
//		}
//		if !afk.HasOneOf(afk.QKT, dw.ToDo) {
//			dw.ZeroPosition()
//			time.Sleep(time.Second)
//			dw.gridTap(2, 18)
//			time.Sleep(time.Second)
//			dw.gridTap(3, 9)
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.gridTap(3, 10)
//			time.Sleep(time.Second)
//			dw.gridTap(3, 18)
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.GridTapOff(1, 15, 0)
//			time.Sleep(time.Second)
//			dw.gridTap(2, 10)
//			time.Sleep(time.Second)
//			dw.MarkDone(afk.QKT)
//		}
//		if !afk.HasOneOf(afk.QCamp, dw.ToDo) {
//			dw.ZeroPosition()
//			time.Sleep(time.Second)
//			dw.gridTap(3, 18)
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.gridTap(3, 17)
//			time.Sleep(time.Second)
//			dw.gridTap(3, 18)
//			time.Sleep(time.Second)
//			dw.gridTap(1, 15)
//			dw.MyLocation()
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.gridTap(2, 10)
//			time.Sleep(time.Second)
//			dw.gridTap(1, 18)
//			time.Sleep(time.Second)
//			dw.MarkDone(afk.QCamp)
//		}
//		if !afk.HasOneOf(afk.Loot, dw.ToDo) {
//			dw.ZeroPosition()
//			time.Sleep(time.Second)
//			dw.gridTap(3, 18)
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.gridTap(3, 16)
//			time.Sleep(time.Second)
//			dw.gridTap(3, 18)
//			dw.Peek()
//			time.Sleep(time.Second)
//			dw.gridTap(5, 17)
//			time.Sleep(time.Second)
//			dw.gridTap(4, 12)
//			time.Sleep(time.Second)
//			dw.gridTap(1, 18)
//			time.Sleep(time.Second)
//			dw.MarkDone(afk.Loot)
//			dw.MarkDone(afk.FastReward)
//		}
		//			if !daily.Loot.Bool {
		//				res := dw.Do(afk.DailyLoOt)
		//				if res == nil {
		//					daily.Loot.Bool = true
		//					dw.User.save()
		//				}
		//			}
		//			if !daily.FastRewards.Bool {
		//				res := dw.Do(afk.FastRewards)
		//				if res == nil {
		//					daily.FastRewards.Bool = true
		//					dw.User.save()
		//				}
		//			}
		//			if !daily.Likes.Bool {
		//				res := dw.Do(afk.FRIENDS)
		//				if res == nil {
		//					daily.Likes.Bool = true
		//					dw.User.save()
		//				}
		//			}
		//		}
	}
	return true, nil
}

// TODO: Handle POPUP Bannera, offers and guild chest
// Ofer ocr example
// ##### Where we? ##############################

func (dw *Daywalker) Tower(tower afk.Tower) {
	if dw.MyLocation() != afk.RANHORNY {
		dw.ZeroPosition()
	}
	dw.WalkTo(afk.DARKFORREST)
	dw.WalkTo(afk.KTower)
	switch tower {
	case afk.Kings:
		//        dw.WalkTo(afk.KTower)
		dw.gridTap(3, 10)
		dw.ktpush()
	}

}

func (dw *Daywalker) ktpush() error {
	e := dw.Do(afk.TOWERCLIMB)
	dw.gridTap(3, 10)
	dw.ktpush()
	return e
}
func (dw *Daywalker) Battle() (bool, error) {
	screentext := dw.Peek()
	c, s := dw.Stage(screentext)
    dw.gridTap(3,18)
	result := dw.Do(afk.FIGHT)
	if result != nil && dw.MyLocation() == afk.BOSSTAGE {
		dw.gridTap(3, 17)
		return false, result
	}

	if dw.isLocation(afk.WIN) {
		// TODO params to control making/uploading screens
		// action to save screens: "screenstats"
		color.HiMagenta(">> PASSED STAGE => %v-%v\n", c, s)
		dw.SetStage(afk.CampainNext(c, s))
		dw.SaveStatsFormation()
		color.HiGreen(">> VICTORY, NEXT\n")
	}

	color.HiYellow("Loooser! ")
	dw.gridTap(3, 17)
	dw.Battle()
	return true, nil
}

func (dw *Daywalker) SaveStatsFormation() {
	if dw.lastLoc.Key == afk.STAT {
		filename := fmt.Sprintf("%v_%v.png", time.Now(), dw.lastLoc.Key)
		bsfname := fmt.Sprintf("stats_%v-%v_%v", dw.User.Chapter, dw.User.Chapter, filename)
		hifname := fmt.Sprintf("info_%v-%v_%v", dw.User.Chapter, dw.User.Chapter, filename)
		e := dw.Screenshot(bsfname)
		dw.WalkTo(afk.HEROINFO)
		e = dw.Screenshot(hifname)
		dw.Back()
		if e != nil {
			log.Errorf("stst err: %v", e.Error())
		}
	}

}

func (dw *Daywalker) Screenshot(name string) error {
	dw.Screencap(name)
	err := dw.Pull(name, ".")
	return err
}

func (dw *Daywalker) ZeroPosition() bool {
	if dw.cnt > maxattempt {
		dw.cnt = 0
		log.Errorf("Reach recursion limit, quitting....")
		return false
	}
	zero := dw.GetLocation(afk.RANHORNY).Position()
    dw.lastLoc = nil
	dw.gridTap(zero.X, zero.Y)
	if dw.MyLocation() != afk.RANHORNY {
		dw.cnt++
		dw.ZeroPosition()

	}
	dw.cnt = 0
	return true
}
