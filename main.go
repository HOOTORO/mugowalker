package main

import (
	"context"
	"fmt"
	"image"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"golang.org/x/sys/windows"

	"worker/adb"
	"worker/afk/activities"
	"worker/bot"
	"worker/ui"

	"github.com/sirupsen/logrus"

	"worker/afk"
	"worker/cfg"
	"worker/ocr"

	"github.com/fatih/color"

	"github.com/erikgeiser/coninput"
	"golang.org/x/exp/slices"
)

var (
	log                        *logrus.Logger
	user                       *cfg.Profile
	red, green, cyan, ylw, mgt func(...interface{}) string
)

func init() {
	user = cfg.ActiveUser()
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	ylw = color.New(color.FgHiYellow).SprintFunc()
	mgt = color.New(color.FgHiMagenta).SprintFunc()
}

func main() {
	log = cfg.Logger()
	fn := func(a string, b string) {
		log.Warnf("%v |>\n %v", mgt(a), b)
	}

	if len(os.Args) > 1 && os.Args[1] == "-t" {
		log.SetLevel(logrus.InfoLevel)

		color.HiRed("%v", "TEST RUN")

		d, e := adb.Connect("127.0.0.1:5555")
		if e != nil {
			// log.Errorf(red("%v"), e)
			// blueprc := cfg.RunProc(cfg.BluestacksExe, user.Bluestacks.Args()...)
			// log.Infof("Bluestacks started: %v, args: %v", blueprc.Process.Pid, user.Bluestacks)
		}
		gw := afk.New(user.User)
		bb := bot.New(d, fn)
		bot := afk.NewArenaBot(bb, gw)

		activities.Push(bot, fn)

		// bot.AltoRun("quests", fn)
		return
	}

	// log.SetLevel(logrus.TraceLevel)
	log.Warnf(red("RUN BEGIN : %v"), time.Now())

	conf := ui.CfgDto(user)

	err := ui.RunMainMenu(conf)
	if err != nil {
		log.Errorf("ERROROR: %v", err)
	}
}

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

func ocrtest() {
	b := afk.New(&cfg.User{Account: "test", Game: "afk", TaskConfigs: []string{"cfg/reactions.yaml"}})

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
	// 		if len(v.CONTENT) > 3 || slices.Contains(user.Exceptions, v.CONTENT) {
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

// startkill example
func runcnd() {

	// hoho, e := cfg.StartProc("HD-Player", cfg.ActiveUser().Bluestacks...)
	cmd := exec.Command("HD-Player", cfg.ActiveUser().Bluestacks.Args()...)
	log.Debugf("run cmd: %v\n", cmd.String())
	cmd.Stdout = os.Stdout
	e1 := cmd.Start()

	e2 := cmd.Wait()
	log.Errorf("wearefinish %v", e1, e2)
	// hoho := 11348
	// if e != nil {
	// 	log.Warn("RAUCREATEDN", e, hoho)
	// }
	// log.Warn("RAUCREATEDN", e, hoho)
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	p, _ := os.FindProcess(hoho)
	// 	e := p.Signal(syscall.SIGKILL)
	// 	log.Warn("kill ee", e)

	// 	log.Warnf("PROCESS --> %+v", p)
	// }()
	// e = p.Kill()
}
