package ocr

import c "worker/cfg"

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
		return ip.recognized
	default:
		log.Warnf("Provide correct PSM num (%v)", psm)
		return nil
	}
}
