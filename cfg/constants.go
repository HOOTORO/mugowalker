package cfg

// Required 3rd-party software
const (
	adbp       = "adb"
	magic      = "magick"
	tesseract  = "tesseract"
	bluestacks = "HD-Player"
)

func thirdparty() []string {
	return []string{adbp, magic, tesseract, bluestacks}
}

// required sys events
const (
	appdataEnv  = "APPDATA"
	programData = "ProgramData"
	temp        = "TEMP"
	userhome    = "USERPROFILE"
)

const (
	programRootDir = ".afkworker"
	logfile        = "app.log"
	// tempImg = "work_images"
	// sqDB =  "db"
	// gameConf =  "cfg"
	// testData: _test

)
