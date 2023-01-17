package cfg

import (
	"fmt"
)

type Profile struct {
	DeviceSerial string `yaml:"connection"`

	User *User `yaml:"userprofile"`
	//  Recognition settings (cmd args for 'Imagick' and 'Tesseract')
	Imagick       []string `yaml:"imagick"`
	AltImagick    []string `yaml:"alt_imagick"`
	Tesseract     []string `yaml:"tesseract"`
	AltTesseract  []string `yaml:"alt_tesseract"`
	Bluestacks    []string `yaml:"bluestacks"`
	UseAltImagick bool     `yaml:"use_alt_imagick"`
	UseAltTess    bool     `yaml:"use_alt_tess"`

	// Dict short word exceptions (>= 3)
	Exceptions []string `yaml:"dict_shrt_except"`

	Loglevel string `yaml:"loglevel"`
	DrawStep bool   `yaml:"draw_step"`
}

type SystemVars struct {
	Logfile                 string `yaml:"logfile"`
	UserConfPath            string
	parties                 []*RunableExe
	App, Userhome, Temp, Db string
}

type RunableExe struct {
	name string
	path string
}

func (ac *Profile) String() string {
	return fmt.Sprintf(
		"%v"+
			"%s\n"+
			"-> Args: \n"+
			" Bluestacks: %v\n"+
			" Magick: %v\n"+
			" Tesseract: %v\n"+
			"%v\n"+
			isStr(ac.DeviceSerial)(" -> Device: "),
		ac.User,
		ac.Bluestacks,
		ac.Imagick,
		ac.Tesseract,
	)
}

func isStr(str string) func(...interface{}) string {
	if str == "" {
		return red
	} else {
		return green
	}
}

type User struct {
	Account     string   `yaml:"account"`
	Game        string   `yaml:"game"`
	TaskConfigs []string `yaml:"taskconfigs"`
}

func (up *User) String() string {
	return fmt.Sprint(isStr(up.Game)("\n -> Game: "+up.Game+"\n ") +
		isStr(up.Account)("     Account: "+up.Account+"\n "))
}

func New(accname, game string, taskcfgpath []string) *User {
	return &User{Account: accname, Game: game, TaskConfigs: taskcfgpath}
}

type ReactiveAlto struct {
	Name      string      `yaml:"name"`
	Taptarget []Areaction `yaml:"reactions"`
}
type Areaction struct {
	If string   `yaml:"if"`
	Do []string `yaml:"do"`
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
	Keywords  []string `yaml:"keywords"`
}

func (l *Location) String() string {
	return fmt.Sprintf("Location key: %v", l.Key)
}

var defUser = &Profile{
	DeviceSerial: "",
	User: &User{
		Account:     "",
		Game:        "AFK Arena",
		TaskConfigs: []string{"cfg/reactions.yaml", "cfg/daily.yaml"},
	},
	Imagick: []string{"-colorspace", "Gray", "-alpha", "off", "-threshold, ", "75%"},
	AltImagick: []string{
		"-colorspace", "Gray",
		"-alpha", "off",
		"-threshold", "75%",
		"-edge", "2",
		"-negate",
		"-black-threshold", "90%",
	},
	Tesseract: []string{
		"--psm", "6",
		"-c", "tessedit_char_blacklist=[“€”\"’^#@™°&!~'‘|<$>«»,¢\\_;§®‘*~.°├⌐ÇöÑ{}",
		"-c", "tessedit_create_alto=1",
		"-c", "tessedit_create_txt=1",
		"quiet",
	},
	AltTesseract: []string{"--psm", "3", "hoot", "quiet"},
	Bluestacks:   []string{"--instance", "Rvc64", "--cmd", "launchApp", "--package", "com.lilithgames.hgame.gp.id"},
	Exceptions:   []string{"Go", "Up ", "In", "Tap"},
	Loglevel:     "FATAL",
	DrawStep:     false,
}
