package cfg

import (
	"fmt"
	"strings"
)

type AppConfig struct {
	DeviceSerial string `yaml:"connection_str"`

	UserProfile *UserProfile
	//  Recognition settings (cmd args for 'Imagick' and 'Tesseract')
	//  Split     []string `yaml:"split"`
	Imagick       []string `yaml:"imagick"`
	AltImagick    []string `yaml:"alt_imagick"`
	UseAltImagick bool     `yaml:"use_alt_imagick"`
	Tesseract     []string `yaml:"tesseract"`
	AltTesseract  []string `yaml:"alt_tesseract"`
	UseAltTess    bool     `yaml:"use_alt_tess"`
	Bluestacks    []string `yaml:"bluestacks"`

	// Dict short word exceptions (>= 3)
	Exceptions []string `yaml:"dict_shrt_except"`

	Logfile  string `yaml:"logfile"`
	Loglevel string `yaml:"loglevel"`
	DrawStep bool   `yaml:"draw_step"`

	Dirs struct {
		Root    string `yaml:"root"`
		TempImg string `yaml:"tempImg"`
		SqDB    string `yaml:"sqDB"`
		//		User     string `yaml:"user"`
		GameConf string `yaml:"gameConf"`
		TestData string `yaml:"testData"`
	}
	RequiredInstalledSoftware []string `yaml:"required_installed_software"`
	Thiscfg                   string   `yaml:"thiscfg"`
}

func (ac *AppConfig) String() string {
	reqsoft := " -> Required software..."
	for _, v := range ac.RequiredInstalledSoftware {
		if strings.Contains(v, adbp) {
			reqsoft += "\n	ADB: " + v
		}
		if strings.Contains(v, magic) {
			reqsoft += "\n	IMAGICK: " + v
		}
		if strings.Contains(v, tesseract) {
			reqsoft += "\n	TESSERACT: " + v
		}
		if strings.Contains(v, bluestacks) {
			reqsoft += "\n	BLUESTACKS: " + v
		}
	}
	return fmt.Sprintf(
		" -> Device: %v	"+
			"%s\n"+
			" -> Args: \n"+
			"	Bluestacks: %v\n"+
			"	Magick: %v\n"+
			"	Tesseract: %v\n"+
			"%v\n"+
			" -> Config: %v",
		ac.DeviceSerial,
		ac.UserProfile,
		ac.Bluestacks,
		ac.Imagick,
		ac.Tesseract,
		reqsoft,
		ac.Thiscfg)
}

type UserProfile struct {
	Account     string
	Game        string
	TaskConfigs []string
}

func (up *UserProfile) String() string {
	return fmt.Sprintf("\n --> Game: %v\n"+
		"     Account: %v\n", up.Game, up.Account)
}

func User(accname, game string, taskcfgpath []string) *UserProfile {
	return &UserProfile{Account: accname, Game: game, TaskConfigs: taskcfgpath}
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

type OcrConfig struct {
	Split      []string `yaml:"split"`
	Imagick    []string `yaml:"imagick"`
	Tesseract  []string `yaml:"tesseract"`
	Exceptions []string `yaml:"dict_shrt_except"`
}

type Location struct {
	Key       string   `yaml:"name"`
	Threshold int      `yaml:"hits"`
	Keywords  []string `yaml:"keywords"`
}

func (l *Location) String() string {
	return fmt.Sprintf("Location key: %v", l.Key)
}

var defaultAppConfig = &AppConfig{
	DeviceSerial: "",
	UserProfile: &UserProfile{
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
	Logfile:      "app.conf",
	Loglevel:     "FATAL",
	DrawStep:     false,
	Dirs: struct {
		Root     string `yaml:"root"`
		TempImg  string `yaml:"tempImg"`
		SqDB     string `yaml:"sqDB"`
		GameConf string `yaml:"gameConf"`
		TestData string `yaml:"testData"`
	}{
		Root:     ".afk_data",
		TempImg:  "work_images",
		SqDB:     "db",
		GameConf: "cfg",
		TestData: "_test",
	},
	RequiredInstalledSoftware: []string{magic, adbp, tesseract, bluestacks},
}
