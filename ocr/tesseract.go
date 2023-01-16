package ocr

import (
	"os/exec"

	"worker/cfg"

	"github.com/sirupsen/logrus"
)

var (
	tesser   string
	user *cfg.Profile

	altoargs = []string{"--psm", "3", "-c", "tessedit_create_alto=1", "quiet"}
	log *logrus.Logger
)

func init() {
	// Fallback to searching on PATH.
	tesser = cfg.LookupPath("tesseract")
	user = cfg.ActiveUser()
	log = cfg.Logger()
}

func OptimizeForOCR(f string) string {
	res, _ := Magick(f, user.Imagick...)
	return res
}

func AltOptimize(f string) string {
	res, _ := Magick(f, user.AltImagick...)
	return res
}

func runOcr(in string, out string) error {
	var tessAgrs []string
	if user.UseAltImagick {
		tessAgrs = user.AltImagick
	} else {
		tessAgrs = user.Tesseract
	}
	args := append([]string{in, out}, tessAgrs...)
	cmd := exec.Command(tesser, args...)
	log.Tracef("cmd tess : %v\n", cmd.String())
	// uncomment for ocr log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}

func tessAlto(in, out string) error {
	args := append([]string{in, out}, altoargs...)
	cmd := exec.Command(tesser, args...)
	log.Tracef("cmd tess : %v\n", cmd.String())

	return cmd.Run()
}
