package afk

type UserField int

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
	BACK  = "back"
	CLOSE = "close"
	ENTRY = "entry"
	FIGHT = "fight"
)

const (
	PUSHc      = "pushcampain"
	BATTLE     = "battlescreen"
	LOSE     = "losecampain"
	BOSSTAGE = "campainBoss"
    BOSSBegin = "bossfight"
	CAMPWIN  = "campvictory"
	BATTLESTAT = "battlestat"
	RETRY      = "retry"
	NEXTSTAGE  = "nxt"
)

const (
	CAMPBegin = "campBegin"
	AFKCHEST  = "afkchest"
	FR       = "fastrewards"
	USEFR    = "usefr"
	RBAN     = "rBanOpen"
	LBAN
	// actions
	QUEST       = "quest"
	BAG         = "bag"
	MAIL        = "mail"
	MAILCOLLECT = "mailcollect"
	FRIENDS     = "friends"
	LIKESBTN    = "sendrecive"
)

const (
	DARKFORREST = "forrest"
	KTower          = "kingstower"
	PVP         = "arena"
	SOLO        = "solo"
)

const (
	RANHORNY = "ranhorn"
	HORNSHOP = "shop"
	GI       = "guild"
	GIBOSSES = "bosses"
	OAK      = "oak"
)

const (
	HEROINFO = "heroinfo"
)

const (
	INFERNAL = "infortress"
)
// Actions


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
