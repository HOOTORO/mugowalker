package datman

import (
	"errors"
	"image"
	"image/png"
	"io/fs"
	"log"
	"os"
)

//TODO: rework this ugly ...
type DataManager interface {
	OpenPng(string) image.Image
	SaveImg(string, image.Image) error
}

type FsMan struct {
	name string
}

func NewFM(name string) *FsMan {
	err := os.MkdirAll(name, os.ModeDir)
	if err != nil && errors.Is(err, fs.ErrExist) {
		log.Fatalf("FSMan: Cannot create folder %v, error: %v", name, err)
		return nil
	}
	SetWD(name)
	return &FsMan{name: name}
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

func SetWD(path string) {
	err := os.Chdir(path)
	if err != nil {
		log.Panicf("cant change wd: %s", err.Error())
	}
}
