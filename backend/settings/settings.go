package settings

import (
	"mugowalker/backend/cfg"
	// "os"
	// "time"
	// "github.com/sirupsen/logrus"
)

var defaultcfg = "backend/assets/cfg.yml"

// var log *logrus.Logger

const (
	TRACE = "Trace"
	INFO  = "Info"
	WARN  = "Warn"
	ERR   = "Error"
	FATAL = "Fatal"
)

type Settings struct {
	DrawStep bool   `yaml:"draw_step"`
	Logfile  string `yaml:"logfile"`
	Loglevel string `yaml:"loglevel"`
	Imagick  struct {
		Colorspace     string `yaml:"colorspace"`
		Alpha          string `yaml:"alpha"`
		Threshold      string `yaml:"threshold"`
		Edge           string `yaml:"edge"`
		Negate         bool   `yaml:"negate"`
		BlackThreshold string `yaml:"black-threshold"`
	} `yaml:"imagick"`
	Tesseract struct {
		Psm  int      `yaml:"psm"`
		Args []string `yaml:"args"`
	} `yaml:"tesseract"`
	Bluestacks struct {
		Instance string `yaml:"instance"`
		Cmd      string `yaml:"cmd"`
		Package  string `yaml:"package"`
	} `yaml:"bluestacks"`
	IgnoredWords []string `yaml:"ignored_words"`
}

func Default() *Settings {
	def := &Settings{}
	cfg.Parse(defaultcfg, def)

	// f, e := os.OpenFile(def.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	// if e != nil {
	// 	panic(e)
	// }
	// ll, e := logrus.ParseLevel(def.Loglevel)
	// if e != nil {
	// 	panic(e)
	// }
	// log = &logrus.Logger{
	// 	Out: f,
	// 	Formatter: &logrus.TextFormatter{
	// 		TimestampFormat: time.Stamp,
	// 		ForceQuote:      true,
	// 		FieldMap:        logrus.FieldMap{},
	// 	},
	// 	Level: ll,
	// }
	return def
}

// func (s *Settings) Log(lvl, msg string) {
// 	l, _ := logrus.ParseLevel(lvl)
// 	log.Logf(l, "%s: %s", lvl, msg)
// }
