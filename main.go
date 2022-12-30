package main

import (
	"fmt"
	"github.com/fatih/color"
	"image"

	//	"image"
	"os"
	"strconv"
	"worker/adb"
	"worker/afk"
	"worker/bot"
	"worker/cfg"
	"worker/ocr"
)

func main() {
	const (
		game = "afkarena"
	)
	log := cfg.Logger()
	log.Warnf("cmd args --> %v", os.Args)
	//	testloc("climb.png", afk.TowerInside)
	//	    ocrtest()
	//    d := cfg.ImageDir(".")
	//
	//    color.HiCyan("%v", d)
	//	    testRegion()
	//	return
	//    player := os.Args[1]
	player := "Devitool"
	devs, e := adb.Devices()
	devno := len(devs)
	if e != nil {
		log.Panicf("\ndevs:%v\nerr:%v", devs, e)
	}
	if len(devs) > 0 {
		var desc string = "Choose, which one will be used by bot\n"
		for i, dev := range devs {
			desc += fmt.Sprintf("%v: Serial-> %v,   id-> %v,    resolution-> %v\n", i, dev.Serial, dev.TransportId, dev.Resolution)
		}
		devno, _ = strconv.Atoi(cfg.UserInput(desc, "0"))
	} else {
		devno = 0
		d, e := adb.Connect("localhost", "5715")
		if e != nil {
			panic("dev err")
		}
		devs = append(devs, d)
	}

	gm := afk.New(game, player)
	bt := bot.New(devs[devno], gm)

	task := bot.Task{
        Name: "CampainPush",
		Actions: []string{afk.DOPUSHCAMP},
		Repeat:  10}
	//	bt.UP()
	//    bt.Tower(afk.Kings)
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
	color.HiBlue("Test overall:\n   Split   --> %v/%v\n   OneImg  --> %v/%v", improved, len(testlocs), casual, len(testlocs))

}

func testloc(img, loc string) (r1, r2 bool) {
	b := afk.New("afk", "test")
	loca := b.GetLocation(loc)
	t := ocr.ImprovedTextExtract(img)
	mt := ocr.TextExtract(img)

	color.HiMagenta("Test location: [%v], source: %v", loca.Key, img)
	color.HiCyan("\nMost Accurate: -->%v", t)
	//	ass := ocr.KeywordHits(b.GetLocation(loc).Keywords, t.Fields())
	ass := t.Intersect(loca.Keywords)
	if r1 = len(ass) >= loca.Threshold; r1 {
		color.HiGreen("Result: %v  Hits --> [%v]", r1, ass)
	} else {
		color.HiRed("Result: %v  Hits --> [%v]", r1, ass)
	}

	color.HiYellow("General: -->%v", mt)
	//	ass = ocr.KeywordHits(b.GetLocation(loc).Keywords, mt.Fields())
	ass = t.Intersect(loca.Keywords)

	if r2 = len(ass) >= loca.Threshold; r2 {
		color.HiGreen("Result: %v Hits --> [%v]", r2, ass)
	} else {
		color.HiRed("Result: %v Hits --> [%v]", r2, ass)
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
	//    ass := ocr.KeywordHits(b.GetLocation(loc).Keywords, r)
	//
	//    if r2 := ass >= b.GetLocation(loc).Threshold; r2 {
	//        color.HiGreen("Result: %v Hits --> [%v]", r2, ass)
	//    } else {
	//        color.HiRed("Result: %v Hits --> [%v]", r2, ass)
	//    }
}
