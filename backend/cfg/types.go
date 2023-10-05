package cfg

import (
	"fmt"
)

type AppUser interface {
	Account() string
	Game() string
	DevicePath() string
}

func (p *Profile) Account() string {
	return p.GameAccount
}

func (p *Profile) Game() string {
	return afkapp
}

func (p *Profile) DevicePath() string {
	return p.DeviceSerial
}

type OcrConfig interface {
	ImagickCfg() []string
	TesseractCfg() []string
}

type Profile struct {
	DeviceSerial string `yaml:"connection"`

	GameAccount string `yaml:"account"`
}

type CmdArgs struct {
	Key string `yaml:"key"`
	Val string `yaml:"val"`
}

type SystemVars struct {
	Logfile       string `yaml:"logfile"`
	UserConfPath  string
	App, Temp, Db string
}

func (ac *Profile) String() string {
	return fmt.Sprintf(
		"\nDevice |> %v"+
			"\n%s",
		ac.DeviceSerial,
		ac.Account(),
	)
}

// // Args upack in same order as packed? Will see
// func (bs *Bluestacks) Args() []string {
// 	return []string{"--instance", bs.Instance, "--cmd", "launchApp", "--package", afkapp}
// }

type ReactiveTask struct {
	Name      string     `yaml:"name"`
	Limit     int        `yaml:"limit"`
	Criteria  string     `yaml:"criteria"`
	Avail     string     `yaml:"avail"`
	Reactions []Reaction `yaml:"reactions"`
}

type Reaction struct {
	If     string `yaml:"if"`
	Before string `yaml:"before"`
	Do     string `yaml:"do"`
	After  string `yaml:"after"`
}

type Location struct {
	Key       string   `yaml:"name"`
	Threshold int      `yaml:"hits"`
	HitKw     []string `yaml:"keywords"`
}

func (l *Location) String() string {
	return fmt.Sprintf("Key: %v | hitwords: %v", l.Key, l.Keywords)
}
func (l *Location) Id() string {
	return l.Key
}

func (l *Location) Keywords() []string {
	return l.HitKw
}

func (l *Location) HitThreshold() int {
	return l.Threshold
}
