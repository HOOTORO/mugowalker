package bot

import (
	"fmt"
	"os"

	// "worker/adb"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Cirrus interface {
	parse(string)
}

type Mission struct {
	Goal  string          `json:"goal,omitempty"`
	Entry string          `json:"entry,omitempty"`
	Plan  map[string]bool `json:"plan,omitempty"`
}

type Location struct {
	// Label    string
	Keywords []string `json:"keywords,omitempty"`
	Actions  []string `json:"actions,omitempty"`
}

// type SupaLocation struct {
// 	// Label    string
// 	Keywords []string          `json:"keywords,omitempty"`
// 	Actions  map[string]Action `json:"actions,omitempty"`
// }

// type Action struct {
// 	*adb.Point
// 	ActionProperties
// 	BaseDelay int
// }

type ActionProperties struct {
	// *Action
	Check  bool
	Delay  int
	Repeat int
}

type Task struct {
	Begin   string // Location key
	Actions map[string]ActionProperties
}

const (
	loccnf  = "../vscode/afkarena/worker/bot/cfg/locations.yaml"
	actcnf  = "../vscode/afkarena/worker/bot/cfg/actions.yaml"
	newconf = "../vscode/afkarena/worker/bot/cfg/newconf.yaml"
)

// var (
// 	scenaries []Mission
// 	actions   map[string]adb.Point
// 	locations []Location
// )

func (d *Daywalker) Load(p string) []Mission {
	// daily := make([]interface{}, 0, 1)
	parse(p, &d.job)
	// parse(newconf, &d.supaloc)
	parse(loccnf, &d.locs)
	parse(actcnf, &d.actions)
	color.HiYellow("GOOD NEWS, EVERY ONE! \n we have a job to do :>\n%v", d.job)
	fmt.Printf("---\n\n%v\n\n", d)

	return d.job
}

func parse(s string, out interface{}) {
	// Load the file; returns []byte
	f, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	// color.HiCyan("Try to load: %v", reflect.TypeOf(out))

	err = yaml.Unmarshal(f, out)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// color.HiRed("MARSHALLED: %v\n\n", out)
	log.Debugf("MARSHALLED: %v\n\n", out)
}

func (l *Location) parse(s string) {
	parse(s, l)
}

func (m *Mission) parse(s string) {
	parse(s, m)
}

// func (d *Daywalker) save(in Cirrus) error {
// 	// p := fmt.Sprintf("%T", in)

// 	res := make(map[string]SupaLocation)
// 	for k, v := range d.locs {
// 		actions := make(map[string]Action)
// 		for _, actionname := range v.Actions {
// 			pnt := d.actions[actionname]
// 			actions[actionname] = Action{Point: &pnt, BaseDelay: 1}
// 		}
// 		res[k] = SupaLocation{Keywords: v.Keywords, Actions: actions}

// 	}
// 	f, err := yaml.Marshal(res)
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}
// 	fmt.Printf("--- t dump:\n%s\n\n", string(f))
// 	err = os.WriteFile(newconf, f, os.ModeDevice)
// 	return err
// }
