/* Game service gives acces to general locations and other stuff */
package afk

import (
	"errors"
	"mugowalker/backend/afk/activities"
)

var ErrLocNotFound = errors.New("unknown location")

// ///////////////////////////
// Global afk activities ///
// /////////////////////////
type Mission int

const (
	PushCampain Mission = iota + 1
	ClimbKings
	ClimbWild
	ClimbGrave
	ClimbInferno
	ClimbMaul
	ClimbLight
	ClimbCelestial
	GuildBosses
	DailyM
)

const (
	AfkAppID     = "com.lilithgame.hgame.gp"
	AfkTestAppID = "com.lilithgames.hgame.gp.id"
)

var (
	locations []activities.Location
)

func Locations() []activities.Location {
	return locations
}
