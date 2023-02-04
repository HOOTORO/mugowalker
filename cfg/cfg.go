package cfg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	ErrLoadInitConf       = errors.New("load sysvars")
)

var (
	log              *logrus.Logger
	activeUser       *Profile
	sysvars          *SystemVars
	red, green, cyan func(...interface{}) string
)
var f = fmt.Sprintf

func init() {
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
	log = Logger()
	sysvars, _ = loadSysconf()

	f, e := os.OpenFile(sysvars.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if e == nil {
		sysvars.Logfile = f.Name()
		log.SetOutput(f)
	}

	activeUser = LastLoaded()

	ll, e := logrus.ParseLevel(activeUser.Loglevel)
	// log.SetReportCaller(true)

	if e != nil {
		panic(e)
	}

	log.SetLevel(ll)

	repository.DbInit(func(x string) string {
		return filepath.Join(sysvars.Db, x)
	})
}

type emum interface {
	~uint
	String() string
	Values() []string
}

func Deserialize[T emum](raw T) []string {
	var result []string
	for i := 0; i < len(raw.Values()); i++ {
		if d := T(1 << i); raw&(1<<uint(i)) != 0 {
			result = append(result, d.String())
		}
	}
	return result
}

// func Deserialize[T int, R string](raw T, opts []R) []T {
// 	var result []T
// 	for i := 0; i < len(opts); i++ {
// 		if d := T(1 << i); raw&(1<<uint(i)) != 0 {
// 			result = append(result, d)
// 		}
// 	}
// 	return result
// }

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

func ActiveUser() *Profile {
	if activeUser != nil {
		return activeUser
	}
	return userTemplate
}
func LastLoaded() *Profile {
	conf := &Profile{}
	lastcfg := mostRecentModifiedYAML(sysvars.Userhome, sysvars.Db)
	if lastcfg != "" {
		e := Parse(lastcfg, conf)
		if e != nil {
			log.Errorf("Err: %v", e)
		}
	}
	return conf
}

func UpdateUserInfo(au AppUser) {
	activeUser.DeviceSerial = au.DevicePath()
	activeUser.Loglevel = au.Loglevel()
	activeUser.User.Account = au.Acccount()
	activeUser.User.Game = au.Game()
	Save(UserFile(au.Acccount()+".yaml"), activeUser)

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
	log.Tracef("\n\n\n\n\nUNMARSHALLED: %v\n\n", out)
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
		log.Warnf("Called f():%v\nError:%v", "cfg.ToInt", e)
	}
	return num
}

func loadSysconf() (sys *SystemVars, e error) {
	sys = &SystemVars{}

	sys.Db, sys.Userhome, sys.App, sys.Temp, e = createDirStructure()
	if e != nil {
		log.Errorf("Create app folders mailfunc: %v", e)
	}
	sys.Logfile = logfile
	return
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
	settings := userTemplate
	cfgpath := UserFile(defaultcfg)
	sysvars.UserConfPath = cfgpath
	Save(defaultcfg, settings)

	return settings
}

func mostRecentModifiedYAML(dirs ...string) string {
	last := time.Time{}
	res := ""
	for _, d := range dirs {
		dir, e := os.ReadDir(d)
		if e != nil {
			log.Errorf("\nerr:%v\nduring run:%v", e, "lookout")
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
	_ = os.RemoveAll(a)
	log.Warnf("DELETED %v\n", a)
}

func absJoin(d, f string) string {
	fb := filepath.Base(f)
	if filepath.IsAbs(d) {
		return filepath.Join(d, fb)
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, fb)
}
