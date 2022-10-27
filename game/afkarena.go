package game

import (
	"errors"
	"fmt"

	"worker/bot"
	"worker/game/repository"

	"github.com/fatih/color"
)

// locations
const (
	CAMPAIN     = "campain"
	BATTLE      = "battlescreen"
	CAMPLOSE    = "losecampain"
	AFKCHEST    = "afkchest"
	BOSSSTAGE   = "campainBoss"
	FR          = "fastrewards"
	DARKFORREST = "forrest"
	KT          = "kingstower"
	PVP         = "arena"
	CAMPWIN     = "campvictory"
	BATTLESTAT  = "stat"
	RANHORNY    = "ranhorn"
	GI          = "guild"
	HORNSHOP    = "shop"
	GIBOSSES    = "bosses"
)

type Game struct {
	Name      string
	Active    bool
	Locations map[string]bot.Location
	b         *bot.Daywalker
}

// type Location {
// 	name string

// }

func New(c, g string, d *bot.Daywalker) *Game {
	color.HiMagenta("Launch %v!", g)
	locs := bot.GameLocations(c)
	// TODO: app check and start

	// e := d.RunTasks(tasks)
	return &Game{Name: g, Locations: locs, Active: true, b: d}
}

func (g *Game) Daily() error {
	return errors.New(fmt.Sprintf("Daily fails! err :> %v", "ffff"))
}

func (g *Game) Push() error {
	for {
		user := repository.GetUser(g.b.Character)
		currentloc := g.b.WhereIs(g.Locations)
		color.HiGreen("#### YOU ARE HERE => %v #####\n", currentloc.Name)
		var nextMove string
		switch {
		case currentloc.Name == g.Locations[CAMPAIN].Name:
			nextMove = "BeginCampain"
		case currentloc.Name == g.Locations[BATTLE].Name:
			nextMove = "Fight"
		case currentloc.Name == g.Locations[CAMPLOSE].Name:
			nextMove = "Retry"
		case currentloc.Name == g.Locations[BOSSSTAGE].Name:
			nextMove = "BeginBoss"
		case currentloc.Name == g.Locations[CAMPWIN].Name:
			nextMove = "screenstats"

		}

		color.HiGreen("#### MOVE TO => %v#####\n", nextMove)
		stgchregex := `Stage:(?P<chapter>\d+)-(?P<stage>\d+)`
		campain := bot.Regex(g.b.Peek(), stgchregex)
		if len(campain) > 0 && campain[0] != user.Chapter && campain[1] != user.Stage {
			color.HiMagenta("##### STAGE: %v-%v ###########", campain[0], campain[1])
			user.Stage = campain[0]
			user.Chapter = campain[1]
			user.SaveUserInfo()

		}

		if nextMove == "screenstats" {
			bsfname, hifname := fmt.Sprintf("stats_%v-%v.png", user.Stage, user.Chapter), fmt.Sprintf("info_%v-%v.png", user.Stage, user.Chapter)
			color.HiGreen("####RUN => VICTORY, GO STATS #####\n")
			currentloc.Actions["battlestat"].Run(g.b)
			g.b.Screencap(bsfname)
			g.b.Pull(bsfname, ".")
			currentloc.Actions["heroinfo"].Run(g.b)
			g.b.Screencap(hifname)
			g.b.Pull(hifname, ".")

			currentloc.Actions["back"].Run(g.b)
			color.HiGreen("####RUN => VICTORY, NEXT #####\n")
			currentloc.Actions["next"].Run(g.b)
		}
		currentloc.Actions[nextMove].Run(g.b)

	}
}

func (g *Game) RunTasks(ts []bot.Task) error {
	for _, task := range ts {
		g.b.SetLocation(g.Locations[task.Entry])
		e := g.b.Do(task)
		if e != nil {
			return e
		}
	}
	return nil
}
