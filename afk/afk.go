package afk

import (
	"fmt"

	"worker/afk/activities"
	"worker/afk/repository"
	"worker/bot"
	"worker/cfg"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = cfg.Logger()
}

type Game struct {
	Name          string
	Active        bool
	Locations     []any
	User          *repository.User
	profile       *cfg.User
	tasks, dailys []cfg.ReactiveTask
	dailysTwo     []cfg.ReactiveAlto
}

func (g *Game) String() string {
	return fmt.Sprintf("Name: %v\n User:%v\n", g.Name, g.User.Username)
}

// New Game for a given User
func New(up *cfg.User) *Game {
	log.Infof("Launch %v", up)
	locs := make([]activities.Location, 1)
	tasks := make([]cfg.ReactiveTask, 1)
	dailys := make([]cfg.ReactiveTask, 1)
	dailysTwo := make([]cfg.ReactiveAlto, 1)

	cfg.Parse(locationsCfg, &locs)
	cfg.Parse(reactionsCfg, &tasks)
	cfg.Parse(dailyCfg, &dailys)
	cfg.Parse(dailyCfgTwo, &dailysTwo)

	anylocs := activities.AllLocations()
	for _, l := range locs {
		for _, kw := range l.Keywords() {
			if kw == "%account" {
				l.Kws = append(l.Kws, up.Account)
			}
		}
		//anylocs = append(anylocs, l)
	}

	log.Infof("Locations: %v", locs)
	log.Warnf("NEW DAILY CONF %+v", dailysTwo)

	user := repository.GetUser(up.Account)

	return &Game{
		Name:      up.Game,
		Locations: anylocs,
		Active:    true,
		User:      user,
		profile:   up,
		tasks:     tasks,
		dailys:    dailys,
		dailysTwo: dailysTwo,
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

func (g *Game) Task(loc bot.Location) *cfg.ReactiveTask {
	var Task cfg.ReactiveTask
	for _, v := range g.tasks {
		if v.Name == loc.Id() {
			return &v
		}
	}
	return &Task
}

func (g *Game) DailyTask(dly activities.DailyQuest) *cfg.ReactiveTask {
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
