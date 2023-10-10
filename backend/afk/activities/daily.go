package activities

import (
	"math"
)

type Gamer interface {
	Id() uint
	Name() string
	Quests() uint
	SetQuests(uint)
}

type DailyQuest uint

var (
	quests = []string{
		"loot", "fastrewards", "friends",
		"wrizz", "arena1x1",
		"oak", "qcamp", "qkt",
		"solo3q", "levelup",
		"enhance", "summon",
	}
	descs = []string{
		"collect loot times", "fast rewards function time", "gift friend companion points",
		"take part guild hunt", "battle arena heroes",
		"claim friend's gift", "begin battle", "begin battle king's tower",
		"begin solo bounty quests", "level up hero time",
		"enhance your gear time", "summon hero tavern",
	}
)

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

func (dq DailyQuest) boardString() string {

	idx := math.Log2(float64(dq))
	return descs[int(idx)]

}
func (dq DailyQuest) String() string {
	idx := math.Log2(float64(dq))
	return quests[int(idx)]
}

func (dq DailyQuest) Id() uint {
	return uint(dq)
}

func ActiveDailies(u Gamer) []DailyQuest {
	return Deserialize(u.Quests())
}

func Deserialize(raw uint) []DailyQuest {
	var result []DailyQuest
	for i := 0; i < len(quests); i++ {
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
