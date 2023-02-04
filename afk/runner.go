package afk

import (
	"worker/afk/activities"
	"worker/bot"
	"worker/cfg"
)

func (d *Daywalker) Run(chosen string) {
	switch chosen {
	case "Push Campain?":
		doActivity(cfg.PushCampain, d)
	case "Kings Tower":
		doActivity(cfg.ClimbKings, d)
	case "World Tree":
		doActivity(cfg.ClimbWild, d)
	case "Forsaken Necropolis":
		doActivity(cfg.ClimbGrave, d)
	case "Towers of Light":
		doActivity(cfg.ClimbLight, d)
	case "Brutal Citadel":
		doActivity(cfg.ClimbMaul, d)
	case "Celestial Sanctum":
		doActivity(cfg.ClimbCelestial, d)
	case "Infernal Fortress":
		doActivity(cfg.ClimbInferno, d)
	default:
		d.NotifyUI("RUN", "Unknown Activity")
	}
}

func doActivity(miss cfg.Mission, ns activities.Nightstalker) {
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
	default:
		outFn("RUN", "Unknown Activity")
	}
}

func Nightstalker(b *bot.BasicBot, user cfg.AppUser) *Daywalker {
	gm := New(user)
	return NewArenaBot(b, gm)
}
