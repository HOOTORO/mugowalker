package bot

import (
    "fmt"
    "path/filepath"

    "worker/adb"
    "worker/afk"
    "worker/cfg"
    "worker/ocr"

    "github.com/fatih/color"
    log "github.com/sirupsen/logrus"
)

// New Instance of bot
func New(d *adb.Device, game *afk.Game) *Daywalker {
	return &Daywalker{
		Device:     d,
		Tasks:      make([]Task, 0, 10),
		CurrentLoc: game.GetLocation(afk.ENTRY),
	}
}

// Peek OCRed Text TODO: maybe add args to peek like peek(data interface) smth like
// this should be w.Peek(Location) \n w.Peek(Stage)
func (dw *Daywalker) Peek() string {
	// TODO: Generate random filname
	filename := "p.png"
	dw.Screencap(filename)
	e := dw.Pull(filename, ".")
	if e != nil {
		fmt.Printf("\nerr:%v\nduring run:%v ", e, "Peek()")
	}
	abspath, _ := filepath.Abs(filename)

	text := ocr.Text(abspath)
	return text
}

func isLocation(loc *cfg.Location, b *Daywalker) bool {
	var ok bool
	recognizedText := ocr.OCRFields(b.Peek())
	color.HiYellow("### Where we? ###\n ## %v ## \n", recognizedText)
	hits := ocr.KeywordHits(loc.Keywords, recognizedText)
	if hits >= loc.Threshold {
		ok = true
	}
	return ok
}

// Tap screen, grid 5x14
func (dw *Daywalker) gridTap(x, y int) {
    // Cell size
	height := dw.WmSize.Y / 14
	width := dw.WmSize.X / 5

    // Center point
	tx := x*width - width/2
	ty := y*height - height/2

	e := dw.Tap(fmt.Sprint(tx), fmt.Sprint(ty))
    if e != nil {
        log.Warnf("Have an error during tap: %v", e.Error())
    }
}

func(dw *Daywalker) Do(a *cfg.Action){
    
}
