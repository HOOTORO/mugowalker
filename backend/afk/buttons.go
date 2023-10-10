package afk

type afkButton struct {
	name         string
	x, y, xo, yo int
}

var (
	Quests          = afkButton{name: "quests"}
	Bag             = afkButton{name: "bag"}
	MailBtn         = afkButton{name: "mail"}
	Go              = afkButton{name: "go"}
	Collect         = afkButton{name: "collect"}
	Begin           = afkButton{name: "begin"}
	BeginB          = afkButton{name: "stage", yo: 739}
	BeginBoss       = afkButton{name: "begin", xo: 1}
	CampainBotPanel = afkButton{name: "campaign"}
	ForrestBotPanel = afkButton{name: "forrest"}
	BattleBtn       = afkButton{name: "battle"}
	TryAgain        = afkButton{name: "again"}
	Next            = afkButton{name: "next"}
	Continue        = afkButton{name: "continue"}
	Challenge       = afkButton{name: "challenge"}
	King            = afkButton{name: "tower"}
	Wld             = afkButton{name: "world"}
	GraveTower      = afkButton{name: "forsaken"}
	InfernalTower   = afkButton{name: "infernal"}
	Mlr             = afkButton{name: "brutal"}
	LightTower      = afkButton{name: "light"}
	CelestialTower  = afkButton{name: "celestial"}
	Any             = afkButton{name: ""}
	Community       = afkButton{name: "community", yo: 80, xo: 40}
)

func (b *afkButton) String() string {
	return b.name
}
func (b *afkButton) Offset() (x, y int) {
	return b.xo, b.yo
}

func (b *afkButton) Position() (x, y int) {
	return b.x, b.y
}
