package esperia

import "worker/navi"

var flatmap map[string]*navi.Location = make(map[string]*navi.Location, 0)

//TODO: rework, map should be provided from ext or in more general way, map

var (
	Main = &navi.Location{
		Name:  "Main",
		Depth: 0,
		Entry: campainMenuPos,
	}
)

var (
	Campain = &navi.Location{
		Name:  "Campain",
		Depth: 1,
		Entry: campainMenuPos,
	}

	DarkForest = &navi.Location{
		Name:  "DarkForest",
		Depth: 1,
		Entry: darkMenuPos,
	}
	Ranhorn = &navi.Location{
		Name:  "Ranhorn",
		Depth: 1,
		Entry: hornyMenuPos,
	}
)

var (
	Guild = &navi.Location{
		Name:   "Guild",
		Depth:  2,
		Entry:  guildPos,
		Parent: Ranhorn,
	}
	Shop = &navi.Location{
		Name:   "Shop",
		Depth:  2,
		Entry:  shopPos,
		Parent: Ranhorn,
	}
	OakInn = &navi.Location{
		Name:   "OakInn",
		Depth:  2,
		Entry:  oakPos,
		Parent: Ranhorn,
	}
)

var (
	Lab = &navi.Location{
		Name:   "ArcaneLab",
		Depth:  2,
		Entry:  labPos,
		Parent: DarkForest,
	}
	KT = &navi.Location{
		Name:   "KingsTower",
		Depth:  2,
		Entry:  kTPos,
		Parent: DarkForest,
	}
	Bounty = &navi.Location{
		Name:   "BountyBoard",
		Depth:  2,
		Entry:  bountyPos,
		Parent: DarkForest,
	}
	Arena = &navi.Location{
		Name:   "Arena",
		Depth:  2,
		Entry:  arenaPos,
		Parent: DarkForest,
	}
	Temporal = &navi.Location{
		Name:   "TemporalRift",
		Depth:  2,
		Entry:  temporalPos,
		Parent: DarkForest,
	}
)

var (
	GuildHunt = &navi.Location{
		Name:   "Guildhunting",
		Depth:  3,
		Entry:  ghuntPos,
		Parent: Guild,
	}
	Hellscape = &navi.Location{
		Name:   "Hellscape",
		Depth:  3,
		Entry:  hellscpPos,
		Parent: Guild,
	}
)

var (
	TwistedRealm = &navi.Location{
		Name:   "TwistedRealm",
		Depth:  4,
		Entry:  twistedPos,
		Parent: Hellscape,
	}
	ClownRealm = &navi.Location{
		Name:   "ClownRealm",
		Depth:  4,
		Entry:  clownPos,
		Parent: Hellscape,
	}
)

var (
	CRLeaderboardship = &navi.Location{
		Name:   "CRLeaders",
		Depth:  5,
		Entry:  CRLeadPos,
		Parent: ClownRealm,
	}
)

func UIMap() map[string]*navi.Location {
	flatmap[Main.Name] = Main
	flatmap[Campain.Name] = Campain
	flatmap[DarkForest.Name] = DarkForest
	flatmap[Ranhorn.Name] = Ranhorn
	flatmap[Guild.Name] = Guild
	flatmap[OakInn.Name] = OakInn
	flatmap[Shop.Name] = Shop
	flatmap[Arena.Name] = Arena
	flatmap[Bounty.Name] = Bounty
	flatmap[KT.Name] = KT
	flatmap[Lab.Name] = Lab
	flatmap[Temporal.Name] = Temporal
	flatmap[GuildHunt.Name] = GuildHunt
	flatmap[Hellscape.Name] = Hellscape
	flatmap[ClownRealm.Name] = ClownRealm
	flatmap[TwistedRealm.Name] = TwistedRealm
	return flatmap
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
