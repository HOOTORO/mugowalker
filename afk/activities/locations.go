package activities

type Button string

const (
	Quests          Button = "Quests"
	Bag                    = "Bag"
	Mail                   = "Mail"
	Go                     = "Go"
	Collect                = "Collect"
	Begin                  = "Begin"
	CampainBotPanel        = "Campain"
	BattleBtn              = "Battle"
	TryAgain               = "Again"
	Next                   = "Next"
	Continue               = "Continue"
)

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
)

/*

, "name: friends
  hits: 2
  keywords:


, "name: mail
  grid: 4:3
  keywords:
    , "collect
    , "all
    , "delete
  transfers:
    , "locid: mailcollect
    , "locid: back

, "name: popextra
  grid: 1:18
  hits: 6
  keywords:
    , "extra
    , "Customize
    , "Bundle
    , "Disappears
    , "Purchase
    , "Anywhere
    , "1999
    , "Customizs
    , "Congratulations
    , "bundle

#######################################################################
##################### CAMPAIN CHILD LOCATIONS #########################
#######################################################################

, "name: bossnode
  hits: 3

  keywords:

, "name: fastrewards
  hits: 3
  keywords:
    , "Collect
    , "Close
    , "Rewards
    , "Fast

, "name: loot
  grid: 3:12
  hits: 4
  delay: 1
  keywords:
    , "Tap
    , "the
    , "blank
    , "area
    , "claim
    , "AFK
    , "Rewards
    , "Timer
    , "Collect
    , "Close

#######################################################################
##################### DARK FOREST CHILD LOCATIONS #####################
#######################################################################

# arena:
, "name: arena
  grid: 4:8
  hits: 6
  delay: 1
  keywords:
    , "Season
    , "Ends
    , "CHALLENGER
    , "TOURNAMENT
    , "Rank
    , "Wins
    , "TREASURE
    , ""LEGENDS:"
    , "Division
    , "Starts
    , "championship
    , "starts
    , "gladiator
    , "Gladiator
    , "Coins
    , "Rating
    , "Required
  actions:
    , "locid: oneonone
    , "locid: 3team
    , "locid: back

, "name: soloarena
  grid: 4:8
  hits: 3
  delay: 1
  keywords:
    , "challendge
    , "Formation
    , "Record
    , "Arena
    , "Heroes
    , "Ladder
    , "Season
    , "Ranking
    , "Ends
  actions:
    , "locid: oneonone
    , "locid: 3team
    , "locid: back

, "name: opponent
  grid: 4:8
  hits: 2
  delay: 1
  keywords:
    , "Challenge
    , "Refresh
    , "Seregi
    , "Mpory

# , "name: arenafight
#   grid: 4:8
#   hits: 1
#   delay: 1
#   keywords:
#     , "Battle

# kingstower:
, "name: kt
  description: Screen afther click on Kings Tower in Dark Forrest
  keywords:
    , ""%account"
    , "Forsaken
    , "Necropolis
    , "Kings
    , "Tower
    , "Light
    , "Floors
    , "Stage
    , "Wed/Sat/Sun
    , "Mon/Fri/Sun
    , "Thu/Sat/Sun
    , "Mon/FriiSun


, "name: king
  description: Inside tower, floors
  keywords:
    , "Floor
    , "Leaderboard
    , "Stage
    , "Info
    , "Cleared
    , "Challendge
    , "Kings's
    , "Kings

  actions:
    , "locid: back

, "name: fn
  grid: 3:9
  hits: 5
  delay: 1
  keywords:
    , "Forsaken
    , "Necropolis
    , "Floor
    , "Leaderboard
    , "Stage
    , "Info
    , "Cleared
    , "Challendge

, "name: wt
  grid: 3:9
  hits: 5
  delay: 1
  keywords:
    , "World
    , "Tree
    , "Floor
    , "Leaderboard
    , "Stage
    , "Info
    , "Cleared
    , "Challendge

, "name: tol
  grid: 3:9
  hits: 5
  delay: 1
  keywords:
    , "Light
    , "Floor
    , "Leaderboard
    , "Stage
    , "Info
    , "Cleared
    , "Challendge

, "name: bc
  hits: 5
  delay: 1
  keywords:
    , "Brutal
    , "Citadel
    , "Floor
    , "Leaderboard
    , "Stage
    , "Info
    , "Cleared
    , "Challendge

, "name: cs
  grid: 3:9
  hits: 5
  delay: 1
  keywords:
    , "Celestial
    , "Sanctum
    , "Floor
    , "Leaderboard
    , "Stage
    , "Info
    , "Cleared
    , "Challendge

, "name: if
  grid: 3:9
  hits: 5
  delay: 1
  keywords:
    , "Infernal
    , "Fortress
    , "Floor
    , "Leaderboard
    , "Stage
    , "Info
    , "Cleared
    , "Challendge

#######################################################################
####################### RANHORNY TAB CHILD LOCATIONS ##################
#######################################################################

# guild:
, "name: guildgrounds
  grid: 3:3
  hits: 2
  delay: 1
  keywords:
    , "Guild
    , "Hall
    , "Hellscape
    , "Grounds
    , "Hunting

, "name: gichest
  hits: 3
  keywords:
    , "FORTUNE
    , "CHESTS
    , "Realm
    , "Fabled
    , "brave
    , "guildmate
    , "share
    , "with
    , "everyone

# wrizz:
, "name: wrizz
  hits: 2
  keywords:
    , "Wrizz
    , "Challenge

, "name: skipf
  hits: 3
  keywords:
    , "Quick
    , "Battle
    , "Sweep
    , "most
    , "recent

# shop:
, "name: shop
  hits: 2
  keywords:
    , "goods
    , "barracks
  actions:
    , "locid: buydust
    , "locid: buypoegold
    , "locid: buypoedia
    , "locid: buycoresy
    , "locid: reset
    , "locid: back

, "name: oak
  hits: 2
  keywords:
    , "Workshop
    , "Tasks
    , "Smart
    , "Selections
    , "Manage
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
