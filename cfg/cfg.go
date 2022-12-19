package cfg

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	savefile = "cfg/save.yaml"
)

type Location struct {
	Key       string   `yaml:"name"`
	Grid      string   `yaml:"grid"`
	Threshold int      `yaml:"hits,omitempty"`
	Keywords  []string `yaml:"keywords"`
	Wait      bool     `yaml:"wait"`
	// Actions   []*Point `yaml:"actions"`
}
type Action struct {
    Name string `yaml:"name"`
    Grid string `yaml:"grid"`
    Loc string  `yaml:"startloc"`
    MidlocId string `yaml:"midloc"`
    FinlocId string `yaml:"final"`
}
// Position on Grid
func (l *Location) Position() (x, y int) {
	sx, sy, success := strings.Cut(l.Grid, ":")
	if success {
		x = toInt(sx)
		y = toInt(sy)
	}
	return
}

func Parse(s string, out interface{}) {
    pwd, _ := os.Getwd()
    log.Warnf("pwd -> %v", pwd)
	f, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, out)
	if err != nil {
		log.Fatalf("MARSHAL WASTED: %v", err)
	}
	log.Debugf("MARSHALLED: %v\n\n", out)
}

func UserInput(desc, def string) string {
	reader := bufio.NewReader(os.Stdin)
	// text := "3"
	color.HiCyan(desc)
	color.HiRed("---------------------")
	fmt.Printf("[default:%v]: ", color.HiGreenString(def))
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)

	return strings.Trim(text, "\r")
}

func toInt(s string) int {
	num, e := strconv.Atoi(s)
	if e != nil {
		fmt.Printf("\nerr:%v\nduring run:%v", e, "intconv")
	}
	return num
}
