package image

import (
	"errors"
	"fmt"
	c "mugowalker/backend/cfg"
	"mugowalker/backend/localstore"
)

var (
	ErrOptimizeImg = errors.New("error during handling optimize image proccess")
	ErrUsedAllPsm  = errors.New("used all psm modes on this profile")
	prefix         = "ghost-tesseract-out"
)

type ImageProfile struct {
	original   string
	prepared   string
	psm        int
	ignored    []string
	recognized []*ScreenWord
}

func (ip *ImageProfile) TesseractResult() []*ScreenWord {
	f := localstore.RandPostfix(prefix)
	e := runTesseract(ip.prepared, f, customPsm(ip.psm)...)
	if e != nil {
		log("Tessereact mailfunc")
	}
	ip.recognized = UnmarshalAlto(f).parse(ip.ignored)

	log(fmt.Sprintf("Words OnScreen: %v\n	Recognized: %v", len(ip.recognized), almoResultsStringer(ip.recognized, ip.psm)))

	return ip.recognized
}

func (ip *ImageProfile) Redo() []*ScreenWord {
	ip.psm = 3
	return ip.TesseractResult()
}

func customPsm(n int) []string {
	return []string{"--psm", fmt.Sprint(n), "-c", "tessedit_create_alto=1", "quiet"}
}

func runTesseract(in, out string, args ...string) error {
	pArgs := append([]string{in, out}, args...)
	return c.RunCmd(TesseractCmd, pArgs)
}
