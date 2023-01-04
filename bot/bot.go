package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"worker/afk/repository"
	"worker/cfg"

	"worker/adb"
	"worker/afk"
	"worker/ocr"

	"github.com/fatih/color"
)

type Offset int

const (
	Center Offset = iota
	Bottom
	Top
)

var ErrLocationMismatch = errors.New("Wrong Location!")

// var errActionFail = errors.New("smthg went wrong during Doing Action")

var (
	tempfile string = "temp"
	step     int    = 0
)

const (
	maxattempt uint8 = 5
	xgrid            = 5
	ygrid            = 18
)

// New Instance of bot
func New(d *adb.Device, game *afk.Game) *Daywalker {
	color.Magenta("Tap grid size : [ %vx%v ]", xgrid, ygrid)
	rand.Seed(time.Now().Unix())
	return &Daywalker{
		id:      rand.Uint32(),
		fprefix: time.Now().Format("2006_01"),
		Device:  d,
		Game:    game,
		lastLoc: game.GetLocation(afk.ENTRY),
		xmax:    xgrid, ymax: ygrid, cnt: 0,
	}
}

// Peek OCRed Text TODO: maybe add args to peek like peek(data interface) smth like
// this should be w.Peek(Location) \n w.Peek(Stage)
func (dw *Daywalker) Peek() ocr.OcrResult {
	// TODO: Generate random filname
	s, e := dw.Screenshot(tempfile)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v ", e, "Peek()")
	}
	text := ocr.TextExtract(s)
	log.Tracef("ocred: %v", text)
	return text
}

func (dw *Daywalker) MyLocation() (locname string) {
WaitForLoc:
	for {
		if !dw.checkLoc(dw.Peek()) {
			time.Sleep(5 * time.Second)
			continue WaitForLoc
		} else {
			break WaitForLoc
		}
	}
	color.HiYellow("My Location most likely -> %v", dw.lastLoc)
	return dw.lastLoc.Key
}

func (dw *Daywalker) Do(action string) error {
	dw.ActiveTask = action
	color.HiBlue("Executing action : %v", action)
	a := dw.Action(action)

	if len(a.Start) > 0 && !slices.Contains(a.Start, dw.CurrentLocation()) {
		dw.ZeroPosition()
	}

	for i, g := range a.Steps {
		step = i + 1
		dw.makeStep(g)
	}
	color.HiGreen("Action done -> %v", action)
	if a.Next != "" {
		e := dw.Do(a.Next)
		return e
	}
	return nil
}

func (dw *Daywalker) checkLoc(o ocr.OcrResult) (ok bool) {
	maxh := 1
	for k, loc := range dw.Locations {
		hit := o.Intersect(loc.Keywords)
		if len(hit) >= loc.Threshold && len(hit) >= maxh {
			maxh = len(hit)
			color.HiYellow("## Keywords hit %v -> %v ##...", loc.Key, hit)
			dw.lastLoc = &dw.Locations[k]
			ok = true
			repository.RawLocData(loc.Key, strings.Join(o.Fields(), ";"))
		}
	}
	return
}

func (dw *Daywalker) makeStep(s cfg.Step) bool {
	var cnt uint8
	px, py, off := s.Target().X, s.Target().Y, s.Target().Offset

	dw.TapGO(px, py, off)

	if s.Skiplocheck {
		return true
	}

	if s.Wait {
		for !(slices.Contains(s.Loc, dw.MyLocation()) || cnt > maxattempt) {
			color.HiCyan("waiting %v location for apear...%v", s.Loc, cnt)
			time.Sleep(2 * time.Second)
			cnt++
		}
		if len(s.Loc) > 1 {
			conditionalAction := dw.Action(dw.ActiveTask).ConditionalStep(dw.CurrentLocation())
			dw.makeStep(conditionalAction)
		}
	}

	if !slices.Contains(s.Loc, dw.MyLocation()) {
		color.HiRed("%v\nWant: %v Have: %v", ErrLocationMismatch.Error(), s.Loc, dw.lastLoc)
		return false
	}
	return true
}

// tapGrid screen, grid 5x1, default Yoffset = 2
func (dw *Daywalker) tapGrid(x, y int) {
	yo := 1
	dw.TapGO(x, y, yo)
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

	e := dw.Tap(fmt.Sprint(px), fmt.Sprint(py))
	color.HiGreen("Tap: Grid-> %v:%v, Point-> %vx%v px", gx, gy, px, py)
	if e != nil {
		log.Warnf("Have an error during tap: %v", e.Error())
	}
}
