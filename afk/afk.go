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

type Game struct {
	Name      string
	User      *repository.User
	Active    bool
	Locations []cfg.Location
	Actions   []cfg.Action
}

func New(g, p string) *Game {
	color.HiMagenta("Launch %v!", g)
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

func (g *Game) Stage(str string) (ch, stg int) {
	stgchregex := `Stage:(?P<chapter>\d+)-(?P<stage>\d+)`

	campain := ocr.Regex(str, stgchregex)
	// campain := ocr.Regex(g.Peek(), stgchregex)

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

func Descripion(loc []*cfg.Location) (str string, usrkeys map[int]string) {
	usrkeys = make(map[int]string)
	for k, v := range loc {
		str += fmt.Sprintf("%v: %v\n", k, v.Key)
		usrkeys[k] = v.Key
	}
	return
}
