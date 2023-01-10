package main

import (
	"fmt"
	"image"
	"os"
	"strings"
	"time"
	"worker/ui"

	"worker/afk"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"

	"golang.org/x/exp/slices"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-t" {
		color.HiRed("%v", "TEST RUN")
        user := cfg.UserProfile{Account: "ss", Game:"aa"}
        ui.UserFillSctructInput(user, "")
		//		ocrtest()
		return
	}

	//    model := ui.MainMenu()
	//
	//    fmt.Print("Chosen one! %v", model)

	//        menu9 := []string{"Current setup",strings.Join(cfg.Env.Imagick,""),"Change threshold?"}

	//	rTaskConf := []string{"cfg/reactions.yaml", "cfg/daily.yaml"}
	//	user := cfg.User("", game, connect, rTaskConf)
	//	device := cfg.Load(user)
	//	device, _ := adb.Devices()
	//	gm := afk.New(cfg.Env.UserProfile)
	//	bt := bot.New(device[0], gm)

	if cfg.Env.UserProfile != nil {
		fmt.Print(cfg.Env.UserProfile)
	} else {
		ui.UserFillSctructInput(cfg.Env.UserProfile, "")
	}
MainMenu:

	choice := ui.UserListInput(mainmenu, "AFK Bot\n What bot should do?", "Exit")
	switch choice {

	case 4:
	Towers:
		choice = ui.UserListInput(tower, "Which one?", "Back")
		switch {
		case choice > 0:
			color.HiYellow("Climbing... %v", tower[choice-1])
		case choice == 0:
			goto MainMenu
		default:
			color.HiRed("DATS WRONG TOWAH MAFAKA!")
		}
		time.Sleep(3 * time.Second)
		goto Towers

		//		push := bt.Task(afk.DOPUSHCAMP)
		//		bt.React(push)
		//	case 4:
		//		kt := bt.Task(afk.Kings)
		//		bt.React(kt)
		//	case 5:
		//		kt := bt.Task(afk.Light)
		//		bt.React(kt)
		//	case 6:
		//		kt := bt.Task(afk.Mauler)
		//		bt.React(kt)
		//	case 7:
		//		kt := bt.Task(afk.Wilder)
		//		bt.React(kt)
		//	case 8:
		//		kt := bt.Task(afk.Graveborn)
		//		bt.React(kt)
	case 5:
	Nine:
		choice = ui.UserListInput(cfg.Env.Imagick, "Current setup", "Back")
		switch {
		case choice > 0:
			cfg.Env.Imagick[choice-1] = ui.ChangeVal(cfg.Env.Imagick[choice-1])
			color.HiBlue("dosomething")
			time.Sleep(2 * time.Second)
			goto Nine
		default:
			goto MainMenu
		}
	case 0:
		os.Exit(0)
	default:
		color.HiRed("DATS WRONG NUMBA MAFAKA!")
		time.Sleep(2 * time.Second)
		goto MainMenu
	}
}

func ocrtest() {
	b := afk.New(&cfg.UserProfile{Account: "test", Game: "afk", TaskConfigs: []string{"cfg/reactions.yaml"}})

	ocrparams := &cfg.OcrConfig{
		Split: nil,
		Imagick: []string{
			"-colorspace", "Gray", "-alpha", "off",
			"-threshold",
			"75%",
			"-edge",
			"2",
			"-negate",
			//            "-canny",
			//            "0x1+10%+30%",
			//            "-unsharp",
			//            "1x1",
			//            "-blur",
			//            "0x1",

			"-black-threshold",
			//			"-white-threshold",
			//			"60%",
			"90%",
			//			"-bordercolor", "black", "-border", "3x3",
			//            "-negate",
		},
		Tesseract: []string{
			//            "--tessdata-dir",
			//			"C:\\Program Files\\Tesseract-OCR\\tessdata\\frmgit\\tessdata_fast",
			"--psm", "3",
			"hoot", "quiet",
		},
		Exceptions: cfg.OcrConf.Exceptions,
	}
	cfg.OcrConf = ocrparams

	testdata := func(lo uint, im string) *struct {
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
		//		testdata(afk.Campain.Id(), "_test/camp.png"),
		//		testdata(afk.RANHORNY.Id(), "_test/h.png"),
		//		testdata(afk.QUESTS.Id(), "_test/quests.png"),
		testdata(afk.QUESTS.Id(), "_test/quests2.png"),
		//		testdata(afk.GUILDCHEST.Id(), "_test/gichest.png"),
		testdata(afk.WIN.Id(), "_test/win_hard.png"),
		//        testdata( afk.Campain, "test/cpn1.png"),
		//        testdata(afk.DARKFORREST, "test/forrest.png"),
		//        testdata( afk.Campain, "test/cpn2.png"),
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

	overall := 0
	for _, v := range testlocs {
		res := testloc(v.img, b.GetLocation(v.loc))
		if res {
			overall++
		}

	}
	//    testRegion("test/btl_onestg_1.png")
	//    testRegion("test/btl_multstg_1.png")
	color.HiBlue("\n	######################\n	#  Test overall:     #\n	#   Basic   --> %v/%v  #\n	######################", overall, len(testlocs))
}

func testloc(img string, loc *cfg.Location) (r1 bool) {
	fail := color.New(color.FgHiRed, color.Bold).SprintfFunc()
	pass := color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	regular := color.New(color.FgHiYellow).SprintFunc()
	divider := color.HiMagentaString("#################################################################")
	regular("%v - %v", img, loc)

	restr := "\n	> /Result	-> %v\n	> /Hits		-> [%v/%v]\n          ______________________________\n\n"

	fmt.Printf("\n%s%s\n		#  |RUN|  #  	--> [TEST] location: [%v], source: %v	#\n%s%s\n\n\n", divider, divider, fail(loc.Key), fail(img), divider, divider)

	mt := ocr.TextExtract(img)
	ass := mt.Intersect(loc.Keywords)

	fmt.Print(regular("\n			<----------- /General/ -------------------> \n\n	"), highlight(mt.Fields(), ass, pass))

	if r1 = len(ass) >= loc.Threshold; r1 {
		fmt.Print(pass(restr, r1, len(ass), loc.Threshold))
	} else {
		fmt.Print(fail(restr, r1, len(ass), loc.Threshold))
	}
	pass("xu")
	al := ocr.TextExtractAlto(img)
	//    fmt.Printf("%v", pass("%v",))
	fmt.Print(regular("\n			<----------- /Alto/ ----------------------> \n\n"))
	tl := al.Layout.Page.PrintSpace.ComposedBlock.TextBlock.TextLine
	//	width := 30
	for _, line := range tl {
		str := fail("	> ")
		for _, v := range line.String {
			if len(v.CONTENT) > 3 || slices.Contains(cfg.OcrConf.Exceptions, v.CONTENT) {

				str += fmt.Sprintf("%s	->	%s | ", pass("%-12s", cutlong(v.CONTENT, 10)), fail("%sx%-4s", v.HPOS, v.VPOS))
			}
		}
		if fail("	> ") != strings.TrimLeft(str, " ") {
			fmt.Printf("%v\n", str)
		}

	}
	fmt.Print("#\n#\n" + divider + "\n\n\n")
	return
}

func cutlong(s string, l int) string {
	if len(s) > l {
		return s[:l-3] + "..."
	}
	return s
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
