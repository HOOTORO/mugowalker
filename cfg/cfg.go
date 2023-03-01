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
)

const (
	defaultcfg = "assets/default.yaml"
	game       = "AFK Arena"
)

type Runnable interface {
	Path() string
	Args() []string
}

var (
	ErrWorkDirFail        = errors.New("working dirictories wasn't created. Exit")
	ErrRequiredProgram404 = errors.New("missing some of required soft")
	ErrLoadInitConf       = errors.New("load sysvars")
)
var (
	// F... format alias func(...interface{}) string
	F = fmt.Sprintf
	// Red coloring Sprint
	Red = color.New(color.FgHiRed).SprintFunc()
	// Green coloring Sprint
	Green = color.New(color.FgHiGreen).SprintFunc()
	// Cyan coloring Sprint
	Cyan = color.New(color.FgHiCyan).SprintFunc()
	// Blue coloring Sprint
	Blue   = color.New(color.FgHiBlue).SprintFunc()
	Ylw    = color.New(color.FgHiYellow).SprintFunc()
	Mgt    = color.New(color.FgHiMagenta).SprintFunc()
	TTrack = color.New(color.BgHiBlue, color.FgCyan, color.Underline, color.Bold).SprintfFunc()
	RFW    = color.New(color.FgHiRed, color.BgWhite).SprintFunc()
	MgCy   = color.New(color.FgHiMagenta, color.BgCyan).SprintFunc()
)

var (
	log        *logrus.Logger
	activeUser *Profile
	sysvars    *SystemVars
)

func init() {
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

// Deserialize bits to string values
func Deserialize[T emum](raw T) []string {
	var result []string
	for i := 0; i < len(raw.Values()); i++ {
		if d := T(1 << i); raw&(1<<uint(i)) != 0 {
			result = append(result, d.String())
		}
	}
	return result
}

// Logger for app to use
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

// ActiveUser or template wwithout name, connect, gameID
func ActiveUser() *Profile {
	if activeUser != nil {
		return activeUser
	}
	return userTemplate
}

// LastLoaded <userconf>.yaml
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

// UpdateUserInfo saves to yaml into Userhome dir
func UpdateUserInfo(au AppUser) {
	activeUser.DeviceSerial = au.DevicePath()
	activeUser.Loglevel = au.Loglevel()
	activeUser.User.Account = au.Acccount()
	activeUser.User.Game = au.Game()
	Save(UserFile(au.Acccount()+".yaml"), activeUser)

}

// GetImages from temp/<appfolder>
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

func RunCmd(r Runnable) error {
	pt := LookupPath(r.Path())

	log.Trace(Blue("  ↓   RunCMD   ↓ \n", Mgt(pt), "\n", Ylw(r.Args())))
	cmd := exec.Command(pt, r.Args()...)

	return cmd.Run()
}

// TempFile in <temp>/<appfolder>/*
func TempFile(f string) string {
	return absJoin(sysvars.Temp, f)
}

// UserFile from $env:USERPROFILE
func UserFile(f string) string {
	return absJoin(sysvars.Userhome, f)
}

// LookupPath for exe s
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
		log.Warnf("Calledc.F():%v\nError:%v", "cfg.ToInt", e)
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

func createDirStructure() (dbfolder, userfolder, appdata, tempfolder string, e error) {
	dbfolder = makeEnvDir(appdataEnv, programRootDir)
	userfolder = makeEnvDir(userhome, programRootDir)
	appdata = makeEnvDir(programData, programRootDir)
	tempfolder = makeEnvDir(temp, programRootDir)

	// Saturday cleaning

	if dbfolder == "" || userfolder == "" || tempfolder == "" || appdata == "" {
		e = ErrWorkDirFail
	}

	log.Infof("\ninit: success; dirs created: \n%v\n%v\n%v\n%v", dbfolder, userfolder, tempfolder, appdata)
	return
}

func makeEnvDir(env, dir string) string {
	envpath := safeEnv(env)
	patyh := filepath.Join(envpath, dir)
	if env == temp {
		truncateDir(patyh)
	}
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

func Shorterer(str string, n int) string {
	if len(str) > n+3 {
		return str[:n] + "..."
	}
	return str
}
