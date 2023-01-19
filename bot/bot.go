package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"worker/adb"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type offset int

const (
	Center offset = iota
	Bottom
	Top
)
const (
	ocrs = "OCR"
	mgc  = "MAGIC"
	tess = "TESSERACT"
)

var (
	user    = cfg.ActiveUser()
	origocr = user.Imagick
)

var ErrLocationMismatch = errors.New("wrong location")

// var errActionFail = errors.New("smthg went wrong during Doing Action")

var (
	tempfile     = "temp"
	step     int = 0
	log      *logrus.Logger
	f        = fmt.Sprintf
)

const (
	startlocation       = "universe"
	maxattempt    uint8 = 3
	xgrid               = 5
	ygrid               = 18
)

type Bot interface {
	Tap(x, y, off int)
	Location() string
	Screenshot(string) string
	ScanText() []ocr.AltoResult
}

type BasicBot struct {
	id           uint32
	xgrid, ygrid uint8
	location     string
	outFn        func(string, string)
	*adb.Device
}

// New Instance of bot
func New(d *adb.Device, altout func(s1, s2 string)) *BasicBot {
	rand.Seed(time.Now().Unix())
	return &BasicBot{
		id:       rand.Uint32(),
		location: startlocation,
		Device:   d,
		outFn:    altout,
	}
}
func init() {
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	ylw = color.New(color.FgHiYellow).SprintFunc()
	mgt = color.New(color.FgHiMagenta).SprintFunc()
	log = cfg.Logger()
}

func (b *BasicBot) Location() (locname string) {
	return b.location
	// WaitForLoc:
	// 	for {
	// 		if !dw.checkLoc(dw.ScanScreen()) {
	// 			time.Sleep(8 * time.Second)
	// 			if step >= dw.maxocrtry {
	// 				Fnotify(ocrs, f(red("\rUsing improved ocr settings")))
	// 				user.UseAltImagick = true
	// 				user.Tesseract = simpletess
	// 				Fnotify("ocr", f(cyan("\rMagick args --> %v\n\r", user.AltImagick)))
	// 			}
	// 			if step >= dw.maxocrtry+2 {
	// 				Fnotify(ocrs, red("\rUsing RANDOM ocr settings xD "))
	// 				user.AltImagick = ocr.MagickArgs()
	// 				// dw.Back()
	// 				Fnotify(mgc, cyan("\rMagick args --> %v\n\r", user.AltImagick))
	// 				// log.Warnf("Magick args --> %v", user.AltImagick)
	// 			}
	// 			step++
	// 			continue WaitForLoc
	// 		} else {
	// 			if step >= dw.maxocrtry {
	// 				Fnotify(ocrs, cyan("Returnin ocr params"))
	// 				user.UseAltImagick = false

	// 			}
	// 			step = 0
	// 			break WaitForLoc

	//		}
	//	}
	//
	// Fnotify("GUESSLOC", ylw("\rBest match -> %v\n\r", dw.lastLoc))
	// // fmt.Printf("My Location most likely -> %v\n\r", dw.lastLoc)
	// return dw.lastLoc.Key
}

// func (dw *Daywalker) checkLoc(o []ocr.AltoResult) (ok bool) {

//		return
//	}
func (b *BasicBot) ScanText() []ocr.AltoResult { // ocr.Result {
	s := b.Screenshot(tempfile)
	text := ocr.TextExtractAlto(s)
	z := func(arr []ocr.AltoResult) string {
		var s string
		for i, elem := range arr {
			if i%2 == 0 {
				s += f("	%-36s	", elem)
			} else {
				s += f("	%36s	\n", elem)
			}
		}
		return s
	}
	// log.Tracef("ocred: %v", cyan(text))
	b.outFn("ALTO OCR RESULT", f("%v", cyan(z(text))))
	return text
}

func (b *BasicBot) Screenshot(name string) string {
	var p, n string
	if filepath.IsAbs(name) {
		p, n = filepath.Split(name)
	} else {
		p = cfg.UserFile("")
	}
	newn := f("%v_%v.png", b.id, n)

	b.Screencap(newn)
	b.Pull(newn, p)
	return filepath.Join(p, newn)
}

// TapGO Grid x,y with y offset
func (b *BasicBot) Tap(gx, gy, off int) {
	o := offset(off)
	// Cell size
	height := b.Resolution.Y / int(b.ygrid)
	width := b.Resolution.X / int(b.xgrid)

	// Center point
	px := gx*width - width/2
	py := gy*height - int(o)*height/2
	if user.DrawStep {
		drawTap(px, py, b)
	}
	e := b.Device.Tap(fmt.Sprint(px), fmt.Sprint(py))
	// fmt.Printf("Tap: Grid-> %v:%v, Point-> %vx%v px\n\r", gx, gy, px, py)
	b.outFn("BOT", f(green("Tap: [Grid > %v:%v] Point-> %vx%v px\n\r", gx, gy, px, py)))
	if e != nil {
		log.Warnf("Have an error during tap: %v", e.Error())
	}
}

func drawTap(tx, ty int, bot Bot) {
	step++
	s := bot.Screenshot(f("%v", step))
	circle := fmt.Sprintf("circle %v,%v %v,%v", tx, ty, tx+20, ty+20)
	no := fmt.Sprintf("+%v+%v", tx-20, ty+20)
	cmd := exec.Command("magick", s, "-fill", "red", "-draw", circle, "-fill", "black", "-pointsize", "60", "-annotate", no, f("%v", step), cfg.UserFile(""))
	e := cmd.Run()

	if e != nil {
		log.Errorf("s:%v", e.Error())
	}
	os.Remove(s)
}
