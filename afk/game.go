/* Game service gives acces to general locations and other stuff */
package afk

import (
	"errors"
	"worker/afk/activities"
	"worker/cfg"
)

var ErrLocNotFound = errors.New("unknown location")

const (
	locationsCfg = "assets/locations.yaml"
	reactionsCfg = "assets/reactions.yaml"
	dailyCfg     = "assets/daily.yaml"
	dailyCfgTwo  = "assets/dailyv2.yaml"
)

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
