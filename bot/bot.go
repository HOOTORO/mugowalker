package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"worker/afk"
	"worker/cfg"

	"worker/afk/repository"

	"worker/adb"
	"worker/ocr"

	"github.com/fatih/color"
)

type Offset int

const (
	Center Offset = iota
	Bottom
	Top
)

var (
	impmagick = []string{
		"-colorspace",
		"Gray",
		"-alpha",
		"off",
		"-threshold",
		"85%",
		"-edge",
		"1",
		"-negate",
		"-blur",
		"1x1",
		"-black-threshold",
		"80%",
	}
	simpletess = []string{"--psm", "3", "-c", "tessedit_create_alto=1", "quiet"}
	origocr    = cfg.Env.Imagick
	//    origmagick = []string{"-colorspace",
	//        "Gray",
	//        "-alpha",
	//        "off",
	//        "-threshold",
	//        "75%",
	//    }
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
		lastLoc: game.GetLocation(afk.Campain),
		xmax:    xgrid, ymax: ygrid, cnt: 0, maxocrtry: 5,
	}
}

// ScanScreen OCRed Text TODO: maybe add args to peek like peek(data interface) smth like
// this should be w.ScanScreen(Location) \n w.ScanScreen(Stage)
func (dw *Daywalker) ScanScreen() ocr.Result {
	// TODO: Generate random filname
	s, e := dw.Screenshot(tempfile)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v ", e, "ScanScreen()")
	}
	text := ocr.TextExtract(s)
	log.Tracef("ocred: %v", text)
	return text
}

func (dw *Daywalker) MyLocation() (locname string) {
WaitForLoc:
	for {
		if !dw.checkLoc(dw.ScanScreen()) {
			time.Sleep(8 * time.Second)
			if step >= dw.maxocrtry {
				color.HiRed("Using improved ocr settings")
				cfg.Env.Imagick = impmagick
				cfg.Env.Tesseract = simpletess
			}
			if step >= dw.maxocrtry+2 {
				dw.Back()
			}
			step++
			continue WaitForLoc
		} else {
			if step >= dw.maxocrtry {
				color.HiCyan("Returnin ocr params")
				cfg.Env.Imagick = origocr

			}
			step = 0
			break WaitForLoc

		}
	}
	color.HiYellow("My Location most likely -> %v", dw.lastLoc)
	return dw.lastLoc.Key
}

func (dw *Daywalker) checkLoc(o ocr.Result) (ok bool) {
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
	cmd := exec.Command("magick", s, "-fill", "red", "-draw", circle, "-fill", "black", "-pointsize", "60", "-annotate", no, fmt.Sprintf("%v", step), cfg.UsrDir(bot.lastscreenshot))
	e = cmd.Run()

	if e != nil {
		log.Errorf("s:%v", e.Error())
	}
	os.Remove(s)
}
