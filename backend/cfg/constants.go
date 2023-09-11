package cfg

// required win sys envs
const (
	appdataEnv  = "APPDATA"
	programData = "ProgramData"
	temp        = "TEMP"
	userhome    = "USERPROFILE"
)

const (
	macEnv = "HOME"
)

const (
	programRootDir = ".afkworker"
	dbfolder       = "db"
	tempfolder     = "temp"
	logfile        = "app.log"
)

// ///////////////////////////
// Global afk activities ///
// /////////////////////////
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
	Daily
)
const (
	afkapp     = "com.lilithgames.hgame.gp"
	afktestapp = "com.lilithgames.hgame.gp.id"
)

var (
	userTemplate = &Profile{
		DeviceSerial: "",
		GameAccount:  "",
		// TaskConfigs: []string{"cfg/reactions.yaml", "cfg/daily.yaml"},

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
		Bluestacks:   &Bluestacks{Instance: "Rvc64"},
		Exceptions:   []string{"Go", "Up ", "In", "Tap"},
		Loglvl:       "TRACE",
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
