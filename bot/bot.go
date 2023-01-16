package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"worker/afk"
	"worker/cfg"

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
	//    origmagick = []string{"-colorspace",
	//        "Gray",
	//        "-alpha",
	//        "off",
	//        "-threshold",
	//        "75%",
	//    }
	user = cfg.ActiveUser()
	origocr    = user.Imagick
)

var ErrLocationMismatch = errors.New("Wrong Location!")

// var errActionFail = errors.New("smthg went wrong during Doing Action")

var (
	tempfile string = "temp"
	step     int    = 0
)

const (
	maxattempt uint8 = 3
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
		xmax:    xgrid, ymax: ygrid, cnt: 0, maxocrtry: 2,
	}
}

// ScanScreen OCRed Text TODO: maybe add args to peek like peek(data interface) smth like
// this should be w.ScanScreen(Location) \n w.ScanScreen(Stage)
func (dw *Daywalker) ScanScreen() []ocr.AltoResult { // ocr.Result {
	// TODO: Generate random filname
	s, e := dw.Screenshot(tempfile)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v ", e, "ScanScreen()")
	}
	// text := ocr.TextExtract(s)
	text := ocr.TextExtractAlto(s)
	log.Tracef("ocred: %v", text)
	color.HiCyan("%s", text)
	return text
}

func (dw *Daywalker) MyLocation() (locname string) {
WaitForLoc:
	for {
		if !dw.checkLoc(dw.ScanScreen()) {
			time.Sleep(8 * time.Second)
			if step >= dw.maxocrtry {
				color.HiRed("\rUsing improved ocr settings")
				user.UseAltImagick = true
				user.Tesseract = simpletess
				color.HiCyan("\rMagick args --> %v\n\r", user.AltImagick)
			}
			if step >= dw.maxocrtry+2 {
				color.HiRed("\rUsing RANDOM ocr settings xD ")
				user.AltImagick = ocr.MagickArgs()
				// dw.Back()
				color.HiCyan("\rMagick args --> %v\n\r", user.AltImagick)
				// log.Warnf("Magick args --> %v", user.AltImagick)
			}
			step++
			continue WaitForLoc
		} else {
			if step >= dw.maxocrtry {
				color.HiCyan("Returnin ocr params")
				user.UseAltImagick = false

			}
			step = 0
			break WaitForLoc

		}
	}
	color.Yellow("\rBest match -> %v\n\r", dw.lastLoc)
	// fmt.Printf("My Location most likely -> %v\n\r", dw.lastLoc)
	return dw.lastLoc.Key
}

func (dw *Daywalker) checkLoc(o []ocr.AltoResult) (ok bool) {
	maxh := 1
	for k, loc := range dw.Locations {
		hit := ocr.Intersect(o, loc.Keywords)
		if len(hit) >= loc.Threshold && len(hit) >= maxh {
			maxh = len(hit)
			// fmt.Printf("## Keywords hit %v -> %v ##...\n\r", loc.Key, hit)
			color.HiYellow("\r## hit %v -> %v ##...\n\r", loc.Key, hit)
			dw.lastLoc = &dw.Locations[k]
			ok = true
			// repository.RawLocData(loc.Key, strings.Join((), ";"))
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
	if user.DrawStep {
		drawTap(px, py, dw)
	}
	e := dw.Tap(fmt.Sprint(px), fmt.Sprint(py))
	// fmt.Printf("Tap: Grid-> %v:%v, Point-> %vx%v px\n\r", gx, gy, px, py)
	color.HiGreen("Tap: Grid-> %v:%v, Point-> %vx%v px\n\r", gx, gy, px, py)
	if e != nil {
		log.Warnf("Have an error during tap: %v", e.Error())
	}
}

func drawTap(tx, ty int, bot *Daywalker) {
	step++
	s, e := bot.Screenshot(fmt.Sprintf("%v", step))
	circle := fmt.Sprintf("circle %v,%v %v,%v", tx, ty, tx+20, ty+20)
	no := fmt.Sprintf("+%v+%v", tx-20, ty+20)
	cmd := exec.Command("magick", s, "-fill", "red", "-draw", circle, "-fill", "black", "-pointsize", "60", "-annotate", no, fmt.Sprintf("%v", step), cfg.UserFile(bot.lastscreenshot))
	e = cmd.Run()

	if e != nil {
		log.Errorf("s:%v", e.Error())
	}
	os.Remove(s)
}
