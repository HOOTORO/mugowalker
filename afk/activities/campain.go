package activities

import (
	"time"
	"worker/ocr"
)

// Nightstalker can push, can climb, can daily.
// Methods required for bots to do so
type Nightstalker interface {
	OcResult() []ocr.AltoResult
	Press(Button) bool
	Location() string
	Back()
}

// Push Campain script (AFK Arena)
func Push(ns Nightstalker, out func(s, s1 string)) {
	for {
		switch ns.Location() {
		case Campain.ID:
			out("NS", "Loook for begin")
			ns.Press(Begin)
		case Bossnode.ID:
			ns.Press(Begin)
		case Result.ID:
			ns.Press(TryAgain)
		case Win.ID:
			ns.Press(Next)
		case Prepare.ID:
			ns.Press(BattleBtn)
		default:
			out("NS", "Unhandled location, wait...")
			if isBaseLoc(ns.Location()) {
				ns.Press(CampainBotPanel)
			}
			time.Sleep(3 * time.Second)
		}

	}

}
