package imaginer

import (
	"fmt"
    "github.com/sirupsen/logrus"
    "image"
	"os/exec"
	"path/filepath"
	"strings"
	"worker/cfg"

	"github.com/vitali-fedulov/images/v2"
)

var magick string
var log *logrus.Logger
const (

	CROP = "-crop"
)

func init() {
	// Fallback to searching on PATH.
    magick =cfg.LookupPath("magick")
    log = cfg.Logger()

}

type Cutter interface {
	Concat(image.Image, int, int)
}

type TextScanner interface {
	ImageText(image.Image) (string, error)
}

type Similizer interface {
	Similarity(image.Image, image.Image) (similar bool, percent int)
}

func Similarity(imgA, imgB image.Image) (similar bool) {
	// Calculate hashes and image sizes.
	hashA, imgSizeA := images.Hash(imgA)
	hashB, imgSizeB := images.Hash(imgB)

	similar = images.Similar(hashA, hashB, imgSizeA, imgSizeB)
	log.Debugf("Are Images similar? --> %v", similar)

	return
}

func OpenImg(fname string) image.Image {
	imgA, err := images.Open(fname)
	if err != nil {
		panic(err)
	}
	if err != nil {
		return nil
	}
	return imgA
}

// func PrepareImg(img string) string {
// 	imagick.Initialize()
// 	defer imagick.Terminate()
// 	dest := "prcsd.png"
// 	mw := imagick.NewMagickWand()
// 	mw.ReadImage(img)
// 	width, height := mw.GetImageWidth(), mw.GetImageHeight()
// 	half := mw.GetImageRegion(0, height/2, int(width), int(height))
// 	half.WriteImage(dest)
// 	half.Destroy()
// 	mw.Destroy()
// 	return dest
// }

func Magick(img string, args ...string) (string, error) {
	out := cfg.ImageDir(img)
	args = append([]string{img}, args...)
	args = append(args, out)
    log.Tracef("Imagick args -> %v", args)
	cmd := exec.Command(magick, args...)

	return out, cmd.Run()
}


func Concat(f string, topleft, bottomright image.Point) string {
	posArg := fmt.Sprintf("%vx%v+%v+%v", bottomright.X, bottomright.Y, topleft.X, topleft.Y)
	res, e := Magick(f, CROP, posArg)
	if e != nil {
		log.Errorf("%v", e)
		return ""
	}
	return res
}

func GridCrop(f string) (crpdImages []string) {
	r, e := Magick(f, cfg.OcrConf.Split...)
    if e != nil {
        log.Errorf("Grid Crop fail -> %v", e)
    }
    origName := strings.TrimRight(filepath.Base(f), filepath.Ext(f))
	for _, file := range cfg.GetImages() {
        if file != r && strings.Contains(file, origName) {
            crpdImages = append(crpdImages,file)
		}
	}
	return crpdImages
}


