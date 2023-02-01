/* Game service gives acces to general locations and other stuff */
package afk

import (
	"errors"
	"worker/afk/activities"
	"worker/bot"
	"worker/cfg"
	"worker/ocr"
)

var ErrLocNotFound = errors.New("unknown location")

const (
	locationsCfg = "assets/locations.yaml"
	reactionsCfg = "assets/reactions.yaml"
	dailyCfg     = "assets/daily.yaml"
	dailyCfgTwo  = "assets/dailyv2.yaml"
)

type GameWorld interface {
	Locations() []bot.Location
}

var (
	locations []activities.Location
	reactions []*cfg.Reaction
)

func init() {
	if locations == nil {
		cfg.Parse(locationsCfg, &locations)

	}

}
func Locations() []activities.Location {
	return locations
}

func LocationStruct(s string) (activities.Location, error) {
	for _, v := range locations {
		if s == v.Id() {
			return v, nil
		}
	}
	return activities.Location{}, ErrLocNotFound
}

func (g *Game) Reactivalto(str string) *cfg.ReactiveAlto {
	for _, v := range g.dailysTwo {
		if v.Name == str {
			return &v
		}
	}
	return nil
}

func (g *Game) GetLocation(l bot.Location) *activities.Location {
	for _, loc := range g.Locations {
		if lo, ok := loc.(activities.Location); ok && lo.Id() == l.Id() {
			return &lo
		}
	}
	return nil
}

func (g *Game) UpdateProgress(loc bot.Location, or ocr.Result) {
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
