package ocr

import (
	"os/exec"

	"worker/cfg"
	"worker/imaginer"

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

func runOcr(in string, out string) error {
	args := append([]string{in, out}, tessAgrs...)
	cmd := exec.Command(tesser, args...)
    log.Tracef("cmd tess : %v\n", cmd.String())
	// uncomment for ocr log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}


//func tessAlto