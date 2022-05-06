package datman

import (
	"errors"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"worker/navi"

	log "github.com/sirupsen/logrus"
)

type FsMan struct {
	basedir string
	datadir string
	locdir  string
}

func (fd *FsMan) OpenPng(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("File: fail to open img")
		return nil
	}

	defer file.Close()
	pngf, err := png.Decode(file)
	if err != nil {
		log.Printf("File: decode error, its not png")
		return nil
	}

	return pngf
}

func (fd *FsMan) SaveImg(fname string, img image.Image) error {
	f, err := os.Create(fname)
	err = png.Encode(f, img)
	f.Close()
	return err
}

//TODO: I thibk this is navi fucntions
func (fd *FsMan) Candidate(loc *navi.Location, img image.Image) {
	dir := filepath.Join(fd.locdir, loc.Name)
	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil && errors.Is(err, fs.ErrExist) {
		log.Fatalf("FSMan: Cannot create folder %v, error: %v", loc.Name, err)
	}
	nn := filepath.Join(dir, "cadidate_"+loc.Name+".png")
	fd.SaveImg(nn, img)
}

// //TODO: generalize?
// func (fd *FsMan) SaveLocation(loc *navi.Location) error {
// 	pulls, err := fd.Pulled()
// 	for _, v := range pulls {
// 		if strings.Contains(v, loc.Name) {
// 			nn := filepath.Join(fd.locdir, v)
// 			err = os.Rename(v, nn)
// 			log.Tracef("SAVED To loc:  --> %v |||  %v ||| %v", v, nn, err)
// 		}
// 	}
// 	for k, v := range loc.Areas {
// 		name := loc.Name + "_" + k + "_entry.png"
// 		nn := filepath.Join(fd.locdir, name)
// 		fd.SaveImg(nn, v)
// 	}
// 	return err
// }

//TODO: refactor, make pulled in a more general way
func (fd *FsMan) LocEtalons(locname string) (locImgs []image.Image, err error) {
	locpath := filepath.Join(fd.locdir, locname)
	locsFolder, err := os.Open(locpath)
	if err != nil {
		return nil, err
	}
	locFiles, err := locsFolder.Readdirnames(0)
	for _, v := range locFiles {
		if strings.Contains(v, locname) {
			fp := filepath.Join(locpath, v)
			img := fd.OpenPng(fp)
			locImgs = append(locImgs, img)
		}
	}
	if locImgs == nil {
		err = errors.New("No Etalons founded!")
	}
	return
}

func (fd *FsMan) SaveData(name string) error {
	panic("TO IMPLEMENT")
}

func (fm *FsMan) Pulled() ([]string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir, _ := os.Open(pwd)

	filenames, err := dir.Readdirnames(0)
	log.Tracef("FSMan: Pulled --> %v", filenames)
	return filenames, err
}

func isExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return os.IsExist(err)
	}
	return false
}
