package bot

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"worker/adb"
	"worker/afk"
	"worker/ocr"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

var errActionFail = errors.New("smthg went wrong during Doing Action")
var step, n int = 0,0
var tempfile string = "step"

const (
	xgrid = 5
	ygrid = 18
)

// New Instance of bot
func New(d *adb.Device, game *afk.Game) *Daywalker {
    color.Magenta("Tap grid size : [ %vx%v ]", xgrid, ygrid)
	return &Daywalker{
		Device:  d,
		Game:    game,
		Tasks:   make([]Task, 0, 10),
		lastLoc: game.GetLocation(afk.ENTRY),
		xmax:    xgrid, ymax: ygrid,
	}
}

// Peek OCRed Text TODO: maybe add args to peek like peek(data interface) smth like
// this should be w.Peek(Location) \n w.Peek(Stage)
func (dw *Daywalker) Peek() string {
	// TODO: Generate random filname
	step++
	filename := fmt.Sprintf("%v_%v.png", step, tempfile)
	e := dw.Screenshot(filename)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v ", e, "Peek()")
	}
	abspath, _ := filepath.Abs(filename)
	text := ocr.Text(abspath)
    log.Tracef("ocred: %v", text)
	return text
}


// Tap screen, grid 5x14
func (dw *Daywalker) gridTap(x, y int) {
	// Cell size
	height := dw.Resolution.Y / dw.ymax
	width := dw.Resolution.X / dw.xmax

	// Center point
	tx := x*width - width/2
	ty := y*height - height/2

	e := dw.Tap(fmt.Sprint(tx), fmt.Sprint(ty))
	color.HiGreen("Tap Grid: in->%v:%v, px-> %vx%v px", x, y, tx, ty)
	showTap(tx, ty)

	if e != nil {
		log.Warnf("Have an error during tap: %v", e.Error())
	}
}


func (dw *Daywalker) GridTapOff(x, y, o int) {
    // Cell size
    height := dw.Resolution.Y / dw.ymax
    width := dw.Resolution.X / dw.xmax

    // Center point
    tx := x*width - width/2
    ty := y*height - (height*o)/2

    e := dw.Tap(fmt.Sprint(tx), fmt.Sprint(ty))
    color.HiGreen("Tap Grid: in->%v:%v (%v:%v), px-> %vx%v px", x, y, xgrid, ygrid, tx, ty)
    showTap(tx, ty)

    if e != nil {
        log.Warnf("Have an error during tap: %v", e.Error())
    }
}

func (dw *Daywalker) Do(action string) error {
	var cnt int8
    color.HiBlue("Current Daily GOAL : %v", action)
	a := dw.Action(action)


	if strings.Contains(a.MidlocId, "overlay") {
		for _, g := range a.OverlayGrids() {
			time.Sleep(time.Duration(a.Delay) * time.Second)
			dw.gridTap(g.X, g.Y)
		}
	}

	for !(dw.isLocation(a.FinlocId) || cnt > 10) {
		color.HiCyan("waiting %v location for apear...%v", a.FinlocId, cnt)
		time.Sleep(3 * time.Second)
		cnt++
	}
	return nil

}
// magick fr.png -fill red -draw "circle 500,1200 520,1220" -fill black -pointsize 60 -annotate +480+1220 '1' touch_4.png
func showTap(tx, ty int) {
	f := fmt.Sprintf("%v_%v.png", step, tempfile)
    n++
	circle := fmt.Sprintf("circle %v,%v %v,%v", tx, ty, tx+20, ty+20)
    no := fmt.Sprintf("+%v+%v",tx-20,ty+20)
//     no,
cmd := exec.Command("magick", f, "-fill", "red", "-draw", circle, "-fill","black", "-pointsize", "60", "-annotate", no, fmt.Sprintf("%v",n) ,f)
	e := cmd.Run()

	if e != nil {
		log.Errorf("s:%v", e.Error())
	}
}
//magick fr.png -fill red -draw "circle 500,1200 520,1230" -fill black -pointsize 70  -draw "text 500,1230 '1'" -annotate +500+1255 '1' touch_4.png
// MyLocation TODO: remove, overkill i think
func (dw *Daywalker) MyLocation() (locname string) {
	var maxh int = 1

	text := dw.Peek()
	recognizedText := ocr.OCRFields(text)
	color.HiYellow("## ocred -> %v ## \nFind matches...", recognizedText)
	for k, loc := range dw.Locations {
		hits := ocr.KeywordHits(loc.Keywords, recognizedText)
		if hits >= loc.Threshold && hits > maxh {
			maxh = hits
			dw.lastLoc = &dw.Locations[k]
		}
	}
	color.HiYellow("My Location most likely -> %v", dw.lastLoc.Key)
	return dw.lastLoc.Key
}

func (dw *Daywalker) isLocation(locname string) bool {
	var ok bool
	text := dw.Peek()
	loc := dw.GetLocation(locname)
	recognizedText := ocr.OCRFields(text)
	hits := ocr.KeywordHits(loc.Keywords, recognizedText)

	if hits >= loc.Threshold {
		dw.lastLoc = loc
		ok = true
	}
	log.Tracef("Location -> %v ; Hits -> %v, %v", locname, hits, ok)
	color.HiYellow("### Is %v? ### -> %v", locname, ok)
	return ok
}
