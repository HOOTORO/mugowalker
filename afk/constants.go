package afk

type Location interface {
	String() string
	Id() uint
}

type UserField uint

var strs = [...]string{"name", "account_id", "vip", "chapter", "stage", "diamonds", "gold"}

const (
	USERNAME UserField = iota + 1
	ACCOUNTID
	VIP
	DIAMONDS
	GOLD
)

func (uf UserField) String() string {
	return strs[uf-1]
}

type ArenaLocation uint

func (al ArenaLocation) String() string {
	return alocs[al-1]
}

func (al ArenaLocation) Id() uint {
	return uint(al)
}

func ArenaLoc(s string) ArenaLocation {
	for i, v := range alocs {
		if v == s {
			return ArenaLocation(i + 1)
		}
	}
	return 0
}

var alocs = [...]string{
	"campain",
	"forrest", "ranhorn", "prepare", "result", "victory",
	"stat", "heroinfo", "bossnode", "lBanOpen", "rBanOpen",
	"friends", "mail", "popextra", "fastrewards", "loot",
	"arena", "soloarena", "opponent", "king",
	"kt", "fn", "wt", "tol", "bc", "cs", "if",
	"guildgrounds", "gichest", "wrizz", "skipf", "shop", "oak", "quests",
}

const (
	Campain ArenaLocation = iota + 1
	DarkForrest
	RANHORNY
	BATTLE
	RESULT
	WIN
	STAT
	HEROINFO
	BOSSTAGE
	LBAN
	RBAN
	FRIENDS
	MAIL
	EDEAL
	FASTR
	AFKCHEST
	ARENA
	SOLO
	OPPO
	KING
	KT
	FN
	WT
	TOL
	BC
	CS
	IF
	GGROUNDS
	GUILDCHEST
	WRIZZ
	SkipF
	SHOP
	OAK
	QUESTS
)

type Level uint

var towers = [...]string{"ch", "stg", "kt", "tol", "bc", "wt", "fn", "cs", "if"}

const (
	Chapter Level = iota + 1
	Stage
	Kings
	Light
	Mauler
	Wilder
	Graveborn
	Celestial
	Infernal
)

func (t Level) String() string {
	return towers[t-1]
}

func LocLvl(s string) Level {
	for i, v := range towers {
		if v == s {
			return Level(i + 1)
		}
	}
	return 0
}

func (t Level) Id() uint {
	return uint(t)
}

// Popout Popouts on locations
type Popout uint

func (p Popout) String() string {
	return popouts[p-1]
}

func (p Popout) Id() uint {
	return uint(p + 1)
}

var popouts = [...]string{"skipf", "gichest", "popextra"}

const (
	SKIPF Popout = iota + 1
	GICHEST
	EXTRAPOPOUT
)

// General locations
type Reaction uint

var rcts = [...]string{"push", "climb"}

func (r Reaction) String() string {
	return rcts[r-1]
}

func (r Reaction) Id() uint {
	return uint(r)
}

// Actions
const (
	DOPUSHCAMP Reaction = iota + 1
	DOTOWERCLIMB

// DOGIBOSSES   = "wrizz"
// DOOAK        = "oak"
)

const (
	GO2KT   = "zero2kt"
	GO2CAMP = "zero2camp"
)

//const (
//	kingone   = 700
//	kingtwo   = 950
//	facone    = 450
//	factwo    = 660
//	godone    = 350
//	stages40  = 19
//	chap1boss = 30
//	chap2boss = 32
//	chap3boss = 34
//	chap4boss = 35
//)

type Action uint

var bas = [...]string{"updProgress", "custom", "deactivate", "gshot", "repeatx"}

const (
	UpdProgress Action = iota + 1
	Custom
	Deactivate
	Gshot

	RepeatX
)

func (ba Action) String() string {
	return bas[ba-1]
}

func IsAction(s string) (Action, bool) {
	for i, act := range bas {
		if act == s {
			return Action(i + 1), true
		}
	}
	return 0, false
}
