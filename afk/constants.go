package afk


type UserField int

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
func  (dq DailyQuest) Indx() uint8{
    return uint8(dq)
}

//000000
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

type Tower int



var strs = [...]string{"", "name", "account_id", "vip", "chapter", "stage", "diamonds", "gold"}

const (
	USERNAME UserField = iota + 1
	ACCOUNTID
	VIP
	CHAPTER
	STAGE
	DIAMONDS
	GOLD
)

const (
    Kings Tower = iota +1
    Celestial
    Infernal
    Light
    Mauler
    Wilder
    Graveborn
)
func (uf UserField) String() string {
	return strs[uf]
}

func BottomPanel() []string {
	return []string{RANHORNY, DARKFORREST, CAMPBegin}
}

/*
	Locations and correspondinf actions config names
*/
//General
const (
	ENTRY = "campain"
)

const (
	BATTLE   = "prepare"
	RESULT = "result"
	WIN    = "victory"
	STAT   = "stat"
	BOSSTAGE = "bossnode"
)

const (
	CAMPBegin = "campBegin"
)

const (

	DARKFORREST = "forrest"
	KTower      = "kt"
    TowerInside = "tower"

)

const (
	RANHORNY = "ranhorn"

)

const (
	HEROINFO = "heroinfo"
)


// Actions
const (
    DOPUSHCAMP    = "pushcamp"
	DOTOWERCLIMB = "climb"
	DOGIBOSSES    = "wrizz"
	DOOAK         = "oak"
)

const (
    GO2KT = "zero2kt"
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
