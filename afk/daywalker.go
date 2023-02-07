package afk

import (
	"errors"
	"strings"
	"time"

	"worker/afk/activities"
	"worker/bot"
	c "worker/cfg"
	"worker/ocr"

	"golang.org/x/exp/slices"
)

const (
	alto = "ALTO"
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
	return &Daywalker{
		BasicBot: b,
		Game:     g,
		fprefix:  time.Now().Format("2006_01"),
		cnt:      0, maxocrtry: 2,
		knownBtns: make(map[activities.Button]sessionBtn, 0),
	}
}

var outFn func(string, string)

func init() {
	log = c.Logger()

}

func (dw *Daywalker) String() string {
	return c.F("Bot status:\n   Game: %v\n ActiveTask: %v", dw.Game, dw.ActiveTask)
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
	// or := dw.ScanText()
	log.Debugf("Known buttons: %+v", dw.knownBtns)
	btn, ok := dw.knownBtns[b]
	if ok {
		dw.NotifyUI(c.Cyan("BTN PRSD"), c.Green(c.F("%v |> %vx%v", b.String(), btn.x, btn.y)))
		dw.Tap(btn.x, btn.y, 1)
		return true
	}
	x, y, e := LookForButton(dw.currentOcr, b)
	if e != nil {
		return false
	}
	dw.knownBtns[b] = sessionBtn{x: x, y: y}
	dw.Press(b)
	// Let location load
	time.Sleep(3 * time.Second)
	return true
}

func LookForButton(or *ocr.ImageProfile, b activities.Button) (x, y int, e error) {
	log.Debugf(c.Red("Looking for BTN: %v"), b.String())
	if b.String() == "" {
		return 500, 100, nil
	}
	for _, r := range or.Result() {
		if strings.Contains(r.Linechars, b.String()) {
			xo, yo := b.Offset()
			return r.X + xo, r.Y + yo, nil
		}
	}
	return 0, 0, errors.New("btn not found here")
}
