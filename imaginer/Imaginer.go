package imaginer

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/harrydb/go/img/grayscale"

	// "github.com/otiai10/gosseract/v2"
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

var tesser string

func init() {
	// Fallback to searching on PATH.
	if p, err := exec.LookPath("tesseract"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			tesser = p
		}
	}
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

// func TextGOS(img string) string {
// 	client := gosseract.NewClient()

// 	defer client.Close()
// 	client.SetImage(img)
// 	text, _ := client.Text()
// 	// fmt.Println(text)
// 	// Hello, World!
// 	return text
// }

func Text(img string) string {
	t, err := ReadTextFromFile(img)
	if err != nil {
		fmt.Printf("OCR Error: %v", err.Error())
	}
	return t
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

// ReadTextFromFile read text from the file. It internally calls ReadText after reading the file.
func ReadTextFromFile(f string) (string, error) {
	r, err := os.Open(f)
	if err != nil {
		return "", err
	}
	defer r.Close()

	return ReadText(r)
}

// ReadText read text from the given io.Reader r. It converts to grayscale first before pass it to tesseract.
// It writes grayscale image and output text file to the os.TempFile.
func ReadText(r io.Reader) (string, error) {
	grayImg, err := covertGrayscale(r)

	outfile, err := ioutil.TempFile("", "ghost-tesseract-out-")
	defer outfile.Close()
	if err != nil {
		return "", err
	}

	if err = runOcr(grayImg.Name(), outfile.Name()); err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(outfile.Name() + ".txt")
	if err != nil {
		return "", err
	}
	result := strings.TrimSpace(string(bytes))
	strings.ReplaceAll("\n", result, result)
	return result, nil
}

func runOcr(in string, out string) error {
	ocr, err := filepath.Abs(tesser)
	if err != nil {
		return err
	}

	cmd := exec.Command(ocr, in, out)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func covertGrayscale(r io.Reader) (*os.File, error) {
	src, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	gray := grayscale.Convert(src, grayscale.ToGrayLuminance)
	grayImg, err := ioutil.TempFile("", "tesseract-gray-")
	defer grayImg.Close()
	if err != nil {
		return nil, err
	}

	err = jpeg.Encode(grayImg, gray, &jpeg.Options{Quality: 100})
	// orig quality 80
	if err != nil {
		return nil, err
	}

	return grayImg, nil
}
