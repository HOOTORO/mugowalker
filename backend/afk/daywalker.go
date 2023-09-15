package afk

import (
	"errors"
	"strings"
	"time"

	"mugowalker/backend/afk/activities"
	"mugowalker/backend/afk/repository"
	"mugowalker/backend/bot"
	c "mugowalker/backend/cfg"
	"mugowalker/backend/ocr"

	"golang.org/x/exp/slices"
)

const (
	alto = "ALTO"
)

var (
	ErrUnknownButton = errors.New("Not found button named:")
)

type Daywalker struct {
	*bot.BasicBot
	*Game

	cnt            uint8
	ActiveTask     string
	Reactive       bool
	currentOcr     *ocr.ImageProfile
	fprefix        string
	lastscreenshot string
	maxocrtry      int
	knownBtns      map[activities.Button]sessionBtn
	knwnBtns       []*repository.Button
}

type sessionBtn struct {
	x, y int
}

func (s sessionBtn) String() string {
	return "ssnBtn"
}
func (s sessionBtn) Offset() (x int, y int) {
	return 0, 0
}

func (s sessionBtn) Position() (x int, y int) {
	return s.x, s.y
}

func NewArenaBot(b *bot.BasicBot, g *Game) *Daywalker {
	btns := repository.GetButtons(b.Resolution.X, b.Resolution.Y)
	return &Daywalker{
		BasicBot: b,
		Game:     g,
		fprefix:  time.Now().Format("2006_01"),
		cnt:      0, maxocrtry: 2,
		knownBtns: make(map[activities.Button]sessionBtn, 0),
		knwnBtns:  btns,
	}
}

var outFn func(string, string)

func init() {
	log = c.Logger()

}

func (dw *Daywalker) String() string {
	return c.F("\nBot status:\n   Game: %v\n ActiveTask: %v\nDevice: %+v", dw.Game, dw.ActiveTask, dw.Device)
}

// ///////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////
func (dw *Daywalker) Location() string {
	dw.currentOcr = dw.ScanText()
	return bot.GuessLocation(dw.currentOcr, dw.Locations)

}
func (dw *Daywalker) TempScreenshot(name string) string {
	imgf := c.F("%v_%v.png", dw.fprefix, name)
	dw.lastscreenshot = c.TempFile(imgf)
	pt := dw.Screenshot(c.TempFile(imgf))
	return pt
}

func availiableToday(days string) bool {
	d := strings.Split(days, "/")
	weekday := time.Now().Weekday().String()
	return slices.Contains(d, weekday[:3])
}

// Press button, search for 'button's text in ocr results
func (dw *Daywalker) Press(b activities.Button) bool {
	log.Debugf("Known buttons: %+v", dw.knwnBtns)
	if len(dw.knwnBtns) > 0 {
		butt, e := dw.button(b.String())
		if e == nil {
			dw.NotifyUI(c.Cyan("BTN PRSD"), c.Green(c.F("%v |> %vx%v", butt, butt.X, butt.Y)))
			dw.Tap(butt.X, butt.Y, 1)
			return true
		}
	}

	res := dw.currentOcr.Result()
Lookin:
	x, y, e := LookForButton(res, b)
	if e != nil {
		res = dw.currentOcr.TryAgain()
		goto Lookin
	}
	// dw.knownBtns[b] = sessionBtn{x: x, y: y}
	nb := repository.NewBtn(b.String(), "", x, y, dw.Resolution.X, dw.Resolution.Y)
	dw.knwnBtns = append(dw.knwnBtns, nb)
	dw.Press(b)
	// Let location load
	time.Sleep(3 * time.Second)
	return true
}

func LookForButton(or []ocr.AlmoResult, b activities.Button) (x, y int, e error) {
	log.Debugf(c.Red("Looking for BTN: %v"), b.String()) //, c.RFW(or))
	if b.String() == "" {
		return 500, 100, nil
	}
	for _, r := range or {
		if strings.Contains(r.Linechars, b.String()) {
			xo, yo := b.Offset()
			return r.X + xo, r.Y + yo, nil
		}
	}
	return 0, 0, errors.New("btn not found here")
}

func (dw *Daywalker) button(s string) (*repository.Button, error) {
	for _, b := range dw.knwnBtns {
		if s == b.Name {
			return b, nil
		}
	}
	return &repository.Button{}, ErrUnknownButton
}
