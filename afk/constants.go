package afk

import "math"

type Const interface {
    String() string
    Id() uint
}
type UserField uint

type DailyQuest uint8

var QuestNames = []string{"loot", "fastrewards", "friends", "wrizz", "arena1x1", "oak", "QCamp", "QKT"}

func QString(k DailyQuest) []string {
	var result []string
	for i := 0; i < len(QuestNames); i++ {
		if k&(1<<uint(i)) != 0 {
			result = append(result, QuestNames[i])
		}
	}
	return result
}

func (dq DailyQuest) String() string {
    idx := math.Log2(float64(dq))
    return QuestNames[int(idx)]
}

func (dq DailyQuest) Id() uint8 {
	return uint8(dq)
}

// 000000
const (
	Loot DailyQuest = 1 << iota
	FastReward
	Friendship
	Wrizz
	Arena1x1
	Oak
	QCamp
	QKT
	Dailies = Loot | FastReward | Friendship | Wrizz | Arena1x1 | Oak | QCamp | QKT
)

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
			return Level(i)
		}
	}
	return 0
}

func (t Level) Id() uint {
	return uint(t)
}

// Popouts on locations
type Popout uint
func (p Popout) String() string {
    return popouts[p-1]
}
func (p Popout) Id() uint {
    return uint(p)
}
var popouts = [...]string{"skipf", "gichest","popextra"}
const (
    SKIPF Popout = iota +1
    GICHEST
	EXTRAPOPOUT

)

//General locations
const (
	ENTRY       = "campain"
	BATTLE   = "prepare"
	RESULT   = "result"
	WIN      = "victory"
	STAT     = "stat"
	BOSSTAGE = "bossnode"

	DARKFORREST = "forrest"
	RANHORNY = "ranhorn"

	HEROINFO = "heroinfo"
)

// Actions
const (
	DOPUSHCAMP   = "push"
	DOTOWERCLIMB = "climb"
	DOGIBOSSES   = "wrizz"
	DOOAK        = "oak"
)

const (
	GO2KT   = "zero2kt"
	GO2CAMP = "zero2camp"
)

const (
	kingone   = 700
	kingtwo   = 950
	facone    = 450
	factwo    = 660
	godone    = 350
	stages40  = 19
	chap1boss = 30
	chap2boss = 32
	chap3boss = 34
	chap4boss = 35
)
