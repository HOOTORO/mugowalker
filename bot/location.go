/* Location service, helpful functions to determinate where bot locates */
package bot

import (
	"strings"
	"worker/cfg"
	"worker/ocr"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

var red, green, cyan, ylw, mgt func(...interface{}) string

func init() {
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	ylw = color.New(color.FgHiYellow).SprintFunc()
	mgt = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()
}
func (dw *Daywalker) ScanScreen() []ocr.AltoResult { // ocr.Result {
	s, e := dw.Screenshot(tempfile)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v ", e, "ScanScreen()")
	}
	text := ocr.TextExtractAlto(s)
	log.Tracef("ocred: %v", cyan(text))
	Fnotify("ALTORES", f("%v", text))
	return text
}

func GuessLocByKeywords(a []ocr.AltoResult, locations []cfg.Location) (locname string) {
	maxh := 1
	var resloc string
	for _, loc := range locations {
		hit := Intersect(a, loc.Keywords)
		if len(hit) >= loc.Threshold && len(hit) >= maxh {
			maxh = len(hit)
			log.Debugf(ylw("\rhit %v -> %v \n\r", loc.Key, hit))
			Fnotify("GUESSHIT", ylw("hit %v -> %v \n\r", loc.Key, hit))
			resloc = loc.Key
		}
	}
	return resloc
}

func TextPosition(str string, alto []ocr.AltoResult) (x, y int) {
	for _, v := range alto {
		if strings.Contains(v.Linechars, str) {
			return v.X, v.Y
		}

	}
	return 0, 0
}
func Intersect(or []ocr.AltoResult, k []string) (r []string) {
NextLine:
	for _, v := range or {
		for _, kw := range k {
			if strings.Contains(v.Linechars, kw) {
				r = append(r, v.Linechars)
				continue NextLine
			}
		}
	}
	return r

}

type notify struct{ s string }

func GuessNotify(str chan interface{}) tea.Msg {
	return func(s string) {
		str <- notify{s}
	}
}
