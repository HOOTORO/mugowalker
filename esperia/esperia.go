package esperia

type Esperia interface {
	New()
	In() *Esperia
}

type Campain struct {
	depth int
	DarkForest
	Ranhorn
	screenpos TPoint
}

func (c *Campain) New() {
	panic("not implemented") // TODO: Implement
}

func (c *Campain) In() (int, int) {
	return c.screenpos.X, c.screenpos.Y
}

type DarkForest struct {
	Ranhorn
	KingsTower
	ArcaneLabyrinth
	ArenaOfHeroes
	TemporalRift
	BountyBoard
	screenpos TPoint
}

type KingsTower struct {
	TowerOfLight
	BrutalCitadel
	WorldTree
	ForsakenNecropolis
	screenpos TPoint
}

type TowerOfLight struct {
	pnt TPoint
}

type BrutalCitadel struct {
	pnt TPoint
}

type WorldTree struct {
	pnt TPoint
}

type ForsakenNecropolis struct {
	pnt TPoint
}

type ArcaneLabyrinth struct {
	pnt TPoint
}

type ArenaOfHeroes struct {
	pnt TPoint
}

type TemporalRift struct {
	pnt TPoint
}

type BountyBoard struct {
	pnt TPoint
}

type Ranhorn struct {
	OakInn
	Store
	Guild
	screenpos TPoint
}

func (c *Ranhorn) In() (int, int) {
	return c.screenpos.X, c.screenpos.Y
}

type OakInn struct {
	pnt TPoint
}

type Store struct {
	pnt TPoint
}

type Guild struct {
	// GHunting
	Hellscape
	screenpos TPoint
}

func (c *Guild) In() (int, int) {
	return c.screenpos.X, c.screenpos.Y
}

// type GHunting struct {
// 	Wrizz
// 	Soren
// 	screenpos TPoint
// }

// type Wrizz struct {
// pnt TPoint
// 	Fight
// }

// type Soren struct {
// pnt TPoint
// 	Fight
// }

type Hellscape struct {
	Cursed
	// Twisted
	screenpos TPoint
}

func (c *Hellscape) In() (int, int) {
	return c.screenpos.X, c.screenpos.Y
}

// type Twisted struct {
// 	TPoint
// 	Fight
// }

type Cursed struct {
	// Leaderboardship
	// Fight
	screenpos TPoint
}

func (c *Cursed) In() (int, int) {
	return c.screenpos.X, c.screenpos.Y
}

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

// type Fight struct {
// }

// type MultiFight struct {
// }

type TPoint struct {
	X int
	Y int
}
