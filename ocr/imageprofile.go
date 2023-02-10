package ocr

import (
	"errors"
	"time"
	c "worker/cfg"
)

var (
	ErrUsedAllPsm = errors.New("used all psm modes on this profile")
)

type ImageProfile struct {
	original   string
	prepArgs   []string
	prepared   string
	psm        int
	recognized []AlmoResult
}

const defpsm = 6

func (ip *ImageProfile) NewResults() []AlmoResult {

	f, _ := tmpFile()

	log.Trace(c.Cyan("SENT TO TESS PREPARED FILE  -> "), c.Mgt(ip.prepared))
	e := ActivateTesseract(ip.prepared, f.Name(), customPsm(defpsm)...)
	if e != nil {
		log.Errorf("Tessereact mailfunc")
	}
	ip.recognized = unique(UnmarshalAlto(f.Name()).parse())

	log.Debug(c.F("Words Onscr: %v\n	Ocred: %v", c.Cyan(len(ip.recognized)), almoResultsStringer(ip.recognized, defpsm)))

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

	blaine.NewRandArgs()
	time.Sleep(3 * time.Second)
	PrepareForRecognize(ip)
	return ip.NewResults()

}
