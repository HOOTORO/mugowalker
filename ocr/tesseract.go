package ocr

import (
	"fmt"
	"os/exec"
)

func OptimizeForOCR(f string) string {
	res, _ := Magick(f, user.ImagickCfg()...)
	return res
}

func customPsm(n int) []string {
	return []string{"--psm", fmt.Sprint(n), "-c", "tessedit_create_alto=1", "hoot", "quiet"}
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
		tessAgrs = user.TesseractCfg()
	}
	args := append([]string{in, out}, tessAgrs...)
	cmd := exec.Command(tesser, args...)
	log.Tracef("cmd tess : %v\n", cmd.String())
	// uncomment for ocr log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}

func tessAlto(in, out string, args ...string) error {
	args = append([]string{in, out}, args...)
	cmd := exec.Command(tesser, args...)
	log.Tracef("cmd tess : %v\n", cmd.String())

	return cmd.Run()
}
