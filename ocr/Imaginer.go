package imaginer

import (
	"image"

	log "github.com/sirupsen/logrus"
	"github.com/vitali-fedulov/images/v2"
	// "gopkg.in/gographics/imagick.v3/imagick"
)

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
