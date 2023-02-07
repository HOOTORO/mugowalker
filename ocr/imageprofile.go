package ocr

import (
	"errors"
	"strings"
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
	used_psms  string
	recognized []AlmoResult
}

func (ip *ImageProfile) Tesseract(psmin int) []AlmoResult {
	switch psmin {
	case 1, 3, 4, 6, 11, 12:
		f, _ := tmpFile()
		e := ActivateTesseract(ip.prepared, f.Name(), customPsm(psmin)...)
		if e != nil {
			log.Errorf("Tessereact mailfunc")
		}
		ip.recognized = unique(UnmarshalAlto(f.Name()).parse())
		ip.used_psms += c.F("%v", psmin)
		z(ip.recognized, psmin)
		return ip.recognized

	default:
		log.Warnf("Provide correct PSM num (%v)", psm)
		return nil
	}
}

func (ip *ImageProfile) Result() []AlmoResult {
	if len(ip.recognized) > 0 {
		return ip.recognized
	} else {
		return ip.Tesseract(6)
	}
}

func (ip *ImageProfile) TryAgain() []AlmoResult {

	newpsm, e := unusedpsm(ip.used_psms)
	if errors.Is(e, ErrUsedAllPsm) {
		blaine.args = MagickArgs()
		ip.used_psms = ""
		PrepareForRecognize(ip)
		return ip.Tesseract(1)
	} else {
		return ip.Tesseract(newpsm)
		// z(ip.recognized, newpsm)
	}
}

func unusedpsm(u string) (int, error) {
	for _, p := range psm {
		if !strings.ContainsAny(u, c.F("%v", p)) {
			return p, nil
		}
	}
	return 0, ErrUsedAllPsm
}
