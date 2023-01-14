package cfg

import (
	"errors"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"worker/afk/repository"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	appdataEnv  = "APPDATA"
	profileEnv  = "USERPROFILE"
	programData = "ProgramData"
	temp        = "TEMP"
	defaultcfg  = "assets/default.yaml"
)

const (
	game = "AFK Arena"
)

var roamdata, userfolder, appdata, tempdir string

var (
	ErrWorkDirFail        = errors.New("working dirictories wasn't created. Exit")
	ErrRequiredProgram404 = errors.New("missing some of required soft")
)

var (
	log        *logrus.Logger
	Env        *AppConfig
	red, green func(...interface{}) string
)

func init() {
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()

	log = Logger()
	if Env == nil {
		Env = defaultAppConfig
	}

	e := createDirStructure()

	Env = loadConf()

	e = Env.validateDependencies()

	f, e := os.OpenFile(Env.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	// defer f.Close()
	if e == nil {
		Env.Logfile = f.Name()
	}

	loglvl, e := logrus.ParseLevel(Env.Loglevel)
	if e != nil {
		log.Errorf("logrus err: %v", e)
	}

	log.SetLevel(loglvl)
	log.SetOutput(f)

	if e != nil {
		panic(e)
	}
	a, _ := f.Stat()
	log.Warnf("somepepeedoor%+v", a)
	repository.DbInit(func(x string) string {
		return filepath.Join(roamdata, x)
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
		Level: logrus.FatalLevel,
	}
}

func (rt ReactiveTask) React(trigger string) (image.Point, int) {
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

func Parse(s string, out interface{}) error {
	f, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, out)
	if err != nil {
		log.Fatalf("UNMARSHAL WASTED: %v", err)
	}
	log.Tracef("UNMARSHALLED: %v\n\n", out)
	return err
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
	_, err = f.Write(b)
	if err != nil {
		log.Errorf("write yaml (e): %v", err)
	}
	log.Tracef("MARSHALLED: %v\n\n", f)
}

//func Load(a *AppConfig) *adb.Device {
//	devs, e := adb.Devices()
//	if e != nil || len(devs) == 0 {
//		d, e := adb.Connect(a.DeviceSerial)
//		if e != nil {
//			panic("dev err")
//		}
//		devs = append(devs, d)
//
//	}
//	num := 0
//	if len(devs) > 1 {
//		var desc string = "Choose, which one will be used by bot\n"
//		for i, dev := range devs {
//			desc += fmt.Sprintf("%v: Serial-> %v,   id-> %v,    resolution-> %v\n", i, dev.Serial, dev.TransportId, dev.Resolution)
//		}
//		num, _ = strconv.Atoi(ui.UserInput(desc, "0"))
//	}
//	return devs[num]
//}

func LoadTask(up *UserProfile) (r []ReactiveTask) {
	for _, t := range up.TaskConfigs {
		reactiveTasks := make([]ReactiveTask, 0)
		Parse(t, &reactiveTasks)
		r = append(r, reactiveTasks...)
	}
	return
}

func GetImages() []string {
	d, e := os.ReadDir(tempdir)
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
	return absJoin(tempdir, f)
}

func UsrDir(f string) string {
	return absJoin(userfolder, f)
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

/*
	Helper func
*/

func safeEnv(n string) string {
	str, ok := os.LookupEnv(n)
	if ok {
		return str
	}
	log.Warnf("$Env:%v NOT FOUND, BE AWARE", n)
	return ""
}

func loadConf() *AppConfig {
	conf := &AppConfig{}
	lastcfg := lookupLastConfig()

	e := Parse(lastcfg, conf)
	if e != nil {
		conf = inputminsettings()
	} else {
		conf.Thiscfg = UsrDir(lastcfg)
	}
	return conf
}

func inputminsettings() *AppConfig {
	settings := defaultAppConfig
	cfgpath := UsrDir(defaultcfg)
	settings.Thiscfg = cfgpath
	Save(defaultcfg, settings)

	return settings
}

func lookupLastConfig() string {
	lookoutdirs := []string{userfolder, appdata}
	last := time.Time{}
	res := ""
	for _, d := range lookoutdirs {
		dir, e := os.ReadDir(userfolder)
		if e != nil {
			fmt.Printf("\nerr:%v\nduring run:%v", e, "lookout")
		}
		for _, entry := range dir {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
				i, _ := entry.Info()
				if i.ModTime().After(last) {
					last = i.ModTime()
					res = filepath.Join(d, i.Name())
				}
			}
		}

	}
	return res
}

func toInt(s string) int {
	num, e := strconv.Atoi(s)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v", e, "intconv")
	}
	return num
}

func cutgrid(str string) (p image.Point, off int) {
	off = 1 // default
	ords := strings.Split(str, ":")
	p = image.Point{
		X: toInt(ords[0]),
		Y: toInt(ords[1]),
	}
	if len(ords) > 2 {
		off = toInt(ords[2])
	}
	return
}

func createDirStructure() error {
	roamdata = makeEnvDir(appdataEnv, Env.Dirs.SqDB)
	userfolder = makeEnvDir(profileEnv, Env.Dirs.GameConf)
	tempdir = makeEnvDir(temp, Env.Dirs.TempImg)
	appdata = makeEnvDir(programData, Env.Dirs.TestData)

	// Saturday cleaning
	if time.Now().Weekday().String() == "Saturday" {
		truncateDir(tempdir)
	}

	if roamdata == safeEnv(appdataEnv) || userfolder == safeEnv(profileEnv) || tempdir == safeEnv(temp) || appdata == safeEnv(programData) {
		return ErrWorkDirFail
	}

	log.Infof("\ninit: success; dirs created: \n%v\n%v\n%v\n%v", roamdata, userfolder, tempdir, appdata)
	return nil
}

func makeEnvDir(env, dir string) string {
	envpath := safeEnv(env)
	patyh := filepath.Join(envpath, Env.Dirs.Root, dir)
	e := os.MkdirAll(patyh, os.ModeDir)
	if e != nil {
		log.Errorf("make dir mailfunc: %v", e)
	}
	return patyh
}

func (ac *AppConfig) validateDependencies() error {
	for i, s := range ac.RequiredInstalledSoftware {
		if pt := LookupPath(s); pt != "" {
			ac.RequiredInstalledSoftware[i] = pt
		} else {
			return ErrRequiredProgram404
		}
	}
	return nil
}

func truncateDir(d string) {
	a, _ := filepath.Abs(d)
	//    _ = os.RemoveAll(a)
	fmt.Printf("DELETED %v\n", a)
}

func absJoin(d, f string) string {
	fb := filepath.Base(f)
	if filepath.IsAbs(d) {
		return filepath.Join(d, fb)
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, fb)
}

func RunBlue() error {
	args := Env.Bluestacks
	cmd := exec.Command(bluestacks, args...)
	log.Tracef("cmd bs : %v\n", cmd.String())
	// uncomment for ocr log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Start()
}

//app := &cfg.AppConfig{
//    DeviceSerial: "192.168.1.7:5555",
//    UserProfile: &cfg.UserProfile{
//        Account:     "E6osh!ro",
//        Game:        "AFK Arena",
//        TaskConfigs: []string{"cfg/reactions.yaml", "cfg/daily.yaml"},
//        },
//        Imagick: cfg.OcrConf.Imagick,
//        AltImagick: []string{"-colorspace", "Gray", "-alpha", "off", "-threshold", "75%", "-edge", "2", "-negate", "-black-threshold",
//            //			"-white-threshold",
//            //			"60%",
//            "90%",
//            },
//            Tesseract:    cfg.OcrConf.Tesseract,
//            AltTesseract: []string{"--psm", "3", "hoot", "quiet"},
//            Bluestacks:   []string{"--instance", "Rvc64_16", "--cmd", "launchApp", "--package", "com.lilithgames.hgame.gp.id"},
//            Exceptions:   cfg.OcrConf.Exceptions,
//            Loglevel:     "INFO",
//            DrawStep:     false,
//            Dirs: struct {
//        Logfile     string `yaml:"logfile"`
//        Root     string `yaml:"rootDir"`
//        TempImg  string `yaml:"tempImgDir"`
//        SqDB    string `yaml:"sqDBDir"`
//        User    string `yaml:"userDir"`
//        GameConf string `yaml:"gameConfDir"`
//        TestData string `yaml:"testDataDir"`
//            }{
//        Logfile:     "app.log",
//        Root:     ".afk_data",
//        TempImg:  "work_images",
//        SqDB:     "db",
//        User:     "usrdata",
//        GameConf: "cfg",
//        TestData: "_test",
//        },
//        }
//
//
//        cfg.Save("runset.yaml", app)
