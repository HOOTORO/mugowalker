package ocr

import (
	"fmt"
	"os"
	"time"

	c "mugowalker/backend/cfg"
	"mugowalker/backend/settings"
)

type engine struct {
	*settings.Settings
}

// AlmoResult parsed xml Almo
type AlmoResult struct {
	Linechars string
	X, Y      int
	LineNo    int
	// Psm       int
}

var (
	log func(string, string)
)

var almoResultsStringer = func(arr []AlmoResult, psm int) string {
	var s string
	s = fmt.Sprintf("	↓	|> PSM %v <|	↓	\n", psm)
	line := 0
	s += "Ln# 0 -> "
	for _, elem := range arr {

		if elem.LineNo == line {
			// "#%2s|>%30s" c.Cyan(i),
			s += fmt.Sprintf("%-48s", elem.String())
		} else {
			line = elem.LineNo
			// log.Debugf("Len S %d", len(elem.String()))
			s += fmt.Sprintf("\nLn#%11s -> %-48s", elem.LineNo, elem.String())
		}
	}
	s += "\n\n"

	return s
}

func (a AlmoResult) String() string {
	return fmt.Sprintf("%s [%s]", fmt.Sprintf("%13vx%-13v", a.X, a.Y), c.Shortener(a.Linechars, 7))
}
func Engine(c *settings.Settings, usr string) *engine {
	// log = c.Log
	// Fallback to searching on PATH.
	return &engine{c}

}

// ExtractText prepare and extract text from img
func (en *engine) ExtractText(img string) *ImageProfile {
	// defer  timeTrack(time.Now(), "AltOcr")
	ip := &ImageProfile{
		original:   img,
		recognized: make([]AlmoResult, 0),
		psm:        en.Tesseract.Psm,
		ignored:    en.IgnoredWords,
	}

	e := PrepareForRecognize(ip, en.Tesseract.Psm, en.Tesseract.Args)
	if e != nil {
		// en.Log(settings.ERR, "IMAGE NOT PREPARED")
	}
	// en.Log(settings.TRACE, "Optimized img -> "+ip.prepared)

	return ip
}

func tmpFile() (*os.File, error) {
	outfile, err := os.CreateTemp("wd/temp", "ghost-tesseract-out-")
	if err != nil {
		return nil, err
	}
	defer outfile.Close()
	return outfile, nil
}

func readTmp(fname string) ([]byte, error) {
	bytes, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func timeTrack(start time.Time, name string) string {
	elapsed := time.Since(start)
	return fmt.Sprintf("%v\n\r", fmt.Sprintf("\r[%s] %s", name, elapsed.Round(time.Millisecond)))
}

func unique(sample []AlmoResult) []AlmoResult {
	var unique []AlmoResult
	type key struct {
		value1 string
		val2   int
	}
	m := make(map[key]int)
	for _, v := range sample {
		k := key{v.Linechars, v.Y}
		if i, ok := m[k]; ok {
			// Overwrite previous value per requirement in
			// question to keep last matching value.
			unique[i] = v
		} else {
			// Unique key found. Record position and collect
			// in result.
			m[k] = len(unique)
			unique = append(unique, v)
		}
	}
	return unique
}
