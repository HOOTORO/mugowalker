package main

import (
	"fmt"

	"github.com/fatih/color"

	"golang.org/x/exp/slices"
)

func main() {
	const (
		name = "Bluestacks"
		host = "127.0.0.1"
		port = "5615"
	)
	// TODO: scaling  adb shell wm size returns resolution
	log.SetLevel(log.InfoLevel)

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
