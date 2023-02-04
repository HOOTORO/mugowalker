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
	NotifyUI(string, string)
}

// Push Campain script (AFK Arena)
func Push(ns Nightstalker) {
	for {
		switch ns.Location() {
		case Campain.ID:
			ns.NotifyUI("NS", "Loook for begin")
			ns.Press(Begin)
		case Bossnode.ID:
			ns.Press(BeginBoss)
		case Result.ID:
			ns.Press(TryAgain)
		case Win.ID:
			ns.Press(Continue)
		case Prepare.ID:
			ns.Press(BattleBtn)
		case RightBanner.ID:
			ns.Press(Community)
		default:
			ns.NotifyUI("NS", "Unhandled location, wait...")
			if isBaseLoc(ns.Location()) {
				ns.Press(CampainBotPanel)
			}
			time.Sleep(3 * time.Second)
		}

	}

}
