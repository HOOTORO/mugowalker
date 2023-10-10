package image

import (
	"fmt"
	"image"
	"mugowalker/backend"
	c "mugowalker/backend/cfg"
	"strings"

	clr "github.com/go-color-term/go-color-term/coloring"
)

const (
	MagickCmd    = "magick"
	TesseractCmd = "tesseract"
)

type OcrEngine struct {
	*backend.Config
}

// ScreenWord parsed xml Almo
type ScreenWord struct {
	*image.Point
	S      string
	LineNo int
}

func SW(s string, xyln ...int) *ScreenWord {
	sw := &ScreenWord{S: s, LineNo: 0, Point: &image.Point{}}
	for i, v := range xyln {
		switch i {
		case 0:
			sw.X = v
		case 1:
			sw.Y = v
		case 2:
			sw.LineNo = v
		}
	}
	return sw
}

func (sw *ScreenWord) Offset(x, y int) {
	sw.X += x
	sw.Y += y
}
func (sw *ScreenWord) String() string {
	return fmt.Sprintf("%2d:%4dx%-4d | %10s\t", sw.LineNo, sw.X, sw.Y, clr.Green(sw.S))
}

var (
	log func(string)
)

func NewEngine(c *backend.Config) *OcrEngine {
	log = func(s string) { c.Log.Debug("[OCR Engine]:" + s) }
	// Fallback to searching on PATH.
	return &OcrEngine{c}

}
func (en *OcrEngine) MagickTransform(img string) (string, error) {
	transformedName := strings.Replace(img, ".png", "_prep.png", 1)
	magickArgs := make([]string, 0)
	magickArgs = append(magickArgs, img)
	magickArgs = append(magickArgs, en.Settings.Imagick.Args()...)
	magickArgs = append(magickArgs, transformedName)
	log(fmt.Sprintf("Imagick args: %v", magickArgs))
	e := c.RunCmd(MagickCmd, magickArgs)
	if e != nil {
		log(e.Error())
		return "", e
	}
	return transformedName, nil
}

// ExtractText prepare and extract text from img
func (en *OcrEngine) ExtractText(img string) *ImageProfile {
	// defer  timeTrack(time.Now(), "AltOcr")
	preparedImg, e := en.MagickTransform(img)
	if e != nil {
		en.Log.Error("IMAGE NOT PREPARED" + fmt.Sprintf("%v", e))
	}

	ip := &ImageProfile{
		original:   img,
		prepared:   preparedImg,
		recognized: make([]*ScreenWord, 0),
		psm:        en.Settings.Tesseract.Psm,
		ignored:    en.Settings.IgnoredWords,
	}
	return ip
}

// func (a ScreenWord) String() string {
// 	return fmt.Sprintf("#%d <%s>[%s] ", a.LineNo, fmt.Sprintf("%vx%v", a.X, a.Y), c.Shortener(a.S, 7))
// }

// func timeTrack(start time.Time, name string) string {
// 	elapsed := time.Since(start)
// 	return fmt.Sprintf("%v\n\r", fmt.Sprintf("\r[%s] %s", name, elapsed.Round(time.Millisecond)))
// }
