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

	if len(os.Args) > 1 && os.Args[1] == "-t" {
		color.HiRed("%v", "TEST RUN")
		ocrtest()
		return
	}

	// USER INPUT DATA
	const (
		player  = "Devitool"
		game    = "afkarena"
		connect = "localhost:5555"
	)
	var rTaskConf = []string{"cfg/reactions.yaml", "cfg/daily.yaml"}
	user := cfg.User(player, game, connect, rTaskConf)

	device := cfg.Load(user)
	gm := afk.New(user)
	bt := bot.New(device, gm)

	choice := cfg.UserInput("What bot should do?\n0. Run all (Default)\n1. Run daily?\n2. Push campain?\n3-7. Push towers (KT,L,M,W,G)", "0")

	switch choice {
	case "0":
		bt.UpAll()
	case "1":
		bt.Daily()
	case "2":
		push := bt.Task(afk.DOPUSHCAMP)
		bt.React(push)
	case "3":
		kt := bt.Task(afk.Kings)
		bt.React(kt)
	case "4":
		kt := bt.Task(afk.Light)
		bt.React(kt)
	case "5":
		kt := bt.Task(afk.Mauler)
		bt.React(kt)
	case "6":
		kt := bt.Task(afk.Wilder)
		bt.React(kt)
	case "7":
		kt := bt.Task(afk.Graveborn)
		bt.React(kt)
    default:
        color.HiRed("DATS WRONG NUMBA MAFAKA!")
	}
}

func ocrtest() {
	b := afk.New(&cfg.UserProfile{Account: "test", Game: "afk", ConnectionStr: "localhost:5555", TaskConfigs: []string{"cfg/reactions.yaml"}})

	var testdata = func(lo uint, im string) *struct {
		loc afk.ArenaLocation
		img string
	} {
		return &struct {
			loc afk.ArenaLocation
			img string
		}{loc: afk.ArenaLocation(lo), img: im}
	}
	testlocs := make([]*struct {
		loc afk.ArenaLocation
		img string
	}, 0)
	testlocs = append(testlocs,
		//        testdata(afk.DARKFORREST, "test/forrest.png"),
		//        testdata( afk.Campain, "test/cpn1.png"),
		//        testdata( afk.Campain, "test/cpn2.png"),
		testdata(afk.RANHORNY.Id(), "_test/h.png"),
		//        testdata( afk.BOSSTAGE, "test/cpnb.png"),
		//        testdata( afk.Kings.String(), "test/towers.png"),
		//        testdata( afk.RESULT, "test/lose.png"),
		//        testdata( afk.Loot.String(), "test/loot.png"),
		//        testdata( afk.FastReward.String(),"test/fr.png"),
		//        testdata( afk.BATTLE ,"test/btl_multstg.png"),
		//        testdata( afk.BATTLE ,"test/btl_onestg.png"),
		//        testdata( afk.STAT ,"test/stt1.png"),
		//        testdata( afk.STAT ,"test/stt2.png"),
		//        testdata( afk.WIN , "test/cpn_win.png"),
		//        testdata( "", ""),
		//        testdata( "", ""),
		//        testdata( "", ""),
		//        testdata( "", ""),
	)

	var overall = 0
	for _, v := range testlocs {
		res := testloc(v.img, b.GetLocation(v.loc))
		if res {
			overall++
		}

	}
	//    testRegion("test/btl_onestg_1.png")
	//    testRegion("test/btl_multstg_1.png")
	color.HiBlue("\nTest overall:\n   Basic   --> %v/%v\n", overall, len(testlocs))
}

func testloc(img string, loc *cfg.Location) (r1 bool) {
	fail := color.New(color.FgHiRed, color.Bold).SprintfFunc()
	pass := color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	regular := color.New(color.FgHiYellow).SprintFunc()
	regular("%v - %v", img, loc)

	restr := "\nResult	-> %v\nHits	-> [%v/%v]\n\n"

	mt := ocr.TextExtract(img)

	fmt.Printf("Test location: [%v], source: %v\n\n", fail(loc.Key), fail(img))
	ass := mt.Intersect(loc.Keywords)

	fmt.Print(regular("General	-> "), highlight(mt.Fields(), ass, pass))

	if r1 = len(ass) >= loc.Threshold; r1 {
		fmt.Print(pass(restr, r1, len(ass), loc.Threshold))
	} else {
		fmt.Print(fail(restr, r1, len(ass), loc.Threshold))
	}
	pass("xu")
	al := ocr.TextExtractAlto(img)
	//    fmt.Printf("%v", pass("%v",))
	tl := al.Layout.Page.PrintSpace.ComposedBlock.TextBlock.TextLine
	for _, line := range tl {
		fmt.Printf("%v : %v", fail("%vx%v", line.HPOS, line.VPOS), pass("%v\n", line.String))
	}
	return
}

func testRegion(img string) {
	//    b := afk.New("afk", "test")
	loc := img
	p1 := image.Point{X: 120, Y: 1500}
	p2 := image.Point{X: 600, Y: 180}
	r := ocr.RegionText(loc, p1, p2)
	color.HiGreen("\nResult --> %v", r.String())
}

func highlight(k, s []string, fn func(s string, a ...interface{}) string) []string {
	for i, v := range k {
		if slices.Contains(s, v) {
			k[i] = fn(v)
		}
	}
	return k
}
