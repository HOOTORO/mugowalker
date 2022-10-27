package bot

import (
	"regexp"
	"strconv"

	"worker/ocr"

	// "github.com/Hootsdev/afkhelperuploader"
	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

func KeywordHits(kw, ocr []string) int {
	res := 0
	for _, word := range kw {
		if slices.Contains(ocr, word) {
			res++
		}
	}
	return res
}

func Regex(s, r string) (res []int) {
	re := regexp.MustCompile(r)
	for _, v := range re.FindStringSubmatch(s) {
		i, err := strconv.Atoi(v)
		if err == nil {
			res = append(res, i)
		}
	}
	return
}

func (d *Daywalker) WhereIs(locs map[string]Location) Location {
	current := ocr.OCRFields(d.Peek())
	color.HiYellow("##### Where we? ##############################\n## %v ##\n", current)
	maxhits, locName := 0, ""
	for name, v := range locs {
		hits := KeywordHits(v.Keywords, current)
		if hits > maxhits {
			maxhits = hits
			locName = name
		}

	}
	if locName != "" {
		color.HiBlue("######## %v ########\n", locName)
	}
	d.SetLocation(locs[locName])
	return locs[locName]
}
