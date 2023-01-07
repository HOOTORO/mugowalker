package cfg

import "fmt"

type AppConfig struct {
	DeviceSerial string `yaml:"connection_str"`

	UserProfile *UserProfile
	//  Recognition settings (cmd args for 'Imagick' and 'Tesseract')
	//  Split     []string `yaml:"split"`
	Imagick      []string `yaml:"imagick"`
	AltImagick   []string `yaml:"alt_imagick"`
	Tesseract    []string `yaml:"tesseract"`
	AltTesseract []string `yaml:"alt_tesseract"`
	Bluestacks   []string `yaml:"bluestacks"`

	// Dict short word exceptions (>= 3)
	Exceptions []string `yaml:"dict_shrt_except"`

	Loglevel string `yaml:"loglevel"`
	DrawStep bool   `yaml:"draw_step"`

	Folders  struct {
        Logfile string `yaml:"logfile"`
		RootDir     string `yaml:"rootDir"`
		TempImgDir  string `yaml:"tempImgDir"`
		SqDBDir     string `yaml:"sqDBDir"`
		UserDir     string `yaml:"userDir"`
		GameConfDir string `yaml:"gameConfDir"`
        TestDataDir string `yaml:"testDataDir"`
	}
}
type UserProfile struct {
	Account     string
	Game        string
	TaskConfigs []string
}

func (up *UserProfile) String() string {
	return fmt.Sprintf("\n--> Game: %v\n--> Acc: %v\n", up.Game, up.Account)
}
func User(accname, game string, taskcfgpath []string) *UserProfile {
	return &UserProfile{Account: accname, Game: game, TaskConfigs: taskcfgpath}
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

type OcrConfig struct {
	Split      []string `yaml:"split"`
	Imagick    []string `yaml:"imagick"`
	Tesseract  []string `yaml:"tesseract"`
	Exceptions []string `yaml:"dict_shrt_except"`
}

type Location struct {
	Key       string   `yaml:"name"`
	Grid      string   `yaml:"grid"`
	Threshold int      `yaml:"hits,omitempty"`
	Keywords  []string `yaml:"keywords"`
	Wait      bool     `yaml:"wait"`
	// Actions   []*Point `yaml:"actions"`
}

type emuConf []struct {
	Cmd  string   `yaml:"cmd"`
	Args []string `yaml:"args"`
}
