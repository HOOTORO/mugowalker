package afk

import (
	"errors"
	"fmt"
	"strings"

	"mugowalker/backend/afk/activities"
	"mugowalker/backend/bot"
	"mugowalker/backend/image"
)

var (
	ErrUnknownButton = errors.New("not found button named")
)

type Daywalker struct {
	*bot.Bot
	*Game
	ActiveTask string
	currentOcr *image.ImageProfile
}

func NewDaywalker(b *bot.Bot, g *Game) *Daywalker {
	return &Daywalker{
		Bot:  b,
		Game: g,
	}
}

func (dw *Daywalker) String() string {
	return fmt.Sprintf("\n\tDevice: %+v%v", dw.Device, dw.Game)
}

// ///////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////
func (dw *Daywalker) Location() string {
	dw.currentOcr = dw.Text()
	return bot.GuessLocation(dw.currentOcr, dw.Locations)

}

func (dw *Daywalker) TapOrBack(word string) bool {
	if err := dw.FindTap(word, 0, 0); err != nil {
		dw.Back()
		return false
	} else {
		return true
	}
}

func LookForButton(or []*image.ScreenWord, b activities.Button) (x, y int, e error) {
	fmt.Printf("\nLooking for BTN: %v", b.String())
	if b.String() == "" {
		return 500, 100, nil
	}
	for _, r := range or {
		if strings.Contains(r.S, b.String()) {
			return r.X, r.Y, nil
		}
	}
	return 0, 0, errors.New("btn not found here")
}
