package cfg

import (
	"bufio"
	"errors"
	"fmt"
    "golang.org/x/exp/slices"
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
	localdir    = ".afk_data"
	adbdir      = ".adb"
	ocrsettings = "cfg/ocr.yaml"
)

var (
    ErrWorkDirFail error = errors.New("working dirictories wasn't created. Exit")
    ErrStepNotFound error = errors.New("Config error. Have conditional step, but no actions for it")
)
var log *logrus.Logger
var OcrConf *ocrConfig

type ocrConfig struct {
	Split     []string `yaml:"split"`
	Imagick   []string `yaml:"imagick"`
	Tesseract []string `yaml:"tesseract"`
}
type Location struct {
	Key       string   `yaml:"name"`
	Grid      string   `yaml:"grid"`
	Threshold int      `yaml:"hits,omitempty"`
	Keywords  []string `yaml:"keywords"`
	Wait      bool     `yaml:"wait"`
	// Actions   []*Point `yaml:"actions"`
}

func (l *Location) String() string {
	return fmt.Sprintf("idloc: %v", l.Key)
}

type Action struct {
	Name      string `yaml:"name"`
    Start []string `yaml:"startloc,omitempty"`
	Steps     []Step `yaml:"steps"`
    Conditional []Step `yaml:"conditional,omitempty"`
    FinlocId string `yaml:"final,omitempty"`
}

type Step struct {
	Grid        string `yaml:"grid"`
	Delay       int    `yaml:"delay,omitempty"`
	Skiplocheck bool   `yaml:"skiplocheck"`
	Wait        bool   `yaml:"wait,omitempty"`
    Loc         []string `yaml:"loc,omitempty"`
}

// Position on Grid
func (l *Location) Position() *adb.Point {
	return cutgrid(l.Grid)
}

func (s *Step) Target() *adb.Point {
	return cutgrid(s.Grid)
}

func (s *Action) ConditionalStep(locid string) Step {
    for _, step := range s.Conditional {
        if slices.Contains(step.Loc, locid) {
            return step
        }
    }
    log.Errorf("%v:%v",locid, ErrStepNotFound.Error() )
    panic(ErrStepNotFound)
}

func init() {
	log = Logger()
	e := createDirStructure()
	f, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE, 0644)
	log.SetLevel(logrus.TraceLevel)
	log.SetOutput(f)
	if e != nil {
		panic(e)
	}
	OcrConf = loadOcr()
	repository.DbInit(func(x string) string {
		return filepath.Join(dbdir, x)
	})
}
func Parse(s string, out interface{}) {
	f, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, out)
	if err != nil {
		log.Fatalf("MARSHAL WASTED: %v", err)
	}
	log.Tracef("MARSHALLED: %v\n\n", out)
}

func UserInput(desc, def string) string {
	//	reader := bufio.NewReader(os.Stdin)
	_ = bufio.NewReader(os.Stdin)
	text := "0"
	color.HiCyan(desc)
	color.HiRed("---------------------")
	fmt.Printf("[default:%v]: ", color.HiGreenString(def))
	//	text, _ := reader.ReadString('\n')
	//	text, _ := reader.ReadString('\n')
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
		fmt.Printf("\nerr:%v\nduring run:%v", e, "intconv")
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

func loadOcr() *ocrConfig {
	res := &ocrConfig{}
	Parse(ocrsettings, res)
	return res
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
			TimestampFormat:           time.Stamp},
		Level: logrus.InfoLevel,
	}

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

	if !filepath.IsAbs(f) {
		return filepath.Join(a, f)
	}
	return filepath.Join(a, filepath.Base(f))
}
func createDirStructure() error {
	//	usr := SafeEnv("USERPROFILE")
	//	wd := filepath.Join(usr, adbdir)
	wd, e := os.Getwd()
	rootd := filepath.Join(wd, localdir)
	dbpath := filepath.Join(wd, localdir, dbdir)
	imgpath := filepath.Join(wd, localdir, tmpdir)
	cfgpath := filepath.Join(wd, localdir, cfgdir)

	truncateDir(imgpath)

	e = os.MkdirAll(cfgpath, os.ModeDir)
	e = os.MkdirAll(dbpath, os.ModeDir)
	e = os.MkdirAll(imgpath, os.ModeDir)

	if e != nil {
		log.Errorf("%v", e)
		return ErrWorkDirFail
	}
	e = os.Chdir(rootd)

	if e == nil {
		pwd, _ := os.Getwd()
		fmt.Printf("\ninit: success; pwd: %v\n\n", pwd)
	}
	//	fmt.Printf("init: fail; error: %v\npwd will be used: %v\n\n", e.Error())
	return e
}

func SafeEnv(name string) string {
	if p, ok := os.LookupEnv(name); ok {
		return p
	}
	log.Warnf("Env [%v] not found!", name)
	return ""
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
