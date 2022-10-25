package bot

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	newconf = "../vscode/afkarena/worker/bot/cfg/config.yaml"
	save    = "../vscode/afkarena/worker/bot/cfg/save.yaml"
)

func init() {
	locs = make(map[string]Location)
	parse(newconf, locs)
	color.HiYellow("Loaded config... \n%v", locs)
}

func (d *Daywalker) Load(p string) []Task {
	parse(taskfile, &d.Tasks)
	if !IsValid(d.Tasks, locs) {
		color.HiRed("Invalid Data")
	}
	return d.Tasks
}

func parse(s string, out interface{}) {
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

func (d *Daywalker) save(in interface{}) error {
	lastLoc := d.loc
	f, err := yaml.Marshal(lastLoc)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("---:\n%s\n\n", string(f))
	err = os.WriteFile(save, f, os.ModeDevice)
	return err
}

func IsValid(mission []Task, locations map[string]Location) bool {
	fl, badentities := false, "Bad entries:\n"

	for _, task := range mission {
		val, ok := locations[task.Entry]
		if !ok {
			fl = ok
			badentities += fmt.Sprintf("%v - doesn't exist.\n", task.Entry)
		}

		val, ok = locations[task.Exit]
		if !ok {
			fl = ok
			badentities += fmt.Sprintf("%v - doesn't exist.\n", task.Exit)
		}
		_ = val
	}
	// TO DO: validate consistnecy
	// if !ok {
	// 	return errors.New(
	// 		fmt.Sprintf("Action<%v> does not exist in <%v>", actionName, d.Location.Label),
	// 	)
	// }
	color.HiRed("MISSION Validatio.\n%v \n\nValid? %v", badentities, !fl)
	return fl
}