package cfg

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"worker/adb"
	"worker/afk/repository"

	"github.com/fatih/color"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	tmpdir      = "imag"
	dbdir       = "db"
	cfgdir      = "cfg"
	stepsdir    = "steps"
	localdir    = ".afk_data"
	adbdir      = ".adb"
	ocrsettings = "cfg/ocr.yaml"
	emulator    = "cfg/emu.yaml"
    appsettings = "cfg/appcfg.yaml"
	usrfolder   = "usrdata"
)

var (
	ErrWorkDirFail  error = errors.New("working dirictories wasn't created. Exit")
	ErrStepNotFound error = errors.New("Config error. Have conditional step, but no actions for it")
)
var (
	log *logrus.Logger
    instanceConfig *AppConfig 
)


var (
    AppConf *AppConfig
	OcrConf *OcrConfig
	EmulatorConf *emuConf
)

func init() {
	log = Logger()
	e := createDirStructure()
	f, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE, 0o644)
	log.SetLevel(logrus.TraceLevel)
	log.SetOutput(f)
	if e != nil {
		panic(e)
	}
    loadConf()
	repository.DbInit(func(x string) string {
		return filepath.Join(dbdir, x)
	})
}

func Logger() *logrus.Logger {
	if log != nil {
		return log
	}
	return &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			PadLevelText:              true,
			TimestampFormat:           time.Stamp,
		},
		Level: logrus.InfoLevel,
	}
}

func (l *Location) String() string {
	return fmt.Sprintf("Location key: %v", l.Key)
}

func (rt ReactiveTask) React(trigger string) *adb.Point {
	for _, v := range rt.Reactions {
		if trigger == v.If {
			return cutgrid(v.Do)
		}
	}
	return cutgrid("1:18")
}

func (rt ReactiveTask) Before(trigger string) string {
	for _, v := range rt.Reactions {
		if trigger == v.If && v.Before != "" {
			return v.Before
		}
	}
	return ""
}

func (rt ReactiveTask) After(trigger string) string {
	for _, v := range rt.Reactions {
		if trigger == v.If && v.After != "" {
			return v.After
		}
	}
	return ""
}

func Parse(s string, out interface{}) {
	f, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, out)
	if err != nil {
		log.Fatalf("UNMARSHAL WASTED: %v", err)
	}
	log.Tracef("UNMARSHALLED: %v\n\n", out)
}

func Save(name string, in interface{}) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
    b, err := yaml.Marshal(in)
	if err != nil {
		log.Fatalf("MARSHAL WASTED: %v", err)
	}
    f.Write(b)
	log.Tracef("MARSHALLED: %v\n\n", f)
}



func UserInput(desc, def string) string {
	reader := bufio.NewReader(os.Stdin)
	bufio.NewReader(os.Stdin)
	//	text := "0"
	color.HiCyan(desc)
	color.HiRed("---------------------")
	fmt.Printf("[default:%v]: ", color.HiGreenString(def))
	text, _ := reader.ReadString('\n')
	//		text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	if len(text) == 0 {
		text = def
	}
	return strings.Trim(text, "\r")
}

func toInt(s string) int {
	num, e := strconv.Atoi(s)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v", e, "intconv")
	}
	return num
}

func cutgrid(str string) (p *adb.Point) {
	ords := strings.Split(str, ":")
	p = &adb.Point{
		Point: image.Point{
			X: toInt(ords[0]),
			Y: toInt(ords[1]),
		},
		Offset: 1,
	}
	if len(ords) > 2 {
		p.Offset = toInt(ords[2])
	}
	return
}

func StrToGrid(str string) (p *adb.Point) {
	ords := strings.Split(str, ":")
	p = &adb.Point{
		Point: image.Point{
			X: toInt(ords[0]),
			Y: toInt(ords[1]),
		},
		Offset: 1,
	}
	if len(ords) > 2 {
		p.Offset = toInt(ords[2])
	}
	return
}

func loadConf() {
    AppConf = &AppConfig{}
    OcrConf = &OcrConfig{}
    EmulatorConf = &emuConf{}
    Parse(appsettings, AppConf)
	Parse(ocrsettings, OcrConf)
	Parse(emulator, EmulatorConf)

}

func truncateDir(d string) {
	a, _ := filepath.Abs(d)
	_ = os.RemoveAll(a)
}

func GetImages() []string {
	d, e := os.ReadDir(tmpdir)
	if e != nil {
		panic("crop err")
	}
	var res []string
	for _, f := range d {
		res = append(res, ImageDir(f.Name()))
	}
	return res
}

func ImageDir(f string) string {
	a, e := filepath.Abs(tmpdir)
	if e != nil {
		panic("no tmpdir")
	}
	return filepath.Join(a, filepath.Base(f))
}

func UsrDir(f string) string {
	a, e := filepath.Abs(usrfolder)
	if e != nil {
		panic("no tmpdir")
	}
	return filepath.Join(a, filepath.Base(f))
}

func createDirStructure() error {
	//	usr := SafeEnv("USERPROFILE")
	//	wd := filepath.Join(usr, adbdir)
	wd, e := os.Getwd()
	rootd := filepath.Join(wd, localdir)
	usr := filepath.Join(wd, localdir, usrfolder)
	dbpath := filepath.Join(wd, localdir, dbdir)
	imgpath := filepath.Join(wd, localdir, tmpdir)
	cfgpath := filepath.Join(wd, localdir, cfgdir)
	stepspath := filepath.Join(wd, localdir, stepsdir)

	truncateDir(imgpath)

	e = os.MkdirAll(cfgpath, os.ModeDir)
	e = os.MkdirAll(dbpath, os.ModeDir)
	e = os.MkdirAll(imgpath, os.ModeDir)
	e = os.MkdirAll(stepspath, os.ModeDir)
	e = os.MkdirAll(usr, os.ModeDir)

	if e != nil {
		log.Errorf("%v", e)
		return ErrWorkDirFail
	}
	e = os.Chdir(rootd)

	if e == nil {
		pwd, _ := os.Getwd()
		log.Infof("\ninit: success; pwd: %v\n\n", pwd)
	}
	return e
}

func LookupPath(name string) (path string) {
	p, err := exec.LookPath(name)
	if err == nil {
		if p, err = filepath.Abs(p); err == nil {
			return p
		}
	}
	panic(fmt.Sprintf("Required programm: %v not found in path\n error: %v", name, err))
}

func Load(a *AppConfig) *adb.Device {
	devs, e := adb.Devices()
	if e != nil || len(devs) == 0 {
		d, e := adb.Connect(a.DeviceSerial)
		if e != nil {
			panic("dev err")
		}
		devs = append(devs, d)

	}
	num := 0
	if len(devs) > 1 {
		var desc string = "Choose, which one will be used by bot\n"
		for i, dev := range devs {
			desc += fmt.Sprintf("%v: Serial-> %v,   id-> %v,    resolution-> %v\n", i, dev.Serial, dev.TransportId, dev.Resolution)
		}
		num, _ = strconv.Atoi(UserInput(desc, "0"))
	}
	return devs[num]
}

func LoadTask(up *UserProfile) (r []ReactiveTask) {
	for _, t := range up.TaskConfigs {
		reactiveTasks := make([]ReactiveTask, 0)
		Parse(t, &reactiveTasks)
		r = append(r, reactiveTasks...)
	}
	return
}
