package activities

import (
	"time"
	c "worker/cfg"
)

type Tower uint

var towers = [...]string{"kt", "tol", "bc", "wt", "fn", "cs", "if"}

const (
	KING Tower = iota + 1
	LIGHTBEARER
	MAULER
	WILDER
	GRAVEBORN
	INFERNAL
	CELESTIAL
)

func (t Tower) String() string {
	return towers[t-1]
}

func LocLvl(s string) Tower {
	for i, v := range towers {
		if v == s {
			return Tower(i + 1)
		}
	}
	return 0
}

func (t Tower) Id() uint {
	return uint(t)
}

// Push Campain script (AFK Arena)
func PushTower(ns Nightstalker, t Tower) {
	for {
		where := ns.Location()
		log.Trace("NS", c.F("Where am I? %v", where))
		switch where {
		case Campain.ID:
			ns.Press(ForrestBotPanel)
		case Forrest.ID:
			ns.Press(King)
		case KTentrance.ID:
			ns.Press(towerBtn(t))
		case KTinside.ID, Graveborn.ID, Hypo.ID, Wilder.ID, Light.ID, Nauler.ID, Celestial.ID:
			ns.Press(Challenge)
		case Result.ID:
			ns.Press(TryAgain)
		case Win.ID:
			ns.Press(Continue)
		case Prepare.ID:
			ns.Press(BattleBtn)
		case PopoutExtra.ID:
			ns.Press(Any)
		default:
			ns.NotifyUI("NS", "Unhandled location, wait...")
			if isBaseLoc(ns.Location()) {
				ns.Press(CampainBotPanel)
			}
			time.Sleep(10 * time.Second)
		}

	}

}

func towerBtn(t Tower) Button {

	switch t.String() {
	case KTentrance.ID:
		return King
	case Wilder.ID:
		return Wld
	case Graveborn.ID:
		return Grvbrn
	case Hypo.ID:
		return Infrl
	case Celestial.ID:
		return Clstl
	case Nauler.ID:
		return Mlr
	case Light.ID:
		return Lght
	default:
		return afkbtn{name: "UNKNOWN TOWER"}
	}
}
