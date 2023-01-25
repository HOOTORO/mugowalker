package activities

type Button interface {
	String() string
	Offset() (x, y int)
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
	CampainBotPanel = afkbtn{name: "Campain"}
	BattleBtn       = afkbtn{name: "Battle"}
	TryAgain        = afkbtn{name: "Again"}
	Next            = afkbtn{name: "Next"}
	Continue        = afkbtn{name: "Continue"}
	Challenge       = afkbtn{name: "Challenge"}
	King            = afkbtn{name: "King's"}
)

func (b afkbtn) String() string {
	return b.name
}
func (b afkbtn) Offset() (x, y int) {
	return b.xo, b.yo
}

type Location struct {
	ID       string
	Keywords []string
	Hit      int
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
		Keywords: []string{"%account", "Guild", "Store", "Library",
			"Resonating", "Crystal", "Ascension", "Beast", "Grounds:", "Ascensigng", "Rickety",
			"WeeRickety", "Cart", "Trading", "BeastiGrounds", "Beast", "Noble", "Tavern",
			"Linrary", "BeastiGrounds)", "hejNoble", "Wall", "Legends", "Quests"},
		Hit: 5,
	}
	Campain = &Location{
		ID:       "campain",
		Keywords: []string{"%account", "Campaign", "World", "Map", "Tales", "Fast", "Rewards", "World’Map", "Camp", "Quests"},
		Hit:      4,
	}
	Forrest = &Location{
		ID:       "forrest",
		Keywords: []string{"%account", "Arena", "Peaks", "Time", "labyrnts", "KingissTower", "emporal", "Temporal", "Peaksiof", "Voyage", "Arcane", "Labyaintn", "Abyssal", "Expedition", "Bounty", "Board"},
		Hit:      5,
	}

	Prepare = &Location{
		ID:       "prepare",
		Keywords: []string{"Forntions", "formations", "Formations", "Stage", "Stage:", "VS)", "Floor", "Battle", "BeginBattle", "must", "defeat", "teams", "advance", "VS"},
		Hit:      2,
	}
	Result = &Location{
		ID: "result", Keywords: []string{"Continue", "Raise", "increase", "strength", "Tier", "using", "methods", "below.", "next", "Fall"},
		Hit: 5,
	}
	Win = &Location{
		ID:       "victory",
		Keywords: []string{"Rewards", "TARY", "Tap", "Continue", "VIEETLARY", "TLARY", "VISELQIRY", "VISFLQIRY", "LQARY", "VIFF", "LARRY", "LORY", "Complete", "Rewards", "Next", "Stage"},
		Hit:      3,
	}
	Stats = &Location{
		ID:       "stat",
		Keywords: []string{"%account", "Statistics", "Battle", "Hero", "Info", "Baltle"},
		Hit:      2,
	}
	RightBanner = &Location{
		ID:       "rBanOpen",
		Keywords: []string{"%account", "Bag", "Mail", "Friends", "Solemn", "Vow", "Community"},
		Hit:      3,
	}
	Quest = &Location{
		ID:       "quiests",
		Keywords: []string{"QUESTS", "Quests", "Refreshes", "Dailies", "Weeklies", "Campaign", "Completed", "Begin", "Battle", "Kings's", "Tower", "QUESTS", "Refreshes", "Level", "Hero", "Time)'", "Enhance", "Your", "Gear", "Timey", "Summon", "Hero", "Tavern"},
		Hit:      7,
	}

	Friends = &Location{
		ID:       "friends",
		Keywords: []string{"Friends", "Garrisoned"},
		Hit:      2,
	}

	Bossnode = &Location{
		ID:       "bossnode",
		Keywords: []string{"Enemy", "Formation", "Stage", "Stage:", "Completition", "Rewards", "ormavion"},
		Hit:      3,
	}
	Mail = &Location{
		ID:       "mail",
		Keywords: []string{"collect", "all", "delete"},
		Hit:      3,
	}

	PopoutExtra = &Location{
		ID:       "popextra",
		Keywords: []string{"extra", "Customize", "Bundle", "Disappears", "Purchase", "Anywhere", "1999", "Customizs", "Congratulations", "bundle"},
		Hit:      6,
	}
	FastRewards = &Location{
		ID:       "fastrewards",
		Keywords: []string{"Collect", "Close", "Rewards", "Fast"},
		Hit:      3,
	}
	Loot = &Location{
		ID:       "loot",
		Keywords: []string{"Tap", "the", "blank", "area", "claim", "AFK", "Rewards", "Timer", "Collect", "Close"},
		Hit:      4,
	}
	Arena = &Location{
		ID:       "arena",
		Keywords: []string{"Season", "Ends", "CHALLENGER", "TOURNAMENT", "Rank", "Wins", "TREASURE", "LEGENDS:", "Division", "Starts", "championship", "starts", "gladiator", "Gladiator", "Coins", "Rating", "Required"},
		Hit:      6,
	}
	Soloarena = &Location{
		ID:       "soloarena",
		Keywords: []string{"challendge", "Formation", "Record", "Arena", "Heroes", "Ladder", "Season", "Ranking", "Ends"},
		Hit:      3,
	}
	OpponentList = &Location{
		ID:       "opponent",
		Keywords: []string{"Challenge", "Refresh", "Seregi", "Mpory"},
		Hit:      2,
	}
	KTentrance = &Location{
		ID:       "kt",
		Keywords: []string{"Forsaken", "Necropolis", "Kings", "Tower", "Light", "Floors", "Stage", "Wed/Sat/Sun", "Mon/Fri/Sun", "Thu/Sat/Sun", "Mon/FriiSun"},
		Hit:      3,
	}
	KTinside = &Location{
		ID:       "king",
		Keywords: []string{"Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge", "Kings's", "Kings"},
		Hit:      5,
	}
	Graveborn = &Location{
		ID:       "fn",
		Keywords: []string{"Forsaken", "Necropolis", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit:      5,
	}
	Wilder = &Location{
		ID:       "wt",
		Keywords: []string{"World", "Tree", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit:      5,
	}
	Light = &Location{
		ID:       "tol",
		Keywords: []string{"Light", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit:      5,
	}
	Nauler = &Location{
		ID:       "bc",
		Keywords: []string{"Brutal", "Citadel", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit:      5,
	}
	Hypo = &Location{
		ID:       "if",
		Keywords: []string{"Infernal", "Fortress", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit:      5,
	}
	Celestial = &Location{
		ID:       "cs",
		Keywords: []string{"Celestial", "Sanctum", "Fortress", "Floor", "Leaderboard", "Stage", "Info", "Cleared", "Challendge"},
		Hit:      5,
	}
	guildgrounds = &Location{
		ID:       "guildgrounds",
		Keywords: []string{"Guild", "Hall", "Hellscape", "Grounds", "Hunting"},
		Hit:      2,
	}
	guildchest = &Location{
		ID:       "gichest",
		Keywords: []string{"FORTUNE", "CHESTS", "Realm", "Fabled", "brave", "guildmate", "share", "with", "everyone"},
		Hit:      3,
	}
	wrizz = &Location{
		ID:       "wrizz",
		Keywords: []string{"Wrizz", "Challenge"},
		Hit:      2,
	}
	skipF = &Location{
		ID:       "skipf",
		Keywords: []string{"Quick", "Battle", "Sweep", "most", "recent"},
		Hit:      3,
	}
	oak = &Location{
		ID:       "oak",
		Keywords: []string{"Workshop", "Tasks", "Smart", "Selections", "Manage"},
		Hit:      2,
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
#             properties:
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
