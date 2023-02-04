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
	xgrid         int   = 5
	ygrid         int   = 18
)

type Bot interface {
	Tap(x, y, off int)
	Location() string
	Screenshot(string) string
	ScanText() []ocr.AltoResult
	DiscoverDevices() []*adb.Device
	Connect(*adb.Device)
}

type BasicBot struct {
	id           uint32
	xgrid, ygrid int
	location     string
	outFn        func(string, string)
	*adb.Device
}

// New Instance of bot
func New(altout func(s1, s2 string)) *BasicBot {
	outFn = altout
	rand.Seed(time.Now().Unix())
	return &BasicBot{
		id:    rand.Uint32(),
		outFn: altout,
		xgrid: xgrid,
		ygrid: ygrid,
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

func (b *BasicBot) NotifyUI(pref, msg string) {
	b.outFn(pref, msg)
}

func (b *BasicBot) Location() (locname string) {
	return b.location
}

func (b *BasicBot) ScanText() []ocr.AltoResult { // ocr.Result {
	s := b.Screenshot(tempfile)
	text := ocr.TextExtractAlto(s)

	// log.Trace(green("OCR-R"), f("Words Onscr: %v lns: %s\nocred: %v", cyan(len(text)), green(text[len(text)].LineNo), cyan(z(text))))
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

// Tap x,y with y offset
func (b *BasicBot) Tap(gx, gy, off int) {

	e := b.Device.Tap(fmt.Sprint(gx), fmt.Sprint(gy))
	// b.outFn(mgt("BOT"), green(f("Tap -> %vx%v px", gx, gy)))
	if e != nil {
		log.Warnf("Have an error during tap: %v", e.Error())
	}
}

func (b *BasicBot) DiscoverDevices() []*adb.Device {
	devs, _ := adb.Devices()
	// Try to connect Bluestacks standart host:port
	bs, e := adb.Connect("localhost:5555")
	if e == nil {
		devs = append(devs, bs)
	}
	// Try to connect Nox standart host:port
	nox, e := adb.Connect("localhost:62001")
	if e == nil {
		devs = append(devs, nox)
	}
	return devs
}

func (b *BasicBot) Connect(d *adb.Device) {
	b.Device = d
}

func (b *BasicBot) IsAppRunnin(app string) int {
	r := b.PS(app)
	if len(r) > 0 {
		return 1
	}
	return 0
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

func (b *BasicBot) OcResult() []ocr.AltoResult {
	return b.ScanText()
}
