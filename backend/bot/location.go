/* Location service, helpful functions to determinate where bot locates */
package bot

import (
	"mugowalker/backend/image"
	"strings"
)

var outFn func(string, string)

type Location interface {
	Id() string
	Keywords() []string
	HitThreshold() int
}

func GuessLocation(a *image.ImageProfile, locations []any) (locname string) {
	maxh := 1
	var resloc string
	var candidates []string
	for _, loc := range locations {
		l, ok := loc.(Location)
		if ok {
			hit := Intersect(a.TesseractResult(), l.Keywords())
			if len(hit) >= l.HitThreshold() && len(hit) >= maxh {
				maxh = len(hit)
				candidates = append(candidates, l.Id())
				resloc = l.Id()
			}
		}
	}
	if maxh == 1 {
		outFn("GUESS |>", f("Bad recognition -> %v ", "retry"))
		a.Redo()
	}

	outFn("GUESS |> ", f(" ↓ Location ↓ \n\t -->  Winner|> %v  Hits|> %v]\n\t --> candidates: %v", resloc, maxh, candidates))
	return resloc
}

func TextPosition(str string, alto []image.ScreenWord) (x, y int) {
	for _, v := range alto {
		if strings.Contains(v.S, str) {
			return v.X, v.Y
		}

	}
	return 0, 0
}

func Intersect(or []*image.ScreenWord, k []string) (r []string) {
NextLine:
	for _, v := range or {
		for _, kw := range k {
			if strings.Contains(v.S, kw) {
				r = append(r, v.S)
				continue NextLine
			}
		}
	}
	return r
}
