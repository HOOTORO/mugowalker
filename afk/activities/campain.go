package activities

import "time"

type Nightstalker interface {
	Press(Button) bool
	Location() string
	Back()
}

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
			if isBaseLoc(ns.Location()) {
				ns.Press(CampainBotPanel)
			}
			time.Sleep(10 * time.Second)
		}
		// // If Campain then bebin
		// out("NIGHTSTALKER", "ON THE HUNT [Campain]")
		// if ns.Location() == Campain.ID {
		// Butt:
		// 	if ns.Press(Begin) && ns.Location() == Bossnode.ID {
		// 		if ns.Press(Begin) {
		// 			continue
		// 		} else {
		// 			goto Butt
		// 		}
		// 	} else {
		// 		if ns.Press(BattleBtn) {
		// 			continue
		// 		} else {
		// 			goto Butt
		// 		}
		// 		// waiting battle for finish
		// 	Bzzzz:
		// 		out("NS", "Time for sleeep a lil bit")
		// 		time.Sleep(10 * time.Second)
		// 		if ns.Location() != Result.ID {
		// 			goto Bzzzz
		// 		}
		// 		out("NS", "Are we wen?")
		// 		if ns.Location() != Win.ID {
		// 			ns.Press(TryAgain)
		// 		} else {
		// 			out("NS", "o")
		// 			ns.Press(Next)
		// 			continue
		// 		}

		// 	}
		// } else {
		// 	// or Back looking for root
		// 	out("NS", "rakamakafo")
		// 	ns.Back()
		// }
		// out("NS", "run check base")
		// if isBaseLoc(ns.Location()) {
		// 	ns.Press(CampainBotPanel)
		// 	continue
		// }
	}

}
