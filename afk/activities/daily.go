package activities

import (
	"math"
	"strings"
)

type Gamer interface {
	Id() uint
	Name() string
	Quests() uint8
	SetQuests(uint8)
}
type DailyQuest uint8

type Routes map[DailyQuest]map[int]string

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

var QuestNames = []string{"loot", "fastrewards", "friends", "wrizz", "arena1x1", "oak", "QCamp", "QKT"}
var allowedQuests = []DailyQuest{Loot, FastReward, Friendship, Wrizz, Arena1x1, Oak}

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

func (dq DailyQuest) BoardString() string {
	switch dq {
	case Loot:
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
	default:
		return ""
	}
}

func Set(p, flag DailyQuest) DailyQuest {
	return p | flag
}

func Clear(p, flag DailyQuest) DailyQuest {
	return p &^ flag
}

func HasAll(p, flag DailyQuest) bool {
	return p&flag == flag
}

func HasOneOf(p, flag DailyQuest) bool {
	return p&flag != 0
}

func ActiveDailies(u Gamer) []DailyQuest {
	var res []DailyQuest
	for i := 0; i < len(QuestNames); i++ {
		if u.Quests()&(1<<uint(i)) == 0 {
			res = append(res, DailyQuest(1<<uint(i)))
		}
	}
	return res
}

func MarkDone(u Gamer, q DailyQuest) {
	userQuests := DailyQuest(u.Quests())
	if !HasOneOf(q, userQuests) {
		u.SetQuests(q.Id())
		// Fnotify("|>",red("--> DAILY <-- \nCurrent: [%08b] \nOverall: [%08b]", quesst, g.ActiveDailies()))
	}
}

func IsBoardQuest(s string) DailyQuest {
	if s == QCamp.BoardString() {
		return QCamp
	}
	words := strings.Fields(s)
	for _, sd := range allowedQuests {
		scoring := 0
		for _, w := range words {
			if strings.Contains(sd.BoardString(), w) {
				scoring++
				if scoring >= 3 {
					return sd
				}
			}
		}

	}
	return 0

}

func Route(q DailyQuest) map[int]string {
	r := make(Routes)

	r[Loot] = make(map[int]string)
	r[Loot][1] = "3:16"
	r[Loot][2] = "Collect"
	r[Loot][3] = "Fast Rewards"
	r[Loot][4] = "Collect"
	r[Loot][5] = "Back"
	r[Loot][6] = "Collect"
	r[Loot][7] = "Back"

	r[FastReward] = make(map[int]string)
	r[FastReward][1] = "3:16"
	r[FastReward][2] = "Collect"
	r[FastReward][3] = "Fast Rewards"
	r[FastReward][4] = "Collect"
	r[FastReward][5] = "Back"
	r[FastReward][6] = "Collect"
	r[FastReward][7] = "Back"
	r[Friendship] = make(map[int]string)
	r[Friendship][1] = "Send Receive"
	r[Friendship][2] = "Back"
	r[Wrizz] = make(map[int]string)
	r[Wrizz][1] = "Quick Battle"
	r[Wrizz][2] = "Sweep"
	r[Wrizz][3] = "Back"
	r[Wrizz][4] = "Back"
	r[Oak] = make(map[int]string)
	r[Oak][1] = "GoOld"
	r[Arena1x1] = make(map[int]string)
	r[Arena1x1][1] = "GoOld"

	return r[q]

}
