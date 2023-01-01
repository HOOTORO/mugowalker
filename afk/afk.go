package afk

import (
	"fmt"

	"worker/afk/repository"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"
)

var (
	locations = "cfg/config.yaml"
	actions   = "cfg/actions.yaml"
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
	Name      string
	User      *repository.User
	Active    bool
	Locations []cfg.Location
	Actions   []cfg.Action
}

func (g *Game) String() string {
	return fmt.Sprintf("Name: %v\n User:%v\n", g.Name, g.User.Username)
}

func New(g, p string) *Game {
	color.HiMagenta("\nLaunch %v!", g)
	locs := make([]cfg.Location, 1, 1)
	acts := make([]cfg.Action, 1, 1)

	cfg.Parse(locations, &locs)
	cfg.Parse(actions, &acts)
	user := repository.GetUser(p)

	return &Game{
		Name:      g,
		Locations: locs,
		Actions:   acts,
		Active:    true,
		User:      user,
	}
}

func (g *Game) GetLocation(id string) *cfg.Location {
	for _, loc := range g.Locations {
		if loc.Key == id {
			return &loc
		}
	}
	return nil
}

func (g *Game) Action(name string) *cfg.Action {
	for _, act := range g.Actions {
		if act.Name == name {
			return &act
		}
	}
	return nil
}

func (g *Game) Stage(ocr ocr.OcrResult) (ch, stg uint) {
	stgchregex := `Stage:(?P<chapter>\d+)-(?P<stage>\d+)`

	campain := ocr.Regex(stgchregex)

	progress := g.User.GetProgress()

	if len(campain) > 0 && campain[0] != int(progress.Chapter) && campain[1] != int(progress.Stage) {
		color.HiMagenta("### Campain data mismatch ###\n Actual STAGE: %v-%v ### \n >>> Fixing...", campain[0], campain[1])
		g.SetStage(campain[0], campain[1])
	}
	return progress.Chapter, progress.Chapter
}

func (g *Game) SetStage(c, s int) {
	p := g.User.GetProgress()
	p.Chapter = uint(c)
	p.Stage = uint(s)
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

func (g *Game) ActiveDailies() DailyQuest {
	return DailyQuest(g.User.ActiveQuests().Quests)
}

func (g *Game) MarkDone(quesst DailyQuest) {
	if !HasOneOf(quesst, g.ActiveDailies()) {
		g.User.
			ActiveQuests().
			Update(
				Set(g.ActiveDailies(), quesst).Indx())
		color.HiRed("--> DAILY <-- \nCurrent: [%08b] \nOverall: [%08b]", quesst, g.ActiveDailies())
	}
}

// 20 Fight Camp
// 10  Loot x2
// 10  Likes         sql.NullBool `gorm:"default:false"`
// 10  FastRewards   sql.NullBool `gorm:"default:false"`
// 10  GuildBoss     sql.NullBool `gorm:"default:false"`
// 20 Arena         sql.NullBool `gorm:"default:false"`
// 10 Fight KT
// 10  Inn sql.NullBool `gorm:"default:false"`

// hard to implement
// 10 Bounty
// 20 summon
// ArenaTopEnemy sql.NullBool `gorm:"default:false"`
//  FRqty         uint8        `gorm:"default:1"`
