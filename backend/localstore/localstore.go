package localstore

import (
	"fmt"
	"math/rand"
	"os"
	"path"
)

const app = "mugowalker"

type LocalStore struct {
	ConfDir string
	WorkDir string
}

func NewLocalStore() *LocalStore {
	return &LocalStore{ConfDir: getConfigHome(), WorkDir: getTempHome()}
}

func TempDir() string {
	return getTempHome()
}
func TempFile() (*os.File, error) {
	outfile, err := os.CreateTemp(TempDir(), "tempo-f-")
	if err != nil {
		return nil, err
	}
	defer outfile.Close()
	return outfile, nil
}
func RandPostfix(str string) string {
	num := rand.Intn(10000000000)
	filename := fmt.Sprintf("%s-%d", str, num)
	return path.Join(TempDir(), filename)
}

func ReadTempFile(fname string) ([]byte, error) {
	bytes, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (ls *LocalStore) Store(data []byte, filename string, isConf bool) error {
	tg := ls.target(isConf)
	p := path.Join(tg, filename)
	if err := ensureDirExists(tg); err != nil {
		return err
	}
	if err := os.WriteFile(p, data, 0777); err != nil {
		return err
	}
	return nil
}
func (ls *LocalStore) Load(filename string, isConf bool) ([]byte, error) {
	p := path.Join(ls.target(isConf), filename)
	d, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return d, err
}

func getConfigHome() string {
	conf, e := os.UserConfigDir()
	if e != nil {
		return getUserHome()
	}
	appFolder := path.Join(conf, app)
	ensureDirExists(appFolder)
	return appFolder
}

func getTempHome() string {

	cache, err := os.UserCacheDir()
	if err != nil {
		return getUserHome()
	}
	appFolder := path.Join(cache, app)
	ensureDirExists(appFolder)
	return appFolder
}

func getUserHome() string {
	home, e := os.UserHomeDir()
	if e != nil {
		return "."
	}
	appFolder := path.Join(home, app)
	ensureDirExists(appFolder)
	return appFolder
}

func (ls *LocalStore) target(b bool) string {
	if b {
		return ls.ConfDir
	} else {
		return ls.WorkDir
	}
}

func ensureDirExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.Mkdir(path, 0777); err != nil {
			return err
		}
	}
	return nil
}
