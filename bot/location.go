/* Location service, helpful functions to determinate where bot locates */
package bot

import (
	"strings"
	"worker/cfg"
	"worker/ocr"

	tea "github.com/charmbracelet/bubbletea"
)

var red, green, cyan, ylw, mgt func(...interface{}) string

type Location interface {
	Id()
	Keywords() []string
}

func GuessLocByKeywords(a []ocr.AltoResult, locations []cfg.Location) (locname string) {
	maxh := 1
	var resloc string
	for _, loc := range locations {
		hit := Intersect(a, loc.Keywords)
		if len(hit) >= loc.Threshold && len(hit) >= maxh {
			maxh = len(hit)
			// log.Debugf(ylw(f("hit: %v -> %v \n", loc.Key, hit)))
			log.Warn(mgt("GUESSHI |> "), ylw(f("Location: %v -> %v \n\r", loc.Key, hit)))
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
