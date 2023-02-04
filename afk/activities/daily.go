package activities

import (
	"fmt"
	"math"
	"strings"
	"time"
	"worker/cfg"
	"worker/ocr"

	"golang.org/x/exp/slices"
)

type Gamer interface {
	Id() uint
	Name() string
	Quests() uint
	SetQuests(uint)
}

var log = cfg.Logger()

type DailyQuest uint

var questNames = []string{"loot", "fastrewards", "friends", "wrizz", "arena1x1", "oak", "QCamp", "QKT", "Solo3Q", "LevelUp", "Enhance", "Summon"}
var BannedQuests = []DailyQuest{Solo3Q, LevelUp, Enhance, Summon}

func DailyRun(ns Nightstalker, g Gamer) {
	todo := Deserialize(g.Quests())
	_ = todo
	// go back until we get to the location with base footer menu and banners
	for !isBaseLoc(ns.Location()) {
		ns.Back()
	}
quests:
	if !ns.Press(Quests) {
		// TO DO: add OCR settings manipulation for better  recognition in case of not founding buttons
		time.Sleep(2 * time.Second)
		goto quests
	}
	bq := BoardsQuests(ns.OcResult())
	ns.Press(Collect)

	for _, q := range bq {
		if !slices.Contains(BannedQuests, q.Quest) {
			ns.Press(q.Btn)
			switch ns.Location() {
			case Campain.ID:
				if q.Quest.String() == "loot" {
					Begin.yo = -100
					if ns.Press(Begin) {
						if ns.Press(Collect) {
							ns.Back()
							markDone(g, q.Quest)
							goto quests
						}
					}
				}
			case KTentrance.ID:
				if ns.Press(King) {
					if ns.Press(Challenge) {
						ns.Press(BattleBtn)
					}
				}

			}
		}
	}

}

func ActiveDailies(u Gamer) []DailyQuest {
	return Deserialize(u.Quests())
}

func Deserialize(raw uint) []DailyQuest {
	var result []DailyQuest
	for i := 0; i < len(questNames); i++ {
		if d := DailyQuest(1 << i); raw&(1<<uint(i)) != 0 {
			result = append(result, d)
		}
	}
	return result
}

func markDone(u Gamer, q DailyQuest) {
	userQuests := DailyQuest(u.Quests())
	if !hasOneOf(q, userQuests) {
		u.SetQuests(q.Id())
	}
}

// 000000
const (
	LootQ DailyQuest = 1 << iota
	FastReward
	Friendship
	Wrizz
	Arena1x1
	Oak
	QCamp
	QKT
	Solo3Q
	LevelUp
	Enhance
	Summon
	Dailies = LootQ | FastReward | Friendship | Wrizz | Arena1x1 | Oak | QCamp | QKT | Solo3Q | LevelUp | Enhance | Summon
)

func (dq DailyQuest) String() string {
	idx := math.Log2(float64(dq))
	return questNames[int(idx)]
}

func (dq DailyQuest) Id() uint {
	return uint(dq)
}

func (dq DailyQuest) boardString() string {
	switch dq {
	case LootQ:
		return "Collect Loot Times"
	case FastReward:
		return "Fast Rewards Function Time"
	case Friendship:
		return "Gift Friend Companion Points"
	case QKT:
		return "Begin Battle King's Tower"
	case Wrizz:
		return "Take Part Guild Hunt"
	case Oak:
		return "Claim Friend's Gift"
	case Arena1x1:
		return "Battle Arena Heroes"
	case QCamp:
		return "Begin Battle"
	case Solo3Q:
		return "Begin Solo Bounty Quests"
	case LevelUp:
		return "Level Up Hero Time"
	case Enhance:
		return "Enhance Your Gear Time"
	case Summon:
		return "Summon Hero Tavern"
	default:
		return ""
	}
}

type BoardQuest struct {
	Quest DailyQuest
	Desc  []string
	Btn   Button
	X, Y  int
}

func (q BoardQuest) String() string {
	return fmt.Sprintf("\n|> Btn[%s]Que[%v]Pos[%vx%v] - Desc: %s", q.Btn, q.Quest, q.X, q.Y, q.Desc)
}

func BoardsQuests(or []ocr.AltoResult) []BoardQuest {
	var res []BoardQuest

	for _, str := range or {
		if str.Linechars == "Go" || strings.Contains(str.Linechars, "Completed") {
			qblock := &BoardQuest{}
			qblock.Desc = wordUpperBlock(str, or)
			qblock.Btn = afkbtn{name: str.Linechars}
			qblock.X = str.X
			qblock.Y = str.Y
			qblock.Quest = isBoardQuest(qblock.Desc)
			if qblock.Quest != 0 {
				res = append(res, *qblock)
			}
		}
	}
	return res
}

func isBoardQuest(s []string) DailyQuest {

	for _, q := range Deserialize(65535) {
		hits := cfg.Intersect(s, strings.Fields(q.boardString()))
		if len(hits) >= 3 {
			return q
		} else if len(s) == 2 && len(hits) == 2 {
			return q
		}
	}
	return 0
}

func wordUpperBlock(word ocr.AltoResult, or []ocr.AltoResult) []string {
	var res []string
	for _, str := range or {
		if diff := word.Y - str.Y; diff < 100 && diff > 40 {
			res = append(res, str.Linechars)
		}
	}
	return res

}

func set(p, flag DailyQuest) DailyQuest {
	return p | flag
}

func clear(p, flag DailyQuest) DailyQuest {
	return p &^ flag
}

func hasAll(p, flag DailyQuest) bool {
	return p&flag == flag
}

func hasOneOf(p, flag DailyQuest) bool {
	return p&flag != 0
}
