package cfg

// Required 3rd-party software
const (
	AdbExe        Executable = "adb"
	MagicExe                 = "magick"
	TessExe                  = "tesseract"
	BluestacksExe            = "HD-Player"
)

type Executable string

func (e Executable) String() string {
	return string(e)
}

// required sys envs
const (
	appdataEnv  = "APPDATA"
	programData = "ProgramData"
	temp        = "TEMP"
	userhome    = "USERPROFILE"
)

const (
	programRootDir = ".afkworker"
	logfile        = "app.log"
)

// Global afk activities

type Mission int

const (
	PushCampain Mission = iota + 1
	ClimbTower
	GuildBosses
)
