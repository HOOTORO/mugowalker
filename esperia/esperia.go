package esperia

import "worker/navi"

var (
	Campain = &navi.Place{
		Name:  "Campain",
		Depth: 0,
		Entry: campainMenuPos,
	}
	DarkForest = &navi.Place{
		Name:  "DarkForest",
		Depth: 0,
		Entry: darkMenuPos,
	}
	Ranhorn = &navi.Place{
		Name:  "Ranhorn",
		Depth: 0,
		Entry: hornyMenuPos,
	}
)

var (
	Guild = &navi.Place{
		Name:   "Guild",
		Depth:  1,
		Entry:  guildPos,
		Parent: Ranhorn,
	}
	Shop = &navi.Place{
		Name:   "Shop",
		Depth:  1,
		Entry:  shopPos,
		Parent: Ranhorn,
	}
	OakInn = &navi.Place{
		Name:   "Oak Inn",
		Depth:  1,
		Entry:  oakPos,
		Parent: Ranhorn,
	}
)

var (
	Lab = &navi.Place{
		Name:   "Arcane Lab",
		Depth:  1,
		Entry:  labPos,
		Parent: DarkForest,
	}
	KT = &navi.Place{
		Name:   "Kings Tower",
		Depth:  1,
		Entry:  kTPos,
		Parent: DarkForest,
	}
	Bounty = &navi.Place{
		Name:   "Bounty Board",
		Depth:  1,
		Entry:  bountyPos,
		Parent: DarkForest,
	}
	Arena = &navi.Place{
		Name:   "",
		Depth:  1,
		Entry:  arenaPos,
		Parent: DarkForest,
	}
	Temporal = &navi.Place{
		Name:   "Temporal Rift",
		Depth:  1,
		Entry:  temporalPos,
		Parent: DarkForest,
	}
)

var (
	GuildHunt = &navi.Place{
		Name:   "Guildhunting",
		Depth:  2,
		Entry:  ghuntPos,
		Parent: Guild,
	}
	Hellscape = &navi.Place{
		Name:   "Hellscape",
		Depth:  2,
		Entry:  hellscpPos,
		Parent: Guild,
	}
)

var (
	TwistedRealm = &navi.Place{
		Name:   "Twisted Realm",
		Depth:  3,
		Entry:  twistedPos,
		Parent: Guild,
	}
	ClownRealm = &navi.Place{
		Name:   "Cursed Realm",
		Depth:  3,
		Entry:  clownPos,
		Parent: Guild,
	}
)

type Esperia struct {
}

type Player struct {
	Name string
	Rank int
	Stat
}

type Stat struct {
	Page int
}

type Fight struct {
}

type MultiFight struct {
}

// type Esperia struct {
// 	Campain struct {
// 		navi.Place
// 	}
// 	DarkForest struct {
// 		navi.Place
// 		KingsTower struct {
// 			navi.Place
// 			TowerOfLight struct {
// 				navi.Place
// 			}
// 			BrutalCitadel struct {
// 				navi.Place
// 			}
// 			WorldTree struct {
// 				navi.Place
// 			}
// 			ForsakenNecropolis struct {
// 				navi.Place
// 			}
// 		}
// 		ArcaneLabyrinth struct {
// 			navi.Place
// 		}
// 		ArenaOfHeroes struct {
// 			navi.Place
// 		}
// 		TemporalRift struct {
// 			navi.Place
// 		}
// 		BountyBoard struct {
// 			navi.Place
// 		}
// 	}
// 	Ranhorn struct {
// 		navi.Place
// 		OakInn struct {
// 			navi.Place
// 		}
// 		Store struct {
// 			navi.Place
// 		}
// 		Guild struct {
// 			navi.Place
// 			GHunting struct {
// 				navi.Place
// 				Wrizz struct {
// 					navi.Place
// 				}
// 				Soren struct {
// 					navi.Place
// 				}
// 			}
// 			Hellscape struct {
// 				navi.Place
// 				Cursed struct {
// 					navi.Place
// 					Leaderboardship struct {
// 						navi.Place
// 					}
// 				}
// 				Twisted struct {
// 					navi.Place
// 				}
// 			}
// 		}
// 	}
// }
// func (to *Esperia) Path(pl navi.Place) (steps []navi.TPoint) {
// 	if pl.Depth > 0 {
// 		for i := 0; i < pl.Depth; i++ {
// 			steps = append(steps, navi.TPoint{X: -1, Y: -1})
// 		}
// 	}

// 	return
// }
// func New() *Esperia {

// 	return &Esperia{}
// }
