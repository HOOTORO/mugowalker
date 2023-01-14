package afk

import (
	"fmt"

	"worker/afk/repository"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"
)

var (
	locations = "assets/locations.yaml"
	reactions = "assets/reactions.yaml"
	daily     = "assets/daily.yaml"
)

func Set(p, flag DailyQuest) DailyQuest {
	return p | flag
}

func Clear(p, flag DailyQuest) DailyQuest {
	return p &^ flag
}

func HasAll(p, flag DailyQuest) bool {
	return p&flag == flag
}

func HasOneOf(p, flag DailyQuest) bool {
	return p&flag != 0
}

type Game struct {
	Name          string
	Active        bool
	Locations     []cfg.Location
	User          *repository.User
	profile       *cfg.UserProfile
	tasks, dailys []cfg.ReactiveTask
}

func (g *Game) String() string {
	return fmt.Sprintf("Name: %v\n User:%v\n", g.Name, g.User.Username)
}

func New(up *cfg.UserProfile) *Game {
	color.HiMagenta("\nLaunch %v!", up)
	locs := make([]cfg.Location, 1, 1)
	tasks := make([]cfg.ReactiveTask, 1, 1)
	dailys := make([]cfg.ReactiveTask, 1, 1)

	cfg.Parse(locations, &locs)
	cfg.Parse(reactions, &tasks)
	cfg.Parse(daily, &dailys)

	user := repository.GetUser(up.Account)

	return &Game{
		Name:      up.Game,
		Locations: locs,
		Active:    true,
		User:      user,
		profile:   up,
		tasks:     tasks,
		dailys:    dailys,
	}
}

func (g *Game) GetLocation(l Location) *cfg.Location {
	for _, loc := range g.Locations {
		if loc.Key == l.String() {
			return &loc
		}
	}
	return nil
}

func (g *Game) UpdateProgress(loc Location, or ocr.Result) {
	u := g.User
	towerEx := `.*[lis|del|ght|ess|um|wer|ree](?P<floor>\d{3}|d{4}) Floors`
	stgchregex := `Stage:(?P<chapter>\d+)-(?P<stage>\d+)`

	switch loc {
	case Chapter, Stage:
		camp := or.Regex(stgchregex)
		if len(camp) == 2 {
			ch := u.GetProgress(Chapter.Id())
			ch.Update(camp[0])
			stg := u.GetProgress(Stage.Id())
			stg.Update(camp[1])
		}
	case Kings, Light, Mauler, Wilder, Graveborn, Celestial, Infernal:
		floor := or.Regex(towerEx)
		if len(floor) == 1 {
			flr := u.GetProgress(loc.Id())
			flr.Update(floor[0])
		}
	}
}

/*
	|			|
pt. |Quest  	| %b
-----------------------
10  |Loot x2	|	1
10	|FastReward	|	1
10  |Friendship	|	1
10	|Wrizz		|	1
20	|Arena1x1	|	1
10  |Inn		|	1
20	|Fight Camp	|	0
10	|Fight KT	|	1


hard to implement
10 	|Bounty		|
20	|summon		|
	|ArenaTopEnemy
	|FRqty		|
*/

func (g *Game) ActiveDailies() []DailyQuest {
	var res []DailyQuest
	userQuests := DailyQuest(g.User.DailyData().Quests)
	for i := 0; i < len(QuestNames); i++ {
		if userQuests&(1<<uint(i)) == 0 {
			res = append(res, DailyQuest(1<<uint(i)))
		}
	}
	return res
}

func (g *Game) MarkDone(quesst DailyQuest) {
	userQuests := DailyQuest(g.User.DailyData().Quests)
	if !HasOneOf(quesst, userQuests) {
		g.User.
			DailyData().
			Update(
				Set(userQuests, quesst).Id())
		color.HiRed("--> DAILY <-- \nCurrent: [%08b] \nOverall: [%08b]", quesst, g.ActiveDailies())
	}
}

func (g *Game) Task(loc Location) *cfg.ReactiveTask {
	var Task cfg.ReactiveTask
	for _, v := range g.tasks {
		if v.Name == loc.String() {
			return &v
		}
	}
	return &Task
}

func (g *Game) DailyTask(dly DailyQuest) *cfg.ReactiveTask {
	var Task cfg.ReactiveTask
	for _, v := range g.dailys {
		if v.Name == dly.String() {
			return &v
		}
	}
	return &Task
}

func (g *Game) Tasks() []cfg.ReactiveTask {
	return g.tasks
}
