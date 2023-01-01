package main

import (
	"fmt"
	"image"
	"os"

	"github.com/fatih/color"

	"worker/afk"
	"worker/bot"
	"worker/cfg"
	"worker/ocr"

	"golang.org/x/exp/slices"
)

func main() {
	// log := cfg.Logger()

	if len(os.Args) > 1 && os.Args[1] == "-t" {
		// testloc("fn.png", afk.TowerFN)
		testloc("statlow.png", afk.STAT)
		// ocrtest()
		return

	}

	// USER INPUT DATA
	const (
		player    = "Devitool"
		game      = "afkarena"
		rTaskConf = "cfg/reactions.yaml"
		connect   = "localhost:5555"
	)
	user := cfg.User(player, game, rTaskConf, connect)

	device := cfg.Load(user)
	gm := afk.New(user.Account, user.Game)
	bt := bot.New(device, gm)
	task := user.Task(afk.TowerFN)

	bt.UP(task)
}

func ocrtest() {
	testlocs := make(map[string]string, 1)

	testlocs[afk.DARKFORREST] = "df.png"
	testlocs[afk.ENTRY] = "cpn.png"
	testlocs[afk.RANHORNY] = "h.png"
	testlocs[afk.BOSSTAGE] = "cpnb.png"
	testlocs[afk.KTower] = "towers.png"
	testlocs[afk.RESULT] = "lose.png"
	var improved, casual int = 0, 0
	for k, v := range testlocs {
		imp, cas := testloc(v, k)
		if imp {
			improved++
		}
		if cas {
			casual++
		}
	}
	color.HiBlue("\nTest overall:\n   Split   --> %v/%v\n   OneImg  --> %v/%v\n", improved, len(testlocs), casual, len(testlocs))
}

func testloc(img, loc string) (r1, r2 bool) {
	fail := color.New(color.FgHiRed, color.Bold).SprintfFunc()
	pass := color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	regular := color.New(color.FgHiYellow).SprintFunc()
	improved := color.New(color.FgHiCyan).SprintFunc()

	b := afk.New("afk", "test")
	loca := b.GetLocation(loc)
	restr := "\nResult	-> %v\nHits	-> [%v/%v]\n\n"

	fmt.Printf("Test location: [%v], source: %v\n\n", loca.Key, img)

	mt := ocr.TextExtract(img)
	ass := mt.Intersect(loca.Keywords)
	fmt.Print(regular("General	-> "), highlight(mt.Fields(), ass, pass))

	if r2 = len(ass) >= loca.Threshold; r2 {
		fmt.Print(pass(restr, r2, len(ass), loca.Threshold))
	} else {
		fmt.Print(fail(restr, r2, len(ass), loca.Threshold))
	}

	t := ocr.ImprovedTextExtract(img)
	ass1 := t.Intersect(loca.Keywords)
	fmt.Print(improved("Most Accurate:	-> "), highlight(t.Fields(), ass1, pass))
	if r1 = len(ass1) >= loca.Threshold; r1 {
		fmt.Print(pass(restr, r1, len(ass1), loca.Threshold))
	} else {
		fmt.Print(fail(restr, r1, len(ass1), loca.Threshold))
	}
	return
}

func testRegion() {
	//    b := afk.New("afk", "test")
	loc := "ch.png"
	p1 := image.Point{X: 400, Y: 1400}
	p2 := image.Point{X: 300, Y: 400}
	r := ocr.RegionText(loc, p1, p2)
	color.HiGreen("\nResult --> %v", r)
}

func highlight(k, s []string, fn func(s string, a ...interface{}) string) []string {
	for i, v := range k {
		if slices.Contains(s, v) {
			k[i] = fn(v)
		}
	}
	return k
}
