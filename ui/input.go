package ui

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"worker/cfg"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	log                           *logrus.Logger
	red, green, cyan, yellow, mag func(...interface{}) string
)

func init() {
	log = cfg.Logger()
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	mag = color.New(color.FgHiMagenta, color.BgHiWhite).SprintFunc()
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

// func CfgDto() map[string]string {
// 	dto := make(map[string]string, 0)
// 	conf := cfg.Env
// }
