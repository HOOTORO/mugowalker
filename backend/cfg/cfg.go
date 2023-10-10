package cfg

import (
	"errors"
	"fmt"
	"mugowalker/backend/localstore"

	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

var temp = localstore.TempDir()

var (
	ErrWorkDirFail        = errors.New("working directories wasn't created. Exit")
	ErrRequiredProgram404 = errors.New("missing some of required soft")
	ErrLoadInitConf       = errors.New("load sysvars")
)

var (
	log *logrus.Logger = logrus.New()
)

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

func RunCmd(command string, args []string) error {
	pt := LookupPath(command)
	log.Trace("  ↓   RunCMD   ↓ \n", pt, "\n", args)
	cmd := exec.Command(pt, args...)
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

func Shortener(str string, n int) string {
	if len(str) > n+3 {
		return str[:n] + "..."
	}
	return str
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
