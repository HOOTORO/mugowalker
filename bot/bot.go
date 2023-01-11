package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"worker/adb"
	"worker/ocr"
	// "github.com/fatih/color"
)

// New Instance of bot
func New(d *adb.Device, game *afk.Game) *Daywalker {
	color.Magenta("Tap grid size : [ %vx%v ]", xgrid, ygrid)
	rand.Seed(time.Now().Unix())
// New Instance of bot
func New(d *adb.Device, game *afk.Game) *Daywalker {
	color.Magenta("Tap grid size : [ %vx%v ]", xgrid, ygrid)
	rand.Seed(time.Now().Unix())
	return &Daywalker{
		id:      rand.Uint32(),
		fprefix: time.Now().Format("2006_01"),
		Device:  d,
		Game:    game,
		lastLoc: game.GetLocation(afk.Campain),
		xmax:    xgrid, ymax: ygrid, cnt: 0, maxocrtry: 5,
	}
}

// ScanScreen OCRed Text TODO: maybe add args to peek like peek(data interface) smth like
// this should be w.ScanScreen(Location) \n w.ScanScreen(Stage)
func (dw *Daywalker) ScanScreen() ocr.Result {
	// TODO: Generate random filname
	s, e := dw.Screenshot(tempfile)
	s, e := dw.Screenshot(tempfile)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v ", e, "ScanScreen()")
		log.Errorf("\nerr:%v\nduring run:%v ", e, "ScanScreen()")
	}
	text := ocr.TextExtract(s)
	log.Tracef("ocred: %v", text)
	text := ocr.TextExtract(s)
	log.Tracef("ocred: %v", text)
	return text
}

func (d *Daywalker) AllowedAction(n string) bool {
	_, ok := locs[n]
	return ok
}

func (d *Daywalker) SetLocation(s string) {
	d.loc = locs[s]
	(d.loc).Label = s
}

func (d *Daywalker) Action(s string, props Properties) error {
	action, ok := d.loc.Actions[s]
	if !ok {
		return errors.New(fmt.Sprintf("NO Action<%v> in context<%v>!", s, d.loc.Label))
	}
	return
}

// TapGO Grid x,y with y offset
func (dw *Daywalker) TapGO(gx, gy, off int) {
	o := Offset(off)
	// Cell size
	height := dw.Resolution.Y / int(dw.ymax)
	width := dw.Resolution.X / int(dw.xmax)

	// Center point
	px := gx*width - width/2
	py := gy*height - int(o)*height/2
	if cfg.Env.DrawStep {
		drawTap(px, py, dw)
	}
	e := dw.Tap(fmt.Sprint(px), fmt.Sprint(py))
	color.HiGreen("Tap: Grid-> %v:%v, Point-> %vx%v px", gx, gy, px, py)
	if e != nil {
		log.Warnf("Have an error during tap: %v", e.Error())
	}
}

func drawTap(tx, ty int, bot *Daywalker) {
	step++
	s, e := bot.Screenshot(fmt.Sprintf("%v", step))
	circle := fmt.Sprintf("circle %v,%v %v,%v", tx, ty, tx+20, ty+20)
	no := fmt.Sprintf("+%v+%v", tx-20, ty+20)
	cmd := exec.Command("magick", s, "-fill", "red", "-draw", circle, "-fill", "black", "-pointsize", "60", "-annotate", no, fmt.Sprintf("%v", step), filepath.Join("steps", filepath.Base(bot.lastscreenshot)))
	e = cmd.Run()

	if e != nil {
		log.Errorf("s:%v", e.Error())
	}
	os.Remove(s)
}
