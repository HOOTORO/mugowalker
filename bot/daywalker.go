package bot

import (
	"errors"
	"fmt"
	"time"

	"worker/adb"
	"worker/afk"
	"worker/cfg"

	"github.com/fatih/color"
)

type Daywalker struct {
	LastActionResult error
	Tasks            []Task
	CurrentLoc       *cfg.Location
	*adb.Device
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

func (dw *Daywalker) WalkTo(afk *afk.Game, locid string) {
	target := afk.GetLocation(locid)
	x, y := target.Position()
	dw.gridTap(x, y)
	result := isLocation(target, dw)
	if result {
		dw.CurrentLoc = target
	} else {
		color.HiMagenta("Retry walk to %v", target.Key)
		dw.WalkTo(afk, locid)
	}
}

func (dw *Daywalker) Push(game *afk.Game) error {
	actualLoc := isLocation(dw.CurrentLoc, dw)
	if !actualLoc {
		return errors.New(fmt.Sprintf("Location mismatch\nWant -> %v", dw.CurrentLoc.Key))
	}

	screentext := dw.Peek()
	c, s := game.Stage(screentext)

	color.HiGreen("> #STATE# => %v\n", dw.CurrentLoc)
	var nextMove string
	switch {
	case dw.CurrentLoc.Key == afk.ENTRY:
		nextMove = afk.CAMPBegin
	case dw.CurrentLoc.Key == afk.CAMPBegin:
		nextMove = afk.BOSSTAGE
	case dw.CurrentLoc.Key == afk.RBAN:
		nextMove = afk.CLOSE
	case dw.CurrentLoc.Key == afk.CAMPBegin:
		nextMove = afk.PUSHc
	case dw.CurrentLoc.Key == afk.BATTLE:
		nextMove = afk.FIGHT
	case dw.CurrentLoc.Key == afk.LOSE:
		nextMove = afk.RETRY
	case dw.CurrentLoc.Key == afk.PUSHc:
		nextMove = afk.BOSSTAGE
	case dw.CurrentLoc.Key == afk.CAMPWIN:
		color.HiMagenta(">> PASSED STAGE => %v-%v\n", c, s)
		game.SetStage(afk.CampainNext(c, s))
		nextMove = afk.NEXTSTAGE

		if nextMove == afk.BATTLESTAT {
			bsfname, hifname := fmt.Sprintf("stats_%v-%v.png", c, s), fmt.Sprintf("info_%v-%v.png", c, s)

			// TODO move run action to Game(?) method
			// currentloc.Actions["battlestat"].Run(g.b)
			dw.WalkTo(game, afk.BATTLESTAT)
			if dw.LastActionResult != nil {
				return dw.LastActionResult
			}

			dw.Screencap(bsfname)
			dw.Pull(bsfname, ".")

			dw.WalkTo(game, afk.HEROINFO)
			dw.Screencap(hifname)
			dw.Pull(hifname, ".")

			dw.Back()

			color.HiGreen(">> VICTORY, NEXT\n")

			game.SetStage(afk.CampainNext(c, s))
			dw.WalkTo(game, afk.NEXTSTAGE)
		} else {
			color.HiGreen("> MOVE TO => [%v] #\n", nextMove)
			dw.WalkTo(game, nextMove)
		}

	}
    return nil
}

// func (g *AfkArena) Daily() error {
// 	var lastDaily repository.Daily
// 	if len(g.User.Daily) == 0 || g.User.Daily[len(g.User.Daily)-1].UpdatedAt.Day() < time.Now().Day() {
// 		lastDaily = repository.Daily{}
// 	} else {
// 		lastDaily = g.User.Daily[len(g.User.Daily)-1]
// 	}

// 	if lastDaily.UpdatedAt.Day() == time.Now().Day() {
// 		return errors.New(fmt.Sprintf("Daily fails! err :> %v", "Already done"))
// 	}

// 	for {
// 		currentloc := g.CurrentLoc.Key
// 		if slices.Contains(BottomPanel(), currentloc) {
// 			if currentloc != CAMPBegin {
// 				g.WalkTo(CAMPBegin)
// 			}
// 			break
// 		}
// 		g.Back()
// 	}
// 	// afkchest
// 	g.WalkTo(AFKCHEST)
// 	if g.LastActionResult == nil {
// 		lastDaily.Loot.Bool = true
// 	}
// 	g.WalkTo(BACK)

// 	// Fast rewards
// 	g.WalkTo(FR)
// 	g.WalkTo(USEFR)
// 	g.WalkTo(BACK)
// 	// g.Action(BACK)

// 	if g.LastActionResult == nil {
// 		lastDaily.FastRewards.Bool = true
// 	}

// 	// frind likes
// 	g.WalkTo(RBAN)
// 	g.WalkTo(FRIENDS)
// 	g.WalkTo(LIKESBTN)
// 	g.WalkTo(BACK)

// 	g.User.SaveUserInfo()
// 	return nil
// }

// TODO: Handle POPUP Bannera, offers and guild chest
// Ofer ocr example
// ##### Where we? ##############################
// ## [Congratulations! You've completed stage 14-40! We've prepared valuable gift help you your way! Sr, Extra Purchase and receive the following rewards Bundle 01:59:28 Tap Anywhere Close] ##

func (dw *Daywalker) Battle(game *afk.Game) (bool, error) {
	if dw.CurrentLoc.Key != afk.BATTLE {
		return false, errors.New("There is nothing to battle with...")
	}
	screentext := dw.Peek()
	c, s := game.Stage(screentext)
	dw.WalkTo(game, afk.FIGHT)
	color.HiGreen("> #STATE# => %v\n", dw.CurrentLoc)
	var nextMove string
	switch {

	case dw.CurrentLoc.Key == afk.BATTLE:
		nextMove = "Fight"
	case dw.CurrentLoc.Key == afk.LOSE:
		nextMove = "Retry"
	case dw.CurrentLoc.Key == afk.BOSSTAGE:
		nextMove = "BeginBoss"
	case dw.CurrentLoc.Key == afk.CAMPWIN:
		// TODO params to control making/uploading screens
		// action to save screens: "screenstats"
		color.HiMagenta(">> PASSED STAGE => %v-%v\n", c, s)
		game.SetStage(afk.CampainNext(c, s))
		nextMove = "next"

	}

	if nextMove == "screenstats" {
		bsfname, hifname := fmt.Sprintf("stats_%v-%v.png", c, s), fmt.Sprintf("info_%v-%v.png", c, s)

		// TODO move run action to Game(?) method
		// currentloc.Actions["battlestat"].Run(g.b)
		dw.WalkTo(game, afk.BATTLESTAT)
		if dw.LastActionResult != nil {
			return false, dw.LastActionResult
		}

		dw.Screencap(bsfname)
		dw.Pull(bsfname, ".")

		dw.WalkTo(game, afk.HEROINFO)
		dw.Screencap(hifname)
		dw.Pull(hifname, ".")

		dw.Back()

		color.HiGreen(">> VICTORY, NEXT\n")

		game.SetStage(afk.CampainNext(c, s))
		dw.WalkTo(game, afk.NEXTSTAGE)
	} else {
		color.HiGreen("> MOVE TO => [%v] #\n", nextMove)
		dw.WalkTo(game, nextMove)
	}
	return true, nil
}

func (dw *Daywalker) SaveStatsFormation(g *afk.Game) {
	if dw.CurrentLoc.Key == afk.BATTLESTAT {
		filename := fmt.Sprintf("%v_%v.jpg", time.Now(), dw.CurrentLoc.Key)
		dw.Screencap("1_" + filename)
		dw.Pull("1_"+filename, ".")
		dw.WalkTo(g, afk.HEROINFO)
		dw.Screencap("2_" + filename)
		dw.Pull("2_"+filename, ".")
		dw.Back()
	}
}
