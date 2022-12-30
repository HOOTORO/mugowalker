package ocr

import (
    "image/jpeg"
	"image/png"
	"io"
	"os"
	"os/exec"
	"worker/cfg"
    "worker/imaginer"

    "github.com/sirupsen/logrus"
    "github.com/harrydb/go/img/grayscale"
	// "github.com/otiai10/gosseract/v2"
)
var tesser string
var tessAgrs []string

var log *logrus.Logger
func init() {
	// Fallback to searching on PATH.
	tesser = cfg.LookupPath("tesseract")
	tessAgrs = cfg.OcrConf.Tesseract
    log = cfg.Logger()
}
func OptimizeForOCR(f string) string {
    res, _ := imaginer.Magick(f, cfg.OcrConf.Imagick...)
    return res
}
func covertGrayscale(r io.Reader) (*os.File, error) {
	src, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	gray := grayscale.Convert(src, grayscale.ToGrayLuminance)
	grayImg, err := os.CreateTemp("", "tesseract-gray-")
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

func runOcr(in string, out string) error {

    args := append([]string{in,out}, tessAgrs...)
    log.Tracef("Tesseract args -> %v", args)
	cmd := exec.Command(tesser, args...)
	// uncomment for ocr log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}
