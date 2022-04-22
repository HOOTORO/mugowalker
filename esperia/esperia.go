package esperia

type Esperia interface {
	Walk() (TPoint, *Campain)
	NewCAMPAIN() *Campain
}

type Campain struct {
	entry TPoint
	depth int
	DarkForest
	Ranhorn
}

func NewCAMPAIN() *Campain {
	return &Campain{entry: CampainEntry, depth: 0}
}

func (c *Campain) Walk() (TPoint, interface{}) {
	return c.entry, c
}

// func (c *Campain) Begin()

// func (c *Campain) DarkForest() (TPoint, *DarkForest) {
// 	return DarkFP, &DarkForest{entry: DarkFP, depth: 0}
// }

type DarkForest struct {
	entry TPoint
	depth int
	// state interface{}
	// Ranhorn
	// KingsTower
	// ArcaneLabyrinth
	// ArenaOfHeroes
	// TemporalRift
	// BountyBoard
}

func (c *DarkForest) Walk() (TPoint, interface{}) {
	return c.entry, c
}
func (df *DarkForest) KingsTower() (TPoint, *KingsTower) {
	return KingsTP, &KingsTower{entry: KingsTP, depth: 1}
}

type KingsTower struct {
	entry TPoint
	depth int
	// TowerOfLight
	// BrutalCitadel
	// WorldTree
	// ForsakenNecropolis
}

// func (kt *KingsTower) KT() (TPoint, *KingsTower) {
// 	return KingsTP, kt
// }

// type TowerOfLight struct {
// 	entry TPoint
// }

// type BrutalCitadel struct {
// 	entry TPoint
// }

// type WorldTree struct {
// 	entry TPoint
// }

// type ForsakenNecropolis struct {
// 	entry TPoint
// }

// type ArcaneLabyrinth struct {
// 	entry TPoint
// }

// type ArenaOfHeroes struct {
// 	entry TPoint
// }

// type TemporalRift struct {
// 	entry TPoint
// }

// type BountyBoard struct {
// 	entry TPoint
// }

type Ranhorn struct {
	entry TPoint
	// OakInn
	// Store
	// Guild
}

// type OakInn struct {
// 	pnt TPoint
// }

// type Store struct {
// 	pnt TPoint
// }

// type Guild struct {
// 	entry TPoint
// 	GHunting
// 	Hellscape
// }

// type GHunting struct {
// 	entry TPoint
// 	Wrizz
// 	Soren
// }

// type Wrizz struct {
// 	pnt TPoint
// 	Fight
// }

// type Soren struct {
// 	pnt TPoint
// 	Fight
// }

// type Hellscape struct {
// 	entry TPoint
// 	Cursed
// 	Twisted
// 	screenpos TPoint
// }

// type Twisted struct {
// 	TPoint
// 	Fight
// }

// type Cursed struct {
// 	entry TPoint
// 	Leaderboardship
// 	Fight
// 	screenpos TPoint
// }

// type Leaderboardship struct {
// 	screenpos TPoint
// }

// type Region struct{}

// type Player struct {
// 	Name string
// 	Rank int
// 	Stat
// }

// type Stat struct {
// 	Page      int
// 	screenpos TPoint
// }

type Fight struct {
}

type MultiFight struct {
}

type TPoint struct {
	X int
	Y int
}
