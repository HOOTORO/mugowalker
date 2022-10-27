package game

import (
	"database/sql"
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
	User      *repository.User
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
	user := repository.GetUser(d.Character)
	// TODO: app check and start

	// e := d.RunTasks(tasks)
	return &Game{Name: g, Locations: locs, Active: true, b: d, User: user}
}

func (g *Game) Daily() error {
	var lastDaily repository.Daily
	// if len(g.User.Daily) == 0 {
	// 	lastDaily = repository.Daily{}
	// } else {
	// 	lastDaily =  g.User.Daily[len(g.User.Daily)-1]
	// }

	// if lastDaily.UpdatedAt.Day() > time.Now().Day(){
	// 	return errors.New(fmt.Sprintf("Daily fails! err :> %v", "Already done"))
	// }
	currentloc := g.b.WhereIs(g.Locations)

	currentloc.Actions[AFKCHEST].Run(g.b)
	currentloc.Actions["back"].Run(g.b)
	lastDaily.Loot = sql.NullBool{Bool: true}
	currentloc.Actions[FR].Run(g.b)
	currentloc.Actions["usefr"].Run(g.b)
	currentloc.Actions["back"].Run(g.b)
	currentloc.Actions["back"].Run(g.b)
	lastDaily.FastRewards = sql.NullBool{Bool: true}
	currentloc.Actions["fiends"].Run(g.b)
	currentloc.Actions["sendrecive"].Run(g.b)
	currentloc.Actions["back"].Run(g.b)
	lastDaily.Likes.Bool = true
	g.User.SaveUserInfo()
	return nil
}

func (g *Game) Push() error {
	for {

		currentloc := g.b.WhereIs(g.Locations)
		if currentloc.Name == "" {
			continue
		}
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

		stgchregex := `Stage:(?P<chapter>\d+)-(?P<stage>\d+)`
		campain := bot.Regex(g.b.Peek(), stgchregex)
		if len(campain) > 0 && campain[0] != g.User.Chapter && campain[1] != g.User.Stage {
			color.HiMagenta("##### STAGE: %v-%v ###########", campain[0], campain[1])
			g.User.Chapter = campain[0]
			g.User.Stage = campain[1]
			g.User.SaveUserInfo()

		}

		if nextMove == "screenstats" {
			bsfname, hifname := fmt.Sprintf("stats_%v-%v.png", g.User.Stage, g.User.Chapter), fmt.Sprintf("info_%v-%v.png", g.User.Stage, g.User.Chapter)
			color.HiGreen("####RUN => VICTORY, GO STATS #####\n")
			// TODO move run action to Game(?) method
			currentloc.Actions["battlestat"].Run(g.b)
			g.b.Screencap(bsfname)
			g.b.Pull(bsfname, ".")
			currentloc.Actions["heroinfo"].Run(g.b)
			g.b.Screencap(hifname)
			g.b.Pull(hifname, ".")

			currentloc.Actions["back"].Run(g.b)
			color.HiGreen("####RUN => VICTORY, NEXT #####\n")
			currentloc.Actions["next"].Run(g.b)
		} else {
			color.HiGreen("#### MOVE TO => %v #####\n", nextMove)
			currentloc.Actions[nextMove].Run(g.b)
		}

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
