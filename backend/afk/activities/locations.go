package activities

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
		Graveborn, Wilder, Light, Mauler, Hypo, Celestial,
		// Ranhorn -> Guild
		guildgrounds, guildchest, WrizzLoc, skipF,
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

var baseLocations = []*Location{Forrest, Ranhorn, Campain}

func isBaseLoc(s string) bool {
	for _, v := range baseLocations {
		if v.ID == s {
			return true
		}
	}
	return false
}

var (
	Ranhorn = &Location{
		ID: "ranhorn",
		Kws: []string{"%account", "guild", "store", "library",
			"resonating", "crystal", "ascension", "beast", "grounds:", "ascensigng", "rickety",
			"weerickety", "cart", "trading", "beastigrounds", "beast", "noble", "tavern",
			"linrary", "beastigrounds)", "hejnoble", "wall", "legends", "quests"},
		Hit: 5,
	}
	Campain = &Location{
		ID:  "campain",
		Kws: []string{"%account", "campaign", "world", "map", "tales", "fast", "rewards", "world’map", "camp", "quests"},
		Hit: 4,
	}
	Forrest = &Location{
		ID:  "forrest",
		Kws: []string{"%account", "arena", "peaks", "time", "labyrnts", "kingissTower", "emporal", "temporal", "peaksiof", "voyage", "arcane", "labyaintn", "abyssal", "expedition", "bounty", "board"},
		Hit: 5,
	}

	Prepare = &Location{
		ID:  "prepare",
		Kws: []string{"forntions", "formations", "formations", "stage", "stage:", "vs)", "floor", "battle", "beginbattle", "must", "defeat", "teams", "advance", "vs"},
		Hit: 2,
	}
	Result = &Location{
		ID:  "result",
		Kws: []string{"continue", "raise", "increase", "strength", "tier", "using", "methods", "below.", "next", "fall", "level", "your", "heroes", "enhance", "gear"},
		Hit: 10,
	}
	Win = &Location{
		ID:  "victory",
		Kws: []string{"rewards", "tary", "tap", "continue", "vieetlary", "tlary", "viselqiry", "visflqiry", "lqary", "viff", "larry", "lory", "complete", "rewards", "next", "stage"},
		Hit: 3,
	}
	Stats = &Location{
		ID:  "stat",
		Kws: []string{"%account", "statistics", "battle", "hero", "info", "baftle"},
		Hit: 2,
	}
	RightBanner = &Location{
		ID:  "rBanOpen",
		Kws: []string{"%account", "bag", "mail", "friends", "solemn", "vow", "community"},
		Hit: 3,
	}
	Quest = &Location{
		ID:  "quests",
		Kws: []string{"quests", "quests", "refreshes", "dailies", "weeklies", "campaign", "completed", "begin", "battle", "kings's", "tower", "quests", "refreshes", "level", "hero", "time)'", "enhance", "your", "gear", "timey", "summon", "hero", "tavern"},
		Hit: 7,
	}

	Friends = &Location{
		ID:  "friends",
		Kws: []string{"friends", "garrisoned"},
		Hit: 2,
	}

	Bossnode = &Location{
		ID:  "bossnode",
		Kws: []string{"enemy", "formation", "stage", "stage:", "completition", "rewards", "ormavion"},
		Hit: 3,
	}
	Mail = &Location{
		ID:  "mail",
		Kws: []string{"collect", "all", "delete"},
		Hit: 3,
	}

	PopoutExtra = &Location{
		ID:  "popextra",
		Kws: []string{"extra", "customize", "bundle", "Disappears", "purchase", "anywhere", "1999", "customizs", "congratulations", "bundle"},
		Hit: 6,
	}
	FastRewards = &Location{
		ID:  "fastrewards",
		Kws: []string{"collect", "close", "rewards", "fast"},
		Hit: 3,
	}
	Loot = &Location{
		ID:  "loot",
		Kws: []string{"tap", "the", "blank", "area", "claim", "aFK", "rewards", "timer", "collect", "close"},
		Hit: 4,
	}
	Arena = &Location{
		ID:  "arena",
		Kws: []string{"season", "ends", "challenger", "tournament", "rank", "wins", "treasure", "legends:", "division", "starts", "championship", "starts", "gladiator", "gladiator", "coins", "rating", "required"},
		Hit: 6,
	}
	Soloarena = &Location{
		ID:  "soloarena",
		Kws: []string{"challendge", "formation", "record", "arena", "Heroes", "ladder", "season", "ranking", "ends"},
		Hit: 3,
	}
	OpponentList = &Location{
		ID:  "opponent",
		Kws: []string{"challenge", "refresh", "seregi", "mpory"},
		Hit: 2,
	}
	KTentrance = &Location{
		ID:  "kt",
		Kws: []string{"forsaken", "necropolis", "kings", "tower", "light", "floors", "stage", "wed/sat/sun", "mon/fri/sun", "thu/sat/sun", "mon/friisun"},
		Hit: 3,
	}
	KTinside = &Location{
		ID:  "king",
		Kws: []string{"floor", "leaderboard", "stage", "Info", "cleared", "challendge", "kings's", "kings"},
		Hit: 5,
	}
	Graveborn = &Location{
		ID:  "fn",
		Kws: []string{"forsaken", "necropolis", "floor", "leaderboard", "stage", "info", "cleared", "challendge"},
		Hit: 5,
	}
	Wilder = &Location{
		ID:  "wt",
		Kws: []string{"world", "tree", "floor", "leaderboard", "stage", "info", "cleared", "challendge"},
		Hit: 5,
	}
	Light = &Location{
		ID:  "tol",
		Kws: []string{"light", "floor", "leaderboard", "stage", "info", "cleared", "challendge"},
		Hit: 5,
	}
	Mauler = &Location{
		ID:  "bc",
		Kws: []string{"brutal", "citadel", "floor", "leaderboard", "stage", "info", "cleared", "challendge"},
		Hit: 5,
	}
	Hypo = &Location{
		ID:  "if",
		Kws: []string{"Infernal", "fortress", "floor", "leaderboard", "stage", "info", "cleared", "challendge"},
		Hit: 5,
	}
	Celestial = &Location{
		ID:  "cs",
		Kws: []string{"celestial", "sanctum", "fortress", "floor", "leaderboard", "stage", "info", "cleared", "challendge"},
		Hit: 5,
	}
	guildgrounds = &Location{
		ID:  "guildgrounds",
		Kws: []string{"guild", "hall", "hellscape", "grounds", "hunting"},
		Hit: 2,
	}
	guildchest = &Location{
		ID:  "gichest",
		Kws: []string{"fortune", "chests", "realm", "fabled", "brave", "guildmate", "share", "with", "everyone"},
		Hit: 3,
	}
	WrizzLoc = &Location{
		ID:  "wrizz",
		Kws: []string{"wrizz", "challenge"},
		Hit: 2,
	}
	skipF = &Location{
		ID:  "skipf",
		Kws: []string{"quick", "battle", "sweep", "most", "recent"},
		Hit: 3,
	}
	oak = &Location{
		ID:  "oak",
		Kws: []string{"workshop", "tasks", "smart", "selections", "manage"},
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
