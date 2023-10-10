package afk

import (
	"fmt"
	"time"

	"mugowalker/backend/afk/activities"
	"mugowalker/backend/bot"
	"mugowalker/backend/cfg"
	"mugowalker/backend/image"

	// "mugowalker/backend/repository"
	"mugowalker/backend/settings"

	"github.com/go-color-term/go-color-term/coloring"
)

type Game struct {
	Name      string
	Active    bool
	Locations []any
	// User      *repository.User
	profile *settings.Pilot
}

type Task struct {
	Name  string
	Steps []*image.ScreenWord
}

var (
	f      = fmt.Sprintf
	finder = func(str string) func(*image.ScreenWord) bool {
		return func(s *image.ScreenWord) bool {
			return s.S == str
		}
	}
)

func (g *Game) String() string {
	return f("\n\tGameId:%v User:%v", g.Name, g.profile)
}

// New Game for a given User
func New(up *settings.Pilot) *Game {
	anylocs := activities.AllLocations()
	for _, l := range anylocs {
		if loc, ok := l.(activities.Location); ok {
			for _, kw := range loc.Keywords() {
				if kw == "%account" {
					loc.Kws = append(loc.Kws, up.Account)
				}
			}

		}
		anylocs = append(anylocs, l)
	}
	// user := &repository.User{Username: up.Account} //repository.GetUser(up.Account)

	return &Game{
		Name:      up.GameId,
		Locations: anylocs,
		Active:    true,
		// User:      user,
		profile: up,
	}
}

func Daily(b *Daywalker) {
	b.NotifyUI("[AFK]", "Daily Run")
quests:
	if err := b.FindTap(Quests.name, 0, 0); err != nil {
		b.Back()
		goto quests
	}
	t := b.Text().TesseractResult()
	board := availableQuests(t)
	b.NotifyUI("[BOARD]", f("%v", board))
	b.FindTap("100", 0, 300)
	for _, v := range board {
		switch v.S {
		case "gocollectloottimes":
			// campsScr := b.TapW(v)
			// cfg.Filter(campsScr, finder("begin"))
			b.NotifyUI("[Q]", coloring.Magenta("Collect loot N times"))
			b.TapSW(v)
			b.FindTap("fast", -200, -200)
			b.FindTap("collect", 0, 0)
			b.FindTap("fast", -200, -200)
			b.FindTap("collect", 0, 0)
			goto quests
		case "gobeginsolobountyquests":
			b.NotifyUI("[Q]", coloring.Magenta("Begin solo bounty quests"))
			b.TapSW(v)
			b.FindTap("all", 0, 0)
			b.FindTap("refresh", 222, 1700)
			b.FindTap("confirm", 0, 0)
			b.FindTap("refresh", 880, 1850)
			b.FindTap("all", 0, 0)
			b.FindTap("support", 222, 1700)
			b.FindTap("confirm", 0, 0)
			b.Back()
			goto quests
		case "gogiftfriendcompanionpointstime":
			//  GOOD
			b.NotifyUI("[Q]", coloring.Magenta("giftfriendcompanionpointstime"))
			b.TapSW(v)
			b.FindTap("receive", 0, 0)
			b.Back()
			goto quests
		case "goclaimfriendsgift":
			b.NotifyUI("[Q]", coloring.Magenta("Claim friends gift"))
			b.TapSW(v)
			b.FindTap("workshop", -130, -370)
			b.FindTap("workshop", -160, -370)
			b.FindTap("workshop", -110, -370)
			b.Back()
			goto quests
		case "gobeginbattlearenaheroes":
			b.NotifyUI("[Q]", coloring.Magenta("Begin battle arena heroes"))
			b.TapSW(v)
			// b.TapSW(v)#
			b.FindTap("glad", 0, -600)
			b.FindTap("confirm", 0, 0)
			b.FindTap("challenge", 0, 0)
			b.FindTap("refresh", 0, 670)
			b.FindTap("battle", 0, 0)
			time.Sleep(4 * time.Second)
			b.Back()
			b.FindTap("exit", 0, 0)
			b.Back()
			goto quests
		case "gobeginbattle":
			b.NotifyUI("[Q]", coloring.Magenta("Begin battle"))
			b.TapSW(v)
			b.FindTap("egin", 0, 0)
			b.FindTap("begin", 0, 0)
			b.FindTap("begin", 0, 0)
			b.FindTap("egin", 0, 0)
			b.Back()
			b.FindTap("exit", 0, 0)
			b.Back()
			goto quests
		case "gofastrewardsfunctiontime":
			b.NotifyUI("[Q]", coloring.Magenta("Fast rewards function N time"))
			b.TapSW(v)
			b.FindTap("collect", 0, 0)
			b.Back()
			goto quests
		case "gotakepartguildhunt":
			b.NotifyUI("[Q]", coloring.Magenta("Take part guild hunt"))
			b.TapSW(v)
			b.FindTap("battle", 0, 0)
			b.FindTap("sweep", 0, 0)
			b.Back()
			b.Back()
			goto quests
		case "gopurchaseitemsfromstoretime":
			b.NotifyUI("[Q]", coloring.Magenta("Purchase items from store N time"))
			b.TapSW(v)
			b.FindTap("barrack", 0, -1530)
			b.FindTap("cancel", 400, 0)
			b.Back()
			b.FindTap("barrack", 30, -2000)
			b.FindTap("confirm", 0, 0)
			b.FindTap("barrack", 0, -1530)
			b.FindTap("cancel", 400, 0)
			b.Back()
			b.FindTap("barrack", 30, -2000)
			b.FindTap("confirm", 0, 0)
			b.FindTap("barrack", 0, -1530)
			b.FindTap("cancel", 400, 0)
			b.Back()
			b.Back()
			goto quests
		}
	}
	goto quests
	b.NotifyUI("[Q]", coloring.Magenta("Take rewards and done!"))
	b.FindTap("100", 0, -40)
	b.Back()
	b.NotifyUI("[DAILY]", fmt.Sprintf("%v", board))
}

func availableQuests(ar []*image.ScreenWord) []*image.ScreenWord {
	// strGO := func(s *image.ScreenWord) bool { return s.S == "go" }

	wordsGO := cfg.Filter(ar, finder("go"))
	for _, b := range wordsGO {
		sameline := func(s *image.ScreenWord) bool { return s.LineNo == b.LineNo-1 }
		for _, v := range cfg.Filter(ar, sameline) {
			b.S += v.S
		}
	}
	return wordsGO
}

func Nightstalker(b *bot.Bot, user *settings.Pilot) *Daywalker {
	gm := New(user)
	return NewDaywalker(b, gm)
}
