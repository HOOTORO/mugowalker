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

type Daywalker struct {
    id              uint32
    xmax, ymax, cnt uint8
    Tasks           []Task
    ActiveTask      string
    lastLoc         *cfg.Location
    fprefix         string
    lastscreenshot  string
    *adb.Device
    *afk.Game
}

func (dw *Daywalker) String() string{
    return fmt.Sprintf("Bot status:\n   Game: %v\n ActiveTask: %v\n Last Location: %v", dw.Game, dw.ActiveTask, dw.lastLoc)
}

func (dw *Daywalker) CurrentLocation()string{
    if dw.lastLoc != nil{
        return dw.lastLoc.Key
    }
    return ""
}

/*
	Representing complex action resulting location change

entry - start location key
exit - finish location
actions - base action key and custom properties
repeat - only if entry = exit
*/

type Task struct {
    Name         string   `yaml:"name"`
    Actions []string `yaml:"actions"`
    Repeat       int      `yaml:"repeat,omitempty"`
}

// UP entry point to run bot
func (dw *Daywalker) UP(t Task) {
    //availible to maintain  daily quiests
    const CanDaily = afk.Wrizz | afk.Oak | afk.Arena1x1 | afk.QKT | afk.Loot | afk.FastReward | afk.QCamp

//    b, e := dw.Daily()
    e := dw.runTask(t)
    if e!= nil{
        panic(e)
    }


}

func (dw *Daywalker) runTask(t Task) error {
    color.HiYellow("Run task --> %v", t.Name )
    for _, a := range t.Actions{
        e := dw.Do(a)
        if e != nil {
            return e
        }
    }
    return nil
}
// Deprecated
func (dw *Daywalker) WalkTo(locid string) {
    target := dw.GetLocation(locid)
    p := target.Position()
    dw.tapGrid(p.X, p.Y)

    if dw.MyLocation() == target.Key {
        dw.lastLoc = target
    } else {
        color.HiMagenta("Retry walk to %v", target.Key)
        dw.WalkTo(locid)
    }
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

// TODO: Handle POPUP Bannera, offers and guild chest
// Ofer ocr example
// ##### Where we? ##############################

func (dw *Daywalker) Tower(tower afk.Tower) {
    color.HiMagenta("Startin tower...")
    if !dw.ZeroPosition() {
        dw.ZeroPosition()
    }
    //	dw.WalkTo(afk.DARKFORREST)
    //	dw.WalkTo(afk.KTower)

    switch tower {
    case afk.Kings:
        err := dw.Do(afk.GO2KT)
        if err != nil {
            return
        }
        dw.ktpush()
    }

}

func (dw *Daywalker) ktpush() error {
    color.HiMagenta("Push KT tower...")

    e := dw.Do(afk.DOTOWERCLIMB)
    if e != nil {
        return e
    }
    if dw.lastLoc.Key == afk.TowerInside{
        e = dw.ktpush()
    } else {
        p := dw.lastLoc.Position()
        dw.tapGrid(p.X, p.Y)
        e = dw.ktpush()
    }

    return e
}

func (dw *Daywalker) SaveStatsFormation() {
    if dw.lastLoc.Key == afk.STAT {
        filename := fmt.Sprintf("%v_%v.png", time.Now(), dw.lastLoc.Key)
        bsfname := fmt.Sprintf("stats_%v-%v_%v", dw.User.Chapter, dw.User.Chapter, filename)
        hifname := fmt.Sprintf("info_%v-%v_%v", dw.User.Chapter, dw.User.Chapter, filename)
        _, e := dw.Screenshot(bsfname)
        dw.WalkTo(afk.HEROINFO)
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
