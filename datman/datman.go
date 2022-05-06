package datman

import (
	"errors"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"worker/navi"

	log "github.com/sirupsen/logrus"
)

const (
	data        = "rawdata"
	locationdir = "locations"
)

//TODO: rework this ugly ...
type DataManager interface {
	OpenPng(string) image.Image
	SaveImg(string, image.Image) error
	// SaveLocation(loc *navi.Location) error
	Candidate(loc *navi.Location, img image.Image)
	LocEtalons(locname string) (locImgs []image.Image, err error)
}

func NewFM(name string) *FsMan {
	if !filepath.IsAbs(name) {
		name, _ = filepath.Abs(name)
	}
	err := os.MkdirAll(name, os.ModeDir)
	if err != nil && errors.Is(err, fs.ErrExist) {
		log.Fatalf("FSMan: Cannot create folder %v, error: %v", name, err)
		return nil
	}
	datadir := filepath.Join(name, data)
	err = os.MkdirAll(datadir, os.ModeDir)
	if err != nil && errors.Is(err, fs.ErrExist) {
		log.Fatalf("FSMan: Cannot create folder %v, error: %v", name, err)
		return nil
	}
	locdi := filepath.Join(name, locationdir)
	err = os.MkdirAll(locdi, os.ModeDir)
	if err != nil && errors.Is(err, fs.ErrExist) {
		log.Fatalf("FSMan: Cannot create folder %v, error: %v", name, err)
		return nil
	}

	SetWD(name)
	return &FsMan{basedir: name, datadir: datadir, locdir: locdi}
}

func SetWD(path string) {
	err := os.Chdir(path)
	if err != nil {
		log.Panicf("cant change wd: %s", err.Error())
	}
}
