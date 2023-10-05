package cfg

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	defaultcfg = "backend/assets/cfg.yml" //"assets/default.yaml"
	game       = "AFK Arena"
	temp       = "wd/temp"
)

type Runnable interface {
	Path() string
	Args() []string
}

var (
	ErrWorkDirFail        = errors.New("working directories wasn't created. Exit")
	ErrRequiredProgram404 = errors.New("missing some of required soft")
	ErrLoadInitConf       = errors.New("load sysvars")
)

var (
	log *logrus.Logger
)

func init() {

	// repository.DbInit(func(x string) string {
	// 	return filepath.Join(sysvars.Db, x)
	// })
}

// type enum interface {
// 	~uint
// 	String() string
// 	Values() []string
// }

// // Deserialize bits to string values
// func Deserialize[T enum](raw T) []string {
// 	var result []string
// 	for i := 0; i < len(raw.Values()); i++ {
// 		if d := T(1 << i); raw&(1<<uint(i)) != 0 {
// 			result = append(result, d.String())
// 		}
// 	}
// 	return result
// }

// Logger for app to use

// GetImages from temp/<appfolder>
func GetImages() []string {
	d, e := os.ReadDir(temp)
	if e != nil {
		panic("get imgs")
	}
	var res []string
	for _, f := range d {
		res = append(res, f.Name())
	}
	return res
}

func RunCmd(r Runnable) error {
	pt := LookupPath(r.Path())
	log.Trace("  ↓   RunCMD   ↓ \n", pt, "\n", r.Args())
	cmd := exec.Command(pt, r.Args()...)
	return cmd.Run()
}

// LookupPath for exe s
func LookupPath(name string) (path string) {
	p, err := exec.LookPath(name)
	if err == nil {
		if p, err = filepath.Abs(p); err == nil {
			return p
		}
	}
	panic(fmt.Sprintf("Required program: %v not found in path\n error: %v", name, err))
}

/*
Helper func
*/

func ToInt(s string) int {
	num, e := strconv.Atoi(s)
	if e != nil {
		log.Warnf("Called.F():%v\nError:%v", "cfg.ToInt", e)
	}
	return num
}

func safeEnv(n string) string {
	str, ok := os.LookupEnv(n)
	if ok {
		return str
	}
	log.Warnf("$Env:%v NOT FOUND, BE AWARE", n)
	return ""
}
func Shortener(str string, n int) string {
	if len(str) > n+3 {
		return str[:n] + "..."
	}
	return str
}
