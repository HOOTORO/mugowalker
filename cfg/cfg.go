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
	defaultcfg = "assets/default.yaml"
	game       = "AFK Arena"
)

var (
	ErrWorkDirFail        = errors.New("working dirictories wasn't created. Exit")
	ErrRequiredProgram404 = errors.New("missing some of required soft")
)

var (
	log        *logrus.Logger
	activeUser *Profile
	sysvars    *SystemVars
	red, green     func(...interface{}) string
)

func init() {
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	log = Logger()
	sysvars = loadSysconf()

	f, e := os.OpenFile(sysvars.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if e == nil {
		sysvars.Logfile = f.Name()
		log.SetOutput(f)
	}

	activeUser = loadUser()

	loglvl, e := logrus.ParseLevel(activeUser.Loglevel)
	if e != nil {
		log.Errorf("logrus err: %v", e)
	}
	log.SetLevel(loglvl)


	if e != nil {
		panic(e)
	}

	repository.DbInit(func(x string) string {
		return filepath.Join(sysvars.Db, x)
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

func ActiveUser() *Profile{
	if activeUser != nil {
		return activeUser
	}
	return defUser
}
func RunBlue() error {
	args := activeUser.Bluestacks
	cmd := exec.Command(bluestacks, args...)
	log.Tracef("cmd bs : %v\n", cmd.String())
	return cmd.Start()
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

//func Load(a *Profile) *adb.Device {
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

func LoadTask(up *User) (r []ReactiveTask) {
	for _, t := range up.TaskConfigs {
		reactiveTasks := make([]ReactiveTask, 0)
		Parse(t, &reactiveTasks)
		r = append(r, reactiveTasks...)
	}
	return
}

func GetImages() []string {
	d, e := os.ReadDir(sysvars.Temp)
	if e != nil {
		panic("get imgs")
	}
	var res []string
	for _, f := range d {
		res = append(res, TempFile(f.Name()))
	}
	return res
}

func TempFile(f string) string {
	return absJoin(sysvars.Temp, f)
}

func UserFile(f string) string {
	return absJoin(sysvars.Userhome, f)
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
func ToInt(s string) int {
	num, e := strconv.Atoi(s)
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v", e, "intconv")
	}
	return num
}

func loadSysconf() *SystemVars {
	sys := &SystemVars{}
	e :=	loadParties(sys)
	if e != nil {
		log.Errorf("Load parties mailfunc: %v", e)
	}
	sys.Db, sys.Userhome, sys.App, sys.Temp, e = createDirStructure()
	if e != nil {
		log.Errorf("Create app folders mailfunc: %v", e)
	}
	sys.Logfile = logfile
	return sys
}
func loadUser() *Profile {
	conf := &Profile{}
	lastcfg := lookupLastConfig(sysvars.Userhome, sysvars.Db)

	e := Parse(lastcfg, conf)
	if e != nil {
		conf = defaultUser()
	} else {
		sysvars.UserConfPath = UserFile(lastcfg)
	}
	return conf
}

func safeEnv(n string) string {
	str, ok := os.LookupEnv(n)
	if ok {
		return str
	}
	log.Warnf("$Env:%v NOT FOUND, BE AWARE", n)
	return ""
}

func defaultUser() *Profile {
	settings := defUser
	cfgpath := UserFile(defaultcfg)
	sysvars.UserConfPath = cfgpath
	Save(defaultcfg, settings)

	return settings
}
func loadParties(sys *SystemVars) error {
	for _, s := range thirdparty() {
		if pt := LookupPath(s); pt != "" {
			sys.parties = append(sys.parties, &RunableExe{name: s, path: pt})
		} else {
			return ErrRequiredProgram404
		}
	}
	return nil
}

func lookupLastConfig(dirs ...string) string {
	last := time.Time{}
	res := ""
	for _, d := range dirs {
		dir, e := os.ReadDir(d)
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

func cutgrid(str string) (p image.Point, off int) {
	off = 1 // default
	ords := strings.Split(str, ":")
	p = image.Point{
		X: ToInt(ords[0]),
		Y: ToInt(ords[1]),
	}
	if len(ords) > 2 {
		off = ToInt(ords[2])
	}
	return
}

func createDirStructure() (dbfolder, userfolder, appdata, tempfolder string, e error) {
	dbfolder = makeEnvDir(appdataEnv, programRootDir)
	userfolder = makeEnvDir(userhome, programRootDir)
	tempfolder = makeEnvDir(temp, programRootDir)
	appdata = makeEnvDir(programData, programRootDir)

	// Saturday cleaning
	if time.Now().Weekday().String() == "Saturday" {
		truncateDir(tempfolder)
	}

	if dbfolder == "" || userfolder == "" || tempfolder == "" || appdata == "" {
		e = ErrWorkDirFail
	}

	log.Infof("\ninit: success; dirs created: \n%v\n%v\n%v\n%v", dbfolder, userfolder, tempfolder, appdata)
	return
}

func makeEnvDir(env, dir string) string {
	envpath := safeEnv(env)
	patyh := filepath.Join(envpath, dir)
	e := os.MkdirAll(patyh, os.ModeDir)
	if e != nil {
		log.Errorf("make dir mailfunc: %v", e)
	}
	return patyh
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


//app := &cfg.Profile{
//    DeviceSerial: "192.168.1.7:5555",
//    User: &cfg.User{
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
