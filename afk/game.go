/* Game service gives acces to general locations and other stuff */
package afk

import (
	"errors"
	"worker/cfg"
)

var ErrLocNotFound = errors.New("unknown location")

const (
	locationsCfg = "assets/locations.yaml"
	reactionsCfg = "assets/reactions.yaml"
	dailyCfg     = "assets/daily.yaml"
	dailyCfgTwo  = "assets/dailyv2.yaml"
)

type GameWorld interface {
	Locations() []*cfg.Location
}

var (
	locations []*cfg.Location
	reactions []*cfg.Reaction
)

func init() {
	if locations == nil {
		cfg.Parse(locationsCfg, &locations)
	}

}
func Locations() []*cfg.Location {
	return locations
}

func LocationStruct(s string) (*cfg.Location, error) {
	for _, v := range locations {
		if s == v.Key {
			return v, nil
		}
	}
	return nil, ErrLocNotFound
}

func (g *Game) Reactivalto(str string) *cfg.ReactiveAlto {
	for _, v := range g.dailysTwo {
		if v.Name == str {
			return &v
		}
	}
	return nil
}
