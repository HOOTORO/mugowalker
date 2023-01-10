package ui

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"worker/cfg"
)

var (
	log                           *logrus.Logger
	red, green, cyan, yellow, mag func(...interface{}) string
	oneIdent                      = "\t"
)

func init() {
	log = cfg.Logger()
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	mag = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()
}

func UserListInput(l []string, title, defvalDesc string) int {
	termClear()
	desc := ListDesc(l, title, defvalDesc, "0")
	fmt.Print(desc)
	return intInput(len(l))
}

func UserFillSctructInput(in any, defvalDesc string) {
	//    termClear()
	//    ty := reflect.ValueOf(in).Type()
    type T struct{
        Hui string
        Pizda string
	}
    t:= T{"hui", "pizda"}
	ty := reflect.TypeOf(&in)
	tyc := reflect.TypeOf(&t)
    fmt.Printf("	!!!HUI : %v , %v\n", ty, tyc)
//	s := reflect.ValueOf(&in).Elem()
    fmt.Printf("	!!!CanSet : %v\n", reflect.ValueOf(&t).Elem().CanSet())
	fmt.Printf("	!!!CanSet : %v\n", reflect.ValueOf(&in).Elem().CanSet())

	fmt.Printf("	!!!BEFIOORE : %v\n", reflect.ValueOf(&t).Elem())
    fmt.Printf("	!!!BEFIOORE : %v\n", reflect.ValueOf(&in).Elem())

    reflect.ValueOf(&t).Elem().Field(1).SetString("AFK Arena")
	fmt.Printf("	!!!AFTERSET : %v\n", reflect.ValueOf(&t).Elem())
    reflect.ValueOf(&in).Elem().Field(1).SetString("hui")
//    reflect.ValueOf(&in).Elem().Elem().Field(1).SetString("AFK Arena")
	fmt.Printf("	!!!AFTER SOSI : %v\n", reflect.ValueOf(&in).Elem())

//fmt.Printf("	!!!ELEM : %v\n", reflect.ValueOf(s))


	//    sv := s.Elem()
	//
	//    fmt.Printf("TypeOF {SV}: %v\n", sv.Type())
	//    fmt.Printf("	CanSet : %v\n", sv.CanSet())
	//    fmt.Printf("\nKind: S  %v  V: %v  SV: %v\n\n",s.Kind(),v.Kind(), s.Kind() )
	//
	//	for i := 0; i < s.NumField(); i++ {
	//	}

	//    desc := ListDesc(l, title, defvalDesc ,"0")
	//    fmt.Print(desc)
	//	return intInput(len(l))
	//return nil
}

func ChangeVal(val string) string {
	fmt.Printf("Enter new value for [%v]:", val)
	return strInput()
}

func ListDesc(l []string, title, defDesc, defVal string) string {
	var res string

	divider := red("-----------------------------------")
	res += fmt.Sprintf("%v\n", mag(title))
	for i, v := range l {
		res += fmt.Sprintf("%v %v\n", yellow("  [", i+1, "]"), cyan(v))
	}
	res += red(" [0] ", defDesc)
	return fmt.Sprintf("%v\n%v\ndefautl[%v] --> ", res, divider, green(defVal))
}

func intInput(maxi int) int {
	reader := bufio.NewReader(os.Stdin)
	bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	r := strings.Trim(text, "\r")
	dig, e := strconv.Atoi(r)
	if e != nil || dig > maxi {
		return 0
	}
	return dig
}
func strInput() string {
	reader := bufio.NewReader(os.Stdin)
	bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	return strings.Trim(text, "\r")
}

func termClear() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Errorf("termclr error: %v", err)
	}
}

//if len(text) == 0 {
//    text = def
//}
