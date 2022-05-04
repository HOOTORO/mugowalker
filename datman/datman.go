package datman

import (
	"errors"
	"image"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const (
	data        = "rawdata"
	locationdir = "Locations"
)

//TODO: rework this ugly ...
type DataManager interface {
	OpenPng(string) image.Image
	SaveImg(string, image.Image) error
	SaveLoc(name string, Loca image.Image) error
}

type FsMan struct {
	basedir string
	datadir string
	locdir  string
}

func NewFM(name string) *FsMan {
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
	f, _ := os.Create(fname)
	return png.Encode(f, img)
}

func (fd *FsMan) SaveLoc(name string, Loca image.Image) error {
	nn := filepath.Join(fd.locdir, name)
	err := fd.SaveImg(nn, Loca)
	return err
}

func (fd *FsMan) SaveData(name string) error {
	panic("TO IMPLEMENT")
}

func SetWD(path string) {
	err := os.Chdir(path)
	if err != nil {
		log.Panicf("cant change wd: %s", err.Error())
	}
}
