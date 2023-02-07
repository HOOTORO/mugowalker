package cfg

import (
	"fmt"
)

type AppUser interface {
	Loglevel() string
	Game() string
	Acccount() string
	DevicePath() string
}

type OcrConfig interface {
	ImagickCfg() []string
	TesseractCfg() []string
}

type Profile struct {
	DeviceSerial string `yaml:"connection"`

	User *User `yaml:"userprofile"`
	//  Recognition settings (cmd args for 'Imagick' and 'Tesseract')
	Imagick       map[int]CmdArgs `yaml:"imagick"`
	AltImagick    []string        `yaml:"alt_imagick"`
	Tesseract     map[int]CmdArgs `yaml:"tesseract"`
	AltTesseract  []string        `yaml:"alt_tesseract"`
	Bluestacks    *Bluestacks     `yaml:"bluestacks"`
	UseAltImagick bool            `yaml:"use_alt_imagick"`
	UseAltTess    bool            `yaml:"use_alt_tess"`

	// Dict short word exceptions (>= 3)
	Exceptions []string `yaml:"dict_shrt_except"`

	Loglevel string `yaml:"loglevel"`
	DrawStep bool   `yaml:"draw_step"`
}

type CmdArgs struct {
	Key string `yaml:"key"`
	Val string `yaml:"val"`
}

type SystemVars struct {
	Logfile                 string `yaml:"logfile"`
	UserConfPath            string
	App, Userhome, Temp, Db string
}

func (ac *Profile) String() string {
	return fmt.Sprintf(
		" -> Device: %v"+
			"%s\n"+
			"-> Args: \n"+
			" Bluestacks: %v\n"+
			" Magick: %v\n"+
			" Tesseract: %v\n"+
			"%v\n"+
			Green(ac.DeviceSerial),
		ac.User,
		ac.Bluestacks,
		ac.Imagick,
		ac.Tesseract,
	)
}

// User profile
type User struct {
	Account     string   `yaml:"account"`
	Game        string   `yaml:"game"`
	TaskConfigs []string `yaml:"taskconfigs"`
}

func (up *User) String() string {
	return F("\n	-> Game: %v\n\t   Account: %v", Green(up.Game), Green(up.Account))
}

// New user profile
func New(accname, game string, taskcfgpath []string) *User {
	return &User{Account: accname, Game: game, TaskConfigs: taskcfgpath}
}

// Bluestacks vm settings
type Bluestacks struct {
	Instance string `yaml:"instance"`
	Package  string `yaml:"package"`
}

func (bs *Bluestacks) String() string {
	return F("\n -> VM: %v\n     App: %v", bs.Instance, bs.Package)
}

// Args upack in same order as packed? Will see
func (bs *Bluestacks) Args() []string {
	return []string{"--instance", bs.Instance, "--cmd", "launchApp", "--package", bs.Package}
}

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
	return F("Key: %v | hitwords: %v", Green(l.Key), Cyan(l.Keywords))
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

func (p *Profile) ImagickCfg() (args []string) {

	for _, v := range p.Imagick {
		args = append(args, v.Key, v.Val)
	}
	return
}
func (p *Profile) TesseractCfg() (args []string) {
	for _, v := range p.Tesseract {
		args = append(args, v.Key, v.Val)
	}
	return
}
