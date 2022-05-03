package esperia

import "worker/navi"

//TODO: rework, map should be provided from ext or in more general way, map
var (
	Campain = &navi.Place{
		Name:  "Campain",
		Depth: 1,
		Entry: campainMenuPos,
	}
	DarkForest = &navi.Place{
		Name:  "DarkForest",
		Depth: 1,
		Entry: darkMenuPos,
	}
	Ranhorn = &navi.Place{
		Name:  "Ranhorn",
		Depth: 1,
		Entry: hornyMenuPos,
	}
)

var (
	Guild = &navi.Place{
		Name:   "Guild",
		Depth:  2,
		Entry:  guildPos,
		Parent: Ranhorn,
	}
	Shop = &navi.Place{
		Name:   "Shop",
		Depth:  2,
		Entry:  shopPos,
		Parent: Ranhorn,
	}
	OakInn = &navi.Place{
		Name:   "OakInn",
		Depth:  2,
		Entry:  oakPos,
		Parent: Ranhorn,
	}
)

var (
	Lab = &navi.Place{
		Name:   "ArcaneLab",
		Depth:  2,
		Entry:  labPos,
		Parent: DarkForest,
	}
	KT = &navi.Place{
		Name:   "KingsTower",
		Depth:  2,
		Entry:  kTPos,
		Parent: DarkForest,
	}
	Bounty = &navi.Place{
		Name:   "BountyBoard",
		Depth:  2,
		Entry:  bountyPos,
		Parent: DarkForest,
	}
	Arena = &navi.Place{
		Name:   "",
		Depth:  2,
		Entry:  arenaPos,
		Parent: DarkForest,
	}
	Temporal = &navi.Place{
		Name:   "TemporalRift",
		Depth:  2,
		Entry:  temporalPos,
		Parent: DarkForest,
	}
)

var (
	GuildHunt = &navi.Place{
		Name:   "Guildhunting",
		Depth:  3,
		Entry:  ghuntPos,
		Parent: Guild,
	}
	Hellscape = &navi.Place{
		Name:   "Hellscape",
		Depth:  3,
		Entry:  hellscpPos,
		Parent: Guild,
	}
)

var (
	TwistedRealm = &navi.Place{
		Name:   "TwistedRealm",
		Depth:  4,
		Entry:  twistedPos,
		Parent: Hellscape,
	}
	ClownRealm = &navi.Place{
		Name:   "CursedRealm",
		Depth:  4,
		Entry:  clownPos,
		Parent: Hellscape,
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
