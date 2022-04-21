package fshelp

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

type Fshelper interface {
	Copy(string, string)
	OpenImg(string)
	Write(interface{})
}

func OpenImg(path string) image.Image {
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

func SaveAsPng(fname string, img image.Image) {
	f, _ := os.Create(fname)
	png.Encode(f, img)
}

func CreateFolder(fname string) error {
	err := os.MkdirAll(fname, os.ModeDir)
	if err != nil && errors.Is(err, fs.ErrExist) {
		return nil
	}
	return err
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err

}

func ReadFile(fname string) ([]byte, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Printf("File read error: %s", err.Error())
	}
	return data, err
}

func SetWD(path string) {
	//NOTE: return after clog refactor
	err := os.Chdir(path)
	if err != nil {
		log.Panicf("cant change wd: %s", err.Error())
	}
}

//Creates file, if not exist, Truncates if exist
func WriteFile(fname string, data []byte) {
	err := ioutil.WriteFile(fname, data, 0777)
	if err != nil {
		panic("Unable to write data into the file" + err.Error())
	}

}
