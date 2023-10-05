package localstore

import (
	"os"
	"path"
)

const app = "mugowalker"

type LocalStore struct {
	ConfDir string
	WorkDir string
}

func NewLocalStore() *LocalStore {
	return &LocalStore{ConfDir: path.Join(getConfigHome(), app), WorkDir: path.Join(getTempHome(), app)}
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
	return conf
}

func getTempHome() string {

	cache, err := os.UserCacheDir()
	if err != nil {
		return getUserHome()
	}
	return cache
}

func getUserHome() string {
	home, e := os.UserHomeDir()
	if e != nil {
		return "."
	}
	return home
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
