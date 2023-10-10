package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"mugowalker/backend/adb"
	"mugowalker/backend/image"
	"strings"
)

var ErrLocationMismatch = errors.New("wrong location")

var (
	f = fmt.Sprintf
)

type Bot struct {
	id    uint32
	outFn func(string, string)
	eyes  func(string) *image.ImageProfile
	*adb.Device
}

// New Instance of bot
func New(altout func(s1, s2 string), ocr *image.OcrEngine) *Bot {
	outFn = altout
	return &Bot{
		id:    rand.Uint32(),
		outFn: altout,
		eyes:  ocr.ExtractText,
	}
}

func (b *Bot) NotifyUI(pref, msg string) {
	b.outFn(pref, msg)
}

func (b *Bot) Text() *image.ImageProfile {
	s := b.screenshot()
	text := b.eyes(s)
	return text
}

// Tap x,y with y offset
func (b *Bot) FindTap(word string, x, y int) error {
	for _, v := range b.Text().TesseractResult() {
		if strings.Contains(v.S, word) {
			v.Offset(x, y)
			b.TapSW(v)
			b.outFn("tappin -> ", f("s:%s  %dx%d", v.S, v.X, v.Y))
			return nil
		}
	}
	return errors.New("UNTAPPED")
}
func (b *Bot) TapSW(screenword *image.ScreenWord) {
	e := b.Device.Tap(fmt.Sprint(screenword.X), fmt.Sprint(screenword.Y))
	if e != nil {
		b.outFn("Have an error during tap: ", e.Error())
	}
}

func (b *Bot) TapW(sw *image.ScreenWord) []*image.ScreenWord {
	b.TapSW(sw)
	return b.Text().TesseractResult()
}

func (b *Bot) Connect(target string) (isConnected bool) {
	dev, err := adb.Connect(target)
	s := fmt.Sprintf("[DEV]: %v", dev)
	b.outFn("ADB", s)
	if err == nil {
		b.Device = dev
		isConnected = true
	}
	return
}

func (b *Bot) AppStatus(app string) int {
	r := b.PS(app)
	if len(r) > 0 {
		return 1
	}
	return 0
}

func (b *Bot) screenshot() string {
	b.Screencap()
	file, err := b.Pull()
	if err != nil {
		b.outFn("BOT", err.Error())
	}
	return file
}
