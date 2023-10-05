package ocr

import (
	"errors"
	"fmt"
	"mugowalker/backend/settings"
	"time"
)

var (
	ErrUsedAllPsm = errors.New("used all psm modes on this profile")
)

type ImageProfile struct {
	original   string
	prepArgs   []string
	prepared   string
	psm        int
	ignored    []string
	recognized []AlmoResult
}

const defpsm = 6

func (ip *ImageProfile) NewResults() []AlmoResult {

	f, _ := tmpFile()

	log(settings.TRACE, "SENT TO TESS PREPARED FILE  -> "+ip.prepared)
	e := ActivateTesseract(ip.prepared, f.Name(), customPsm(defpsm)...)
	if e != nil {
		log(settings.ERR, "Tessereact mailfunc")
	}
	ip.recognized = unique(UnmarshalAlto(f.Name()).parse(ip.ignored))

	log(settings.TRACE, fmt.Sprintf("Words Onscr: %v\n	Ocred: %v", len(ip.recognized), almoResultsStringer(ip.recognized, defpsm)))

	return ip.recognized

}

func (ip *ImageProfile) Result() []AlmoResult {
	if len(ip.recognized) > 0 {
		return ip.recognized
	} else {
		return ip.NewResults()
	}
}

func (ip *ImageProfile) TryAgain() []AlmoResult {
	time.Sleep(3 * time.Second)
	PrepareForRecognize(ip, ip.psm, ip.prepArgs)
	return ip.NewResults()

}
