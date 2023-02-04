package activities

import "fmt"

type Button interface {
	String() string
	Offset() (x, y int)
	Position() (x, y int)
}

type afkbtn struct {
	name         string
	x, y, xo, yo int
}

var (
	Quests          = afkbtn{name: "Quests"}
	Bag             = afkbtn{name: "Bag"}
	MailBtn         = afkbtn{name: "Mail"}
	Go              = afkbtn{name: "Go"}
	Collect         = afkbtn{name: "Collect"}
	Begin           = afkbtn{name: "Begin"}
	BeginB          = afkbtn{name: "Stage", yo: 739}
	BeginBoss       = afkbtn{name: "Begin", xo: 1}
	CampainBotPanel = afkbtn{name: "Campain"}
	ForrestBotPanel = afkbtn{name: "Forrest"}
	BattleBtn       = afkbtn{name: "Battle"}
	TryAgain        = afkbtn{name: "Again"}
	Next            = afkbtn{name: "Next"}
	Continue        = afkbtn{name: "Continue"}
	Challenge       = afkbtn{name: "Challenge"}
	King            = afkbtn{name: "Tower"}
	Wld             = afkbtn{name: "World"}
	Grvbrn          = afkbtn{name: "Forsaken"}
	Infrl           = afkbtn{name: "Infernal"}
	Mlr             = afkbtn{name: "Brutal"}
	Lght            = afkbtn{name: "Light"}
	Clstl           = afkbtn{name: "Celestial"}
	Any             = afkbtn{name: ""}
	Community       = afkbtn{name: "Community", yo: 80, xo: 40}
)

var (
	f = fmt.Sprintf
)

func (b afkbtn) String() string {
	return b.name
}
func (b afkbtn) Offset() (x, y int) {
	return b.xo, b.yo
}

func (b afkbtn) Position() (x, y int) {
	return b.x, b.y
}

func AllLocations() []any {
	return []any{
		// Bottom panel
		Ranhorn, Campain, Forrest,
		// Fight interface
		Prepare, Result, Win, Stats,
		// Banners
		RightBanner, Quest, Friends, Mail, PopoutExtra,
		// Campain
		Bossnode, FastRewards, Loot,
		// Dark Forrest
		Arena, KTentrance, KTinside,
		// DF -> Arena
		Soloarena, OpponentList,
		// DF -> Factional towers
		Graveborn, Wilder, Light, Nauler, Hypo, Celestial,
		// Ranhorn -> Guild
		guildgrounds, guildchest, wrizz, skipF,
		// Ranhorn -> Oak
		oak,
	}
}

type Location struct {
	ID  string
	Kws []string
	Hit int
}

func (l *Location) Id() string {
	return l.ID
}

func (l *Location) Keywords() []string {
	return l.Kws
}

func (l *Location) HitThreshold() int {
	return l.Hit
}

var baselocs = []*Location{Forrest, Ranhorn, Campain}

func isBaseLoc(s string) bool {
	for _, v := range baselocs {
		if v.ID == s {
			return true
		}
	}
	return false
}

var (
	Ranhorn = &Location{
		ID: "ranhorn",
		Kws: []string{"%account", "Guild", "Store", "Library",
			"Resonating", "Crystal", "Ascension", "Beast", "Grounds:", "Ascensigng", "Rickety",
			"WeeRickety", "Cart", "Trading", "BeastiGrounds", "Beast", "Noble", "Tavern",
			"Linrary", "BeastiGrounds)", "hejNoble", "Wall", "Legends", "Quests"},
		Hit: 5,
	}
	Campain = &Location{
		ID:  "campain",
		Kws: []string{"%account", "Campaign", "World", "Map", "Tales", "Fast", "Rewards", "World’Map", "Camp", "Quests"},
		Hit: 4,
	}
	Forrest = &Location{
		ID:  "forrest",
		Kws: []string{"%account", "Arena", "Peaks", "Time", "labyrnts", "KingissTower", "emporal", "Temporal", "Peaksiof", "Voyage", "Arcane", "Labyaintn", "Abyssal", "Expedition", "Bounty", "Board"},
		Hit: 5,
	}

	Prepare = &Location{
		ID:  "prepare",
		Kws: []string{"Forntions", "formations", "Formations", "Stage", "Stage:", "VS)", "Floor", "Battle", "BeginBattle", "must", "defeat", "teams", "advance", "VS"},
		Hit: 2,
	}
	Result = &Location{
		ID:  "result",
		Kws: []string{"Continue", "Raise", "increase", "strength", "Tier", "using", "methods", "below.", "next", "Fall", "Level", "Your", "Heroes", "Enhance", "Gear"},
		Hit: 10,
	}
	Win = &Location{
		ID:  "victory",
		Kws: []string{"Rewards", "TARY", "Tap", "Continue", "VIEETLARY", "TLARY", "VISELQIRY", "VISFLQIRY", "LQARY", "VIFF", "LARRY", "LORY", "Complete", "Rewards", "Next", "Stage"},
		Hit: 3,
	}
	Stats = &Location{
		ID:  "stat",
		Kws: []string{"%account", "Statistics", "Battle", "Hero", "Info", "Baltle"},
		Hit: 2,
	}
	RightBanner = &Location{
		ID:  "rBanOpen",
		Kws: []string{"%account", "Bag", "Mail", "Friends", "Solemn", "Vow", "Community"},
		Hit: 3,
	}
	Quest = &Location{
		ID:  "quiests",
		Kws: []string{"QUESTS", "Quests", "Refreshes", "Dailies", "Weeklies", "Campaign", "Completed", "Begin", "Battle", "Kings's", "Tower", "QUESTS", "Refreshes", "Level", "Hero", "Time)'", "Enhance", "Your", "Gear", "Timey", "Summon", "Hero", "Tavern"},
		Hit: 7,
	}

	Friends = &Location{
		ID:  "friends",
		Kws: []string{"Friends", "Garrisoned"},
		Hit: 2,
	}

	Bossnode = &Location{
		ID:  "bossnode",
		Kws: []string{"Enemy", "Formation", "Stage", "Stage:", "Completition", "Rewards", "ormavion"},
		Hit: 3,
	}
	Mail = &Location{
		ID:  "mail",
		Kws: []string{"collect", "all", "delete"},
		Hit: 3,
	}

	PopoutExtra = &Location{
		ID:  "popextra",
		Kws: []string{"extra", "Customize", "Bundle", "Disappears", "Purchase", "Anywhere", "1999", "Customizs", "Congratulations", "bundle"},
		Hit: 6,
	}
	FastRewards = &Location{
		ID:  "fastrewards",
		Kws: []string{"Collect", "Close", "Rewards", "Fast"},
		Hit: 3,
	}
	Loot = &Location{
		ID:  "loot",
		Kws: []string{"Tap", "the", "blank", "area", "claim", "AFK", "Rewards", "Timer", "Collect", "Close"},
		Hit: 4,
	}
	Arena = &Location{
		ID:  "arena",
		Kws: []string{"Season", "Ends", "CHALLENGER", "TOURNAMENT", "Rank", "Wins", "TREASURE", "LEGENDS:", "Division", "Starts", "championship", "starts", "gladiator", "Gladiator", "Coins", "Rating", "Required"},
		Hit: 6,
	}
	Soloarena = &Location{
		ID:  "soloarena",
		Kws: []string{"challendge", "Formation", "Record", "Arena", "Heroes", "Ladder", "Season", "Ranking", "Ends"},
		Hit: 3,
	}
	OpponentList = &Location{
		ID:  "opponent",
		Kws: []string{"Challenge", "Refresh", "Seregi", "Mpory"},
		Hit: 2,
	}
	KTentrance = &Location{
		ID:  "kt",
		Kws: []string{"Forsaken", "Necropolis", "Kings", "Tower", "Light", "Floors", "Stage", "Wed/Sat/Sun", "Mon/Fri/Sun", "Thu/Sat/Sun", "Mon/FriiSun"},
		Hit: 3,
	}
	KTinside = &Location{
		ID:  "king",
		Kws: []string{"Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge", "Kings's", "Kings"},
		Hit: 5,
	}
	Graveborn = &Location{
		ID:  "fn",
		Kws: []string{"Forsaken", "Necropolis", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit: 5,
	}
	Wilder = &Location{
		ID:  "wt",
		Kws: []string{"World", "Tree", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit: 5,
	}
	Light = &Location{
		ID:  "tol",
		Kws: []string{"Light", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit: 5,
	}
	Nauler = &Location{
		ID:  "bc",
		Kws: []string{"Brutal", "Citadel", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit: 5,
	}
	Hypo = &Location{
		ID:  "if",
		Kws: []string{"Infernal", "Fortress", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit: 5,
	}
	Celestial = &Location{
		ID:  "cs",
		Kws: []string{"Celestial", "Sanctum", "Fortress", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit: 5,
	}
	guildgrounds = &Location{
		ID:  "guildgrounds",
		Kws: []string{"Guild", "Hall", "Hellscape", "Grounds", "Hunting"},
		Hit: 2,
	}
	guildchest = &Location{
		ID:  "gichest",
		Kws: []string{"FORTUNE", "CHESTS", "Realm", "Fabled", "brave", "guildmate", "share", "with", "everyone"},
		Hit: 3,
	}
	wrizz = &Location{
		ID:  "wrizz",
		Kws: []string{"Wrizz", "Challenge"},
		Hit: 2,
	}
	skipF = &Location{
		ID:  "skipf",
		Kws: []string{"Quick", "Battle", "Sweep", "most", "recent"},
		Hit: 3,
	}
	oak = &Location{
		ID:  "oak",
		Kws: []string{"Workshop", "Tasks", "Smart", "Selections", "Manage"},
		Hit: 2,
	}
)

/*

	sample = &Location{
		ID:       "smplsId",
		Keywords: []string{""},
		Hit:      3,
	}

#######################################################################
##################### CAMPAIN CHILD LOCATIONS #########################
#######################################################################




#######################################################################
##################### DARK FOREST CHILD LOCATIONS #####################
#######################################################################




# kingstower:
, "name: kt
  description: Screen afther click on Kings Tower in Dark Forrest
  keywords:



, "name: king
  description: Inside tower, floors
  keywords:




#######################################################################
####################### RANHORNY TAB CHILD LOCATIONS ##################
#######################################################################


# shop:
, "name: shop
  hits: 2
  keywords:
    , "goods
    , "barracks


, "name: oak
  hits: 2
  keywords:




# <location>:
#   , "name:
#     hits:
#     keywords:
#         , "#word1
#         , "#word2
#     actions:
#         <act1>:
#
#                 x:
#                 "y":
#             properties:
#                 check: false
#                 delay: 1
#         <act2>:
#
#                 x:
#                 "y":
#             propeSprintf()
#                 check: false
#                 delay: 1
# back:
#     point:
#         x: "930"
#         "y": "2300"
#     properties:
#         check: false
#         delay: 1
#         repeat: 0

# entry:

)
*/
