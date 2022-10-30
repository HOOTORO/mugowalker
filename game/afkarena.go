package game

import (
	"errors"
	"fmt"
	"time"

	"worker/bot"
	"worker/game/repository"
	"worker/ocr"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

// locations

type Game struct {
	Name      string
	User      *repository.User
	Active    bool
	Locations map[string]bot.Location
	*bot.Daywalker
}

type Position struct {
	x, y string
}

func (p *Position) Point() (x, y string) {
	return p.x, p.y
}

// type Location {
// 	name string

// }

func New(c, g string, d *bot.Daywalker) *Game {
	color.HiMagenta("Launch %v!", g)
	locs := bot.GameLocations(c)
	user := repository.GetUser(d.Character)
	// TODO: app check and start

	// e := d.RunTasks(tasks)
	return &Game{Name: g, Locations: locs, Active: true, Daywalker: d, User: user}
}

func (g *Game) Daily() error {
	var lastDaily repository.Daily
	if len(g.User.Daily) == 0 || g.User.Daily[len(g.User.Daily)-1].UpdatedAt.Day() < time.Now().Day() {
		lastDaily = repository.Daily{}
	} else {
		lastDaily = g.User.Daily[len(g.User.Daily)-1]
	}

	if lastDaily.UpdatedAt.Day() == time.Now().Day() {
		return errors.New(fmt.Sprintf("Daily fails! err :> %v", "Already done"))
	}

	for {
		currentloc := g.WhereIs(g.Locations)
		if slices.Contains(BottomPanel(), currentloc.Name) {
			if currentloc.Name != CAMPAIN {
				g.ChangeLoc(CAMPAIN)
			}
			break
		}
		g.Back()
	}
	// afkchest

	e2 := g.Action(AFKCHEST)
	if e2 == nil {
		lastDaily.Loot.Bool = true
	}
	g.Action(BACK)

	// Fast rewards
	e1 := g.Action(FR)
	e2 = g.Action(USEFR)
	g.Action(BACK)
	// g.Action(BACK)

	if e1 == nil && e2 == nil {
		lastDaily.FastRewards.Bool = true
	}

	// frind likes
	g.ChangeLoc(RBANNER)
	g.Action(FRIENDS)
	e2 = g.Action(LIKESBTN)
	g.Action(BACK)

	g.User.SaveUserInfo()
	return nil
}

// TODO: Handle POPUP Bannera, offers and guild chest
// Ofer ocr example
// ##### Where we? ##############################
// ## [Congratulations! You've completed stage 14-40! We've prepared valuable gift help you your way! Sr, Extra Purchase and receive the following rewards Bundle 01:59:28 Tap Anywhere Close] ##
func (g *Game) Push() error {
	for {

		currentloc := g.WhereIs(g.Locations)
		if currentloc.Name == "" {
			continue
		}
		c, s := g.Stage()

		color.HiGreen("#### YOU ARE HERE => %v #####\n", currentloc.Name)
		var nextMove string
		switch {
		case currentloc.Name == g.Locations[RBANNER].Name:
			nextMove = "close"
		case currentloc.Name == g.Locations[CAMPAIN].Name:
			nextMove = "BeginCampain"
		case currentloc.Name == g.Locations[BATTLE].Name:
			nextMove = "Fight"
		case currentloc.Name == g.Locations[LOSE].Name:
			nextMove = "Retry"
		case currentloc.Name == g.Locations[BOSSSTAGE].Name:
			nextMove = "BeginBoss"
		case currentloc.Name == g.Locations[CAMPWIN].Name:
			// TODO params to control making/uploading screens
			// action to save screens: "screenstats"
			color.HiMagenta("#### PASSED STAGE => %v-%v #######\n", c, s)
			g.SetStage(CampainNext(c, s))
			nextMove = "next"

		}

		if nextMove == "screenstats" {
			bsfname, hifname := fmt.Sprintf("stats_%v-%v.png", c, s), fmt.Sprintf("info_%v-%v.png", c, s)

			// TODO move run action to Game(?) method
			// currentloc.Actions["battlestat"].Run(g.b)
			err := g.Action(BATTLESTAT)
			if err != nil {
				return err
			}
			g.SetLocation(g.Locations[BATTLESTAT])
			g.Screencap(bsfname)
			g.Pull(bsfname, ".")

			g.Action(HEROINFO)
			g.Screencap(hifname)
			g.Pull(hifname, ".")

			g.Action("back")
			g.SetLocation(g.Locations[CAMPWIN])
			color.HiGreen("#### RUN => VICTORY, NEXT #####\n")

			g.SetStage(CampainNext(c, s))
			g.Action("next")
		} else {
			color.HiGreen("#### MOVE TO => %v #####\n", nextMove)
			g.Action(nextMove)
		}

	}
}

func (g *Game) Stage() (ch, stg int) {
	stgchregex := `Stage:(?P<chapter>\d+)-(?P<stage>\d+)`

	campain := ocr.Regex(g.Peek(), stgchregex)

	if len(campain) > 0 && campain[0] != g.User.Chapter && campain[1] != g.User.Stage {
		color.HiMagenta("### Campain data mismatch ###\n Actual STAGE: %v-%v ### \n >>> Fixing...", campain[0], campain[1])
		g.SetStage(campain[0], campain[1])
	}
	return g.User.Chapter, g.User.Stage
}

func (g *Game) SetStage(c, s int) {
	g.User.Chapter = c
	g.User.Stage = s
	g.User.SaveUserInfo()
}

func CampainNext(c, s int) (int, int) {
	if s <= stages40 {
		if s < 40 {
			s++
		} else {
			s = 1
			c++

		}
	} else {
		if s < 60 {
			s++
		} else {
			s = 1
			c++
		}
	}
	return c, s
}

// func (g *Game) RunTasks(ts []bot.Task) error {
// 	for _, task := range ts {
// 		g.SetLocation(g.Locations[task.Entry])
// 		e := g.Do(task)
// 		if e != nil {
// 			return e
// 		}
// 	}
// 	return nil
// }

func (g *Game) ChangeLoc(s string) {
	err := g.Action(s)
	if err != nil {
		panic("Cannot Action! >_< " + s + "\n" + err.Error())
	}
	g.Actlike(&Position{})
	g.SetLocation(g.Locations[s])
}
