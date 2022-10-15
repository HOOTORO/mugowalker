package datman

import (
	"errors"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
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

func (fd *FsMan) Candidate(loc *navi.Location, img image.Image) {
	nn := "cadidate_" + loc.Name + ".png"
	fd.saveToLocFolder(nn, loc, img)
}

func (fd *FsMan) Unknownplace(loc *navi.Location, img image.Image, clickSector navi.TPoint) {
	nn := "unknown_sector" + strconv.Itoa(clickSector.X) + "-" + strconv.Itoa(clickSector.Y) + ".png"
	clkarea := loc.CutSector(loc.Etalons[0], clickSector.X, clickSector.Y, 60, "unknown")
	fd.saveToLocFolder(nn, loc, img)
	fd.saveToLocFolder("secotr_unkno.png", loc, clkarea)

}

func (fd *FsMan) saveToLocFolder(fname string, loc *navi.Location, img image.Image) {
	dir := filepath.Join(fd.locdir, loc.Name)
	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil && errors.Is(err, fs.ErrExist) {
		log.Fatalf("FSMan: Cannot create folder %v, error: %v", loc.Name, err)
	}
	nn := filepath.Join(dir, fname)
	fd.SaveImg(nn, img)
}

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
