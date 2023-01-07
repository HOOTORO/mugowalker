package main

import (
	"fmt"
	"image"
	"os"
	"strings"
	"worker/ui"

	"worker/afk"
//	"worker/bot"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"

	"golang.org/x/exp/slices"
)

func main() {
	// r := cfg.List()
	//	r := cfg.SimpleMenu()
	//	fmt.Printf("VAL BIII --> %v", r.Myval())
	//	return
	if len(os.Args) > 1 && os.Args[1] == "-t" {

		color.HiRed("%v", "TEST RUN")
		ocrtest()
		return
	}

	app := &cfg.AppConfig{
		DeviceSerial: "192.168.1.7:5555",
		UserProfile: &cfg.UserProfile{
			Account:     "E6osh!ro",
			Game:        "AFK Arena",
            TaskConfigs: []string{"cfg/reactions.yaml", "cfg/daily.yaml"},
		},
		Imagick: cfg.OcrConf.Imagick,
		AltImagick: []string{"-colorspace", "Gray", "-alpha", "off", "-threshold", "75%", "-edge", "2", "-negate", "-black-threshold",
			//			"-white-threshold",
			//			"60%",
			"90%",
		},
		Tesseract:    cfg.OcrConf.Tesseract,
		AltTesseract: []string{"--psm", "3", "hoot", "quiet"},
		Bluestacks:   []string{"--instance", "Rvc64_16", "--cmd", "launchApp", "--package", "com.lilithgames.hgame.gp.id"},
		Exceptions:   cfg.OcrConf.Exceptions,
		Loglevel:     "INFO",
		DrawStep:     false,
		Folders: struct {
			Logfile     string `yaml:"logfile"`
			RootDir     string `yaml:"rootDir"`
			TempImgDir  string `yaml:"tempImgDir"`
			SqDBDir     string `yaml:"sqDBDir"`
			UserDir     string `yaml:"userDir"`
			GameConfDir string `yaml:"gameConfDir"`
			TestDataDir string `yaml:"testDataDir"`
		}{
			Logfile:     "app.log",
            RootDir:     ".afk_data",
            TempImgDir:  "work_images",
            SqDBDir:     "db",
            UserDir:     "usrdata",
            GameConfDir: "cfg",
			TestDataDir: "_test",
		},
	}


	cfg.Save("runset.yaml", app)
return
	model := ui.SimpleMenu()

	fmt.Print("Chosen one! %v", model.Myval())
	//    e := ui.MainMenu()
	return

	//    if e!=nil {
	//        return
	//	}

//	rTaskConf := []string{"cfg/reactions.yaml", "cfg/daily.yaml"}
//	user := cfg.User(player, game, connect, rTaskConf)

//	device := cfg.Load(user)
//	gm := afk.New(user)
//	bt := bot.New(device, gm)
	cfg.TermClear()
UserMenu:
	choice := cfg.UserInput("What bot should do?\n0. Run all (Default)\n1. Run daily?\n2. Push campain?\n3-7. Push towers (KT,L,M,W,G)\n9. OCR Settings", "0")

	if choice == "9" {
		cfg.TermClear()

		choice = cfg.UserInput(fmt.Sprintf("Current setup:\n %v\n Change threshold?\n", cfg.OcrConf), "75%")
		if choice != "0" {
			color.HiRed("Definetly do something...")
			for i, s := range cfg.OcrConf.Imagick {
				if s == "75%" {
					cfg.OcrConf.Imagick[i] = fmt.Sprintf("%v%", choice)
					cfg.TermClear()

					color.HiGreen("New ocr settings --> %v", cfg.OcrConf.Imagick)
					goto UserMenu
				}
			}
		}
	}

//	switch choice {
//	case "0":
//		bt.UpAll()
//	case "1":
//		bt.Daily()
//	case "2":
//		push := bt.Task(afk.DOPUSHCAMP)
//		bt.React(push)
//	case "3":
//		kt := bt.Task(afk.Kings)
//		bt.React(kt)
//	case "4":
//		kt := bt.Task(afk.Light)
//		bt.React(kt)
//	case "5":
//		kt := bt.Task(afk.Mauler)
//		bt.React(kt)
//	case "6":
//		kt := bt.Task(afk.Wilder)
//		bt.React(kt)
//	case "7":
//		kt := bt.Task(afk.Graveborn)
//		bt.React(kt)
//	default:
//		color.HiRed("DATS WRONG NUMBA MAFAKA!")
//	}
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
				//                toadd := v.CONTENT
				//                if len(toadd)<width{
				//                    toadd += "  "+ strings.Repeat("-", width - len(toadd))
				//				} else {
				//                    toadd = toadd[:width-3]+"..."
				//				}
				//
				//                str += fmt.Sprintf("%v->	%v |", pass("%v", toadd), fail("%vx%v", v.HPOS, v.VPOS))
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
