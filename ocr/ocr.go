package ocr

import (
	"os"
	"time"

	c "worker/cfg"

	"github.com/sirupsen/logrus"
)

// AlmoResult parsed xml Almo
type AlmoResult struct {
	Linechars string
	X, Y      int
	LineNo    int
	// Psm       int
}

var (
	user *c.Profile
	log  *logrus.Logger
)

var (
	send func(string, string)
	z    = func(arr []AlmoResult, psm int) string {
		var s string
		s = c.Red(c.F("	↓	|> PSM %v <|	↓	\n", psm))
		line := 0

		for i, elem := range arr {

			if elem.LineNo == line {
				s += c.Cyan(c.F("{idx:%d}%s ", i, elem))
			} else {
				line = elem.LineNo
				s += c.Cyan(c.F("\n{idx:%d}%s ", i, elem))
			}
		}
		s += "\n\n"

		return s
	}
)

func (a AlmoResult) String() string {
	return c.F("[%2d|%4dx%4d <| %s|", a.LineNo, a.X, a.Y, a.Linechars)
}
func init() {
	// Fallback to searching on PATH.
	user = c.ActiveUser()
	log = c.Logger()
}

// ExtractText prepare and extract text from img
func ExtractText(img string) *ImageProfile {
	// defer  timeTrack(time.Now(), "AltOcr")
	ip := &ImageProfile{
		original:   img,
		recognized: make([]AlmoResult, 0),
	}

	imgPrep := PrepareForRecognize(ip)
	log.Debug(c.Red("Optimized img -> "), c.Cyan(imgPrep))
	// s := ip.Tesseract(1)
	// log.Debug(c.F("Words Onscr: %v\n	Ocred: %v", c.Cyan(len(s)), z(s, 1)))

	// }

	return ip
}

func tmpFile() (*os.File, error) {
	outfile, err := os.CreateTemp(c.TempFile(""), "ghost-tesseract-out-")
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
	return c.F("%v\n\r", c.TTrack("\r[%s] %s", name, elapsed.Round(time.Millisecond)))
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
