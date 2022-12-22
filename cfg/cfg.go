package cfg

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"worker/adb"

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
	Name     string `yaml:"name"`
	Grid     string `yaml:"grid"`
	Delay    int    `yaml:"delay"`
	Loc      string `yaml:"startloc"`
	MidlocId string `yaml:"midloc"`
	FinlocId string `yaml:"final"`
}

// Position on Grid
func (l *Location) Position() *adb.Point {
	return cutgrid(l.Grid)
}

func (act *Action) StartXY() *adb.Point {
	return cutgrid(act.Grid)
}

func (act *Action) OverlayGrids() (taps []*adb.Point) {
	if strings.Contains(act.MidlocId, "overlay") {
		grids := strings.Split(strings.Trim(act.MidlocId, "overlay "), ";")
		for _, v := range grids {
			taps = append(taps, cutgrid(v))
		}
	}
	return taps
}

func Parse(s string, out interface{}) {
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
	if len(text) == 0 {
		text = def
	}
	return strings.Trim(text, "\r")
}

func toInt(s string) int {
	num, e := strconv.Atoi(s)
	if e != nil {
		fmt.Printf("\nerr:%v\nduring run:%v", e, "intconv")
	}
	return num
}

func cutgrid(str string) (p *adb.Point) {
	ords := strings.Split(str, ":")
		p = &adb.Point{
			X: toInt(ords[0]),
			Y: toInt(ords[1]),
            Offset: 1,
		}
	if len(ords)>2 {
        p.Offset = toInt(ords[2])
	}
	return
}
