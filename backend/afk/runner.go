package afk

import (
	"mugowalker/backend/afk/activities"
	"mugowalker/backend/bot"
	"mugowalker/backend/cfg"
)

func (d *Daywalker) Run(chosen string) {
	g := d.User
	switch chosen {
	case "Push Campain?":
		doActivity(cfg.PushCampain, d, g)
	case "Kings Tower":
		doActivity(cfg.ClimbKings, d, g)
	case "World Tree":
		doActivity(cfg.ClimbWild, d, g)
	case "Forsaken Necropolis":
		doActivity(cfg.ClimbGrave, d, g)
	case "Towers of Light":
		doActivity(cfg.ClimbLight, d, g)
	case "Brutal Citadel":
		doActivity(cfg.ClimbMaul, d, g)
	case "Celestial Sanctum":
		doActivity(cfg.ClimbCelestial, d, g)
	case "Infernal Fortress":
		doActivity(cfg.ClimbInferno, d, g)
	case "Do daily?":
		doActivity(cfg.Daily, d, g)
	default:
		d.NotifyUI("RUN", "Unknown Activity")
	}
}

func doActivity(miss cfg.Mission, ns activities.Nightstalker, g activities.Gamer) {
	switch miss {
	case cfg.PushCampain:
		activities.Push(ns)
	case cfg.ClimbKings:
		activities.PushTower(ns, activities.KING)
	case cfg.ClimbWild:
		activities.PushTower(ns, activities.WILDER)
	case cfg.ClimbGrave:
		activities.PushTower(ns, activities.GRAVEBORN)
	case cfg.ClimbLight:
		activities.PushTower(ns, activities.LIGHTBEARER)
	case cfg.ClimbMaul:
		activities.PushTower(ns, activities.MAULER)
	case cfg.ClimbCelestial:
		activities.PushTower(ns, activities.CELESTIAL)
	case cfg.ClimbInferno:
		activities.PushTower(ns, activities.INFERNAL)
	case cfg.Daily:
		activities.DailyRun(ns, g)
	default:
		outFn("RUN", "Unknown Activity")
	}
}

func Nightstalker(b *bot.BasicBot, user cfg.AppUser) *Daywalker {
	gm := New(user)
	return NewArenaBot(b, gm)
}
