package cfg

// Required 3rd-party software
const (
	AdbExe   Executable = "adb"
	MagicExe            = "magick"
	TessExe             = "tesseract"
	// dup in emulator
	BluestacksExe = "HD-Player"
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
	ClimbKings
	ClimbWild
	ClimbGrave
	ClimbInferno
	ClimbMaul
	ClimbLight
	ClimbCelestial
	GuildBosses
)

var (
	userTemplate = &Profile{
		DeviceSerial: "",
		User: &User{
			Account:     "",
			Game:        "",
			TaskConfigs: []string{"cfg/reactions.yaml", "cfg/daily.yaml"},
		},
		Imagick: imgksArggs(),
		AltImagick: []string{
			"-colorspace", "Gray",
			"-alpha", "off",
			"-threshold", "75%",
			"-edge", "2",
			"-negate",
			"-black-threshold", "90%",
		},
		Tesseract:    tssA(),
		AltTesseract: []string{"--psm", "3", "hoot", "quiet"},
		Bluestacks:   &Bluestacks{Instance: "Rvc64", Package: "com.lilithgames.hgame.gp.id"},
		Exceptions:   []string{"Go", "Up ", "In", "Tap"},
		Loglevel:     "FATAL",
		DrawStep:     false,
	}
	imgksArggs = func() map[int]CmdArgs {
		r := make(map[int]CmdArgs, 0)
		r[1] = CmdArgs{Key: "-colorspace", Val: "Gray"}
		r[2] = CmdArgs{Key: "-alpha", Val: "off"}
		r[3] = CmdArgs{Key: "-threshold", Val: "75%"}

		return r
	}
	tssA = func() map[int]CmdArgs {
		r := make(map[int]CmdArgs, 0)
		r[1] = CmdArgs{Key: "--psm", Val: "6"}
		r[2] = CmdArgs{Key: "-c", Val: "tessedit_create_alto=1"}
		r[3] = CmdArgs{Key: "hoot", Val: "quiet"}

		return r
	}
)
