package cfg

import "fmt"

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
		Root     string `yaml:"root"`
		TempImg  string `yaml:"tempImg"`
		SqDB     string `yaml:"sqDB"`
//		User     string `yaml:"user"`
		GameConf string `yaml:"gameConf"`
		TestData string `yaml:"testData"`
	}
	RequiredInstalledSoftware []string `yaml:"required_installed_software"`
}

func (ac *AppConfig) String() string {
    return fmt.Sprintf("DeviceId: %v\nConfolder: %v", ac.DeviceSerial, ac.Dirs.GameConf )

}
type UserProfile struct {
	Account     string
	Game        string
	TaskConfigs []string
}

func (up *UserProfile) String() string {
	return fmt.Sprintf("\n --> Account: %v\n" +
        "     Game: %v\n", up.Account, up.Game)
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
	Threshold int      `yaml:"hits"`
	Keywords  []string `yaml:"keywords"`
}

type emuConf []struct {
	Cmd  string   `yaml:"cmd"`
	Args []string `yaml:"args"`
}
