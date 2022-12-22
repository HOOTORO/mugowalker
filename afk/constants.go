package afk

type UserField int

type DailyQuest uint8

func (dq DailyQuest) String() string {
    if dq < Dailies{
        return [...]string{"Loot", "FastRewards", "Friendship", "Wrizz", "Arena", "OakInn"}[dq-1]
    }    else {
        return ""
    }

}
func  (dq DailyQuest) Indx() uint8{
    return uint8(dq)
}


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
	WIN      = "victory"
	LOSE     = "losecampain"
	STAT     = "stat"
	BOSSTAGE = "bossnode"
)

const (
	CAMPBegin = "campBegin"

	RBAN = "rBanOpen"
	LBAN
	// actions
	QUEST = "quest"
	BAG   = "bag"
)

const (
    CAMPAIN = "campain"
	DARKFORREST = "forrest"
	KTower      = "kt"
	PVP         = "arena"
	SOLO        = "solo"
)

const (
	RANHORNY = "ranhorn"
	HORNSHOP = "shop"
	GI       = "guild"
)

const (
	HEROINFO = "heroinfo"
)

const (
	INFERNAL = "infortress"
)

// Actions
const (
	FIGHT = "Battle"
	TOWERCLIMB = "climb"

	DailyLoOt   = "loot"
	FastRewards = "fastrewards"
	MAIL        = "mail"
	FRIENDS     = "friends"
	GIBOSSES    = "wrizz"
	OAK         = "oak"
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
