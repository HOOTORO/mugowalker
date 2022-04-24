package esperia

type Esperia struct {
	Campain struct {
		entry TPoint
		depth int
	}
	DarkForest struct {
		entry      TPoint
		depth      int
		KingsTower struct {
			entry        TPoint
			depth        int
			TowerOfLight struct {
				entry TPoint
			}
			BrutalCitadel struct {
				entry TPoint
			}
			WorldTree struct {
				entry TPoint
			}
			ForsakenNecropolis struct {
				entry TPoint
			}
		}
		ArcaneLabyrinth struct {
			entry TPoint
		}
		ArenaOfHeroes struct {
			entry TPoint
		}
		TemporalRift struct {
			entry TPoint
		}
		BountyBoard struct {
			entry TPoint
		}
	}
	Ranhorn struct {
		entry  TPoint
		OakInn struct {
			pnt TPoint
		}
		Store struct {
			pnt TPoint
		}
		Guild struct {
			entry    TPoint
			GHunting struct {
				entry TPoint
				Wrizz struct {
					pnt TPoint
				}
				Soren struct {
					pnt TPoint
				}
			}
			Hellscape struct {
				entry  TPoint
				Cursed struct {
					entry           TPoint
					Leaderboardship struct {
						entry  TPoint
						Region struct {
							no int
						}
					}
				}
				Twisted struct {
					entry TPoint
				}
			}
		}
	}
}

func (from *Esperia) Path(to Esperia) []TPoint {
	return []TPoint{}
}

type Player struct {
	Name string
	Rank int
	Stat
}

type Stat struct {
	Page int
}

type TPoint struct {
	X int
	Y int
}

type Fight struct {
}

type MultiFight struct {
}
