package ocr

import (
    "image/jpeg"
	"image/png"
	"io"
	"os"
	"os/exec"

	"worker/cfg"
	"worker/imaginer"

	"github.com/harrydb/go/img/grayscale"
	"github.com/sirupsen/logrus"
)

var (
	tesser   string
	tessAgrs []string
)

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
	args := append([]string{in, out}, tessAgrs...)
	cmd := exec.Command(tesser, args...)
    log.Tracef("cmd tess : %v\n", cmd.String())
	// uncomment for ocr log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}
