/* Game service gives acces to general locations and other stuff */
package afk

import (
	"errors"
	"worker/afk/activities"
)

var ErrLocNotFound = errors.New("unknown location")

const (
	AfkAppID     = "com.lilithgames.hgame.gp"
	AfkTestAppID = "com.lilithgames.hgame.gp.id"
)

const (
	reactionsCfg = "assets/reactions.yaml"
	dailyCfg     = "assets/daily.yaml"
	dailyCfgTwo  = "assets/dailyv2.yaml"
)

var (
	locations []activities.Location
)

func Locations() []activities.Location {
	return locations
}
