package main

import (
	"context"
	"fmt"
	"image"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sys/windows"

	"worker/ui"

	"github.com/sirupsen/logrus"

	"worker/afk"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"

	"github.com/erikgeiser/coninput"
	"golang.org/x/exp/slices"
)

var log *logrus.Logger

func run() (err error) {
	con, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		return fmt.Errorf("get stdin handle: %w", err)
	}

	var originalConsoleMode uint32

	err = windows.GetConsoleMode(con, &originalConsoleMode)
	if err != nil {
		return fmt.Errorf("get console mode: %w", err)
	}

	fmt.Println(
		"Input mode:",
		coninput.DescribeInputMode(originalConsoleMode),
	)

	newConsoleMode := coninput.AddInputModes(
		windows.ENABLE_MOUSE_INPUT,
		windows.ENABLE_WINDOW_INPUT,
		windows.ENABLE_PROCESSED_INPUT,
		windows.ENABLE_EXTENDED_FLAGS,
	)

	fmt.Println(
		"Setting mode to:",
		coninput.DescribeInputMode(newConsoleMode),
	)

	err = windows.SetConsoleMode(con, newConsoleMode)
	if err != nil {
		return fmt.Errorf("set console mode: %w", err)
	}

	defer func() {
		fmt.Println("Resetting input mode to:", coninput.DescribeInputMode(originalConsoleMode))

		resetErr := windows.SetConsoleMode(con, originalConsoleMode)
		if err == nil && resetErr != nil {
			err = fmt.Errorf("reset console mode: %w", resetErr)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	for {
		if ctx.Err() != nil {
			break
		}

		events, err := coninput.ReadNConsoleInputs(con, 16)
		if err != nil {
			return fmt.Errorf("read input events: %w", err)
		}

		fmt.Printf("Read %d events:\n", len(events))
		for _, event := range events {
			fmt.Println("  ", event)
		}
	}

	return nil
}

func testWinEvents() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)

		os.Exit(1)
	}
}

func main() {
	log = cfg.Logger()
	if len(os.Args) > 1 && os.Args[1] == "-t" {
		color.HiRed("%v", "TEST RUN")
		// user := cfg.UserProfile{Account: "ss", Game: "aa"}
		ocrtest()
		return
	}

	// log.Error("knock knock")
	// ocrtest()
	log.Warnf(color.RedString("RUN BEGIN : %v"), time.Now())
	testselect()
	// dev, _ := adb.Connect(cfg.Env.DeviceSerial)
	// gm := afk.New(cfg.Env.UserProfile)
	// b := bot.New(dev, gm)
	// t := b.Task(afk.DOPUSHCAMP)
	// b.React(t)
	// testWinEvents()
}

// func RunBot(choice string, confg *cfg.AppConfig) {
// 	rTaskConf := []string{"assets/reactions.yaml", "assets/daily.yaml"}

// 	gm := afk.New(confg.UserProfile)
// 	bt := bot.New(, gm)

// 	switch choice {
// 	case "0":
// 		bt.UpAll()
// 	case "1":
// 		bt.Daily()
// 	case "2":
// 		push := bt.Task(afk.DOPUSHCAMP)
// 		bt.React(push)
// 	case "3":
// 		kt := bt.Task(afk.Kings)
// 		bt.React(kt)
// 	case "4":
// 		kt := bt.Task(afk.Light)
// 		bt.React(kt)
// 	case "5":
// 		kt := bt.Task(afk.Mauler)
// 		bt.React(kt)
// 	case "6":
// 		kt := bt.Task(afk.Wilder)
// 		bt.React(kt)
// 	case "7":
// 		kt := bt.Task(afk.Graveborn)
// 		bt.React(kt)
// 	}
// }

func testselect() {
	conf := ui.CfgDto(cfg.Env)

	err := ui.RunMainMenu(conf)
	if err != nil {
		log.Errorf("ERROROR: %v", err)
	}

	//	ui.SoloStrInput()
}

func ocrtest() {
	b := afk.New(&cfg.UserProfile{Account: "test", Game: "afk", TaskConfigs: []string{"cfg/reactions.yaml"}})

	cfg.Env.Imagick = []string{
		//		to try convert test.tif -fill black -fuzz 30% +opaque "#FFFFFF" result.tif
		//		convert test.tif -brightness-contrast -40x10 -units pixelsperinch -density 300 -negate -noise 10 -threshold 70% result.tif
		//		convert test.tif -fill black -fuzz 30% +opaque "#FFFFFF" result.tif
		// convert test.tif -negate -threshold 100 -negate result.tif
		// textcleaner -g -e normalize -f 30 -o 12 -s 2 http://i.stack.imgur.com/ficx7.jpg out.png
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
	}

	cfg.Env.Tesseract = []string{
		//            "--tessdata-dir",
		//			"C:\\Program Files\\Tesseract-OCR\\tessdata\\frmgit\\tessdata_fast",
		"--psm", "3",
		"hoot", "quiet",
	}

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
	fmt.Print(regular("\n			<----------- /Alto/ ----------------------> \n\n\t"))
	fmt.Printf("%v", al)
	// tl := al.Layout.Page.PrintSpace.ComposedBlock.TextBlock.TextLine
	// //	width := 30
	// for _, line := range tl {
	// 	str := fail("	> ")
	// 	for _, v := range line.String {
	// 		if len(v.CONTENT) > 3 || slices.Contains(cfg.Env.Exceptions, v.CONTENT) {
	// 			str += fmt.Sprintf("%s	->	%s | ", pass("%-12s", cutlong(v.CONTENT, 10)), fail("%sx%-4s", v.HPOS, v.VPOS))
	// 		}
	// 	}
	// 	if fail("	> ") != strings.TrimLeft(str, " ") {
	// 		fmt.Printf("%v\n", str)
	// 	}

	// }
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
