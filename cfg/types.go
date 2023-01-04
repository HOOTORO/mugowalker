package cfg

import "fmt"

type Action struct {
	Name        string   `yaml:"name"`
	Start       []string `yaml:"startloc,omitempty"`
	Steps       []Step   `yaml:"steps"`
	Conditional []Step   `yaml:"conditional,omitempty"`
	Next        string   `yaml:"nexta,omitempty"`
}
type Step struct {
	Grid        string   `yaml:"grid"`
	Delay       int      `yaml:"delay,omitempty"`
	Skiplocheck bool     `yaml:"skiplocheck"`
	Wait        bool     `yaml:"wait,omitempty"`
	Loc         []string `yaml:"loc,omitempty"`
}

type ReactiveTask struct {
	Name      string     `yaml:"name"`
	Limit     int        `yaml:"limit"`
	Criteria  string     `yaml:"criteria`
	Avail     string     `yaml:"avail"`
	Reactions []Reaction `yaml:"reactions"`
}

type Reaction struct {
	If     string `yaml:"if"`
	Before string `yaml:"before"`
	Do     string `yaml:"do"`
	After  string `yaml:"after"`
}

type Task struct {
	Name    string   `yaml:"name"`
	Actions []string `yaml:"actions"`
	Repeat  int      `yaml:"repeat,omitempty"`
}

type ocrConfig struct {
	Split     []string `yaml:"split"`
	Imagick   []string `yaml:"imagick"`
	Tesseract []string `yaml:"tesseract"`
}

type emuConf []struct {
	Cmd  string   `yaml:"cmd"`
	Args []string `yaml:"args"`
}
type Location struct {
	Key       string   `yaml:"name"`
	Grid      string   `yaml:"grid"`
	Threshold int      `yaml:"hits,omitempty"`
	Keywords  []string `yaml:"keywords"`
	Wait      bool     `yaml:"wait"`
	// Actions   []*Point `yaml:"actions"`
}

func (emu emuConf) Command(name string) []string {
	for _, v := range emu {
		if v.Cmd == name {
			r := []string{v.Cmd}
			r = append(r, v.Args...)
			return r
		}
	}
	return []string{}
}

type UserProfile struct {
	Account       string
	Game          string
    TaskConfigs   []string
	ConnectionStr string
}

func (up *UserProfile) String() string {
    return fmt.Sprintf("\n--> Game: %v\n--> Acc: %v\n", up.Game, up.Account)
}
func User(accname, game, connect string, taskcfgpath []string) *UserProfile {
	return &UserProfile{Account: accname, Game: game, TaskConfigs: taskcfgpath, ConnectionStr: connect}
}
