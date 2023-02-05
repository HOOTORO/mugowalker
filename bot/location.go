/* Location service, helpful functions to determinate where bot locates */
package bot

import (
	"strings"
	c "worker/cfg"
	"worker/ocr"

	tea "github.com/charmbracelet/bubbletea"
)

var outFn func(string, string)

type Location interface {
	Id() string
	Keywords() []string
	HitThreshold() int
}

func GuessLocation(a []ocr.AlmoResult, locations []any) (locname string) {

	maxh := 1
	var resloc string
	var candidates []string
	for _, loc := range locations {
		l, ok := loc.(Location)
		if ok {
			hit := Intersect(a, l.Keywords())
			if len(hit) >= l.HitThreshold() && len(hit) >= maxh {
				maxh = len(hit)
				candidates = append(candidates, l.Id())
				resloc = l.Id()
			}
		}
	}
	outFn(c.Mgt("GUESSHI |>"), c.Ylw(f("Location -> %v [hits:%v]", resloc, maxh)))

	log.Debug(c.Mgt("GUESSHI |> "), c.Ylw(f(" ↓ Location ↓ \n\t -->  Winner|> %v  Hits|> %v]\n\t --> candidates: %v", resloc, maxh, candidates)))
	return resloc
}

func TextPosition(str string, alto []ocr.AlmoResult) (x, y int) {
	for _, v := range alto {
		if strings.Contains(v.Linechars, str) {
			return v.X, v.Y
		}

	}
	return 0, 0
}

func Intersect(or []ocr.AlmoResult, k []string) (r []string) {
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
