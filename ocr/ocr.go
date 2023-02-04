package ocr

import (
	"fmt"
	"image"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"worker/cfg"
)

type AltoResult struct {
	Linechars string
	X, Y      int
	LineNo    int
	// Psm       int
}

var (
	tesser   string
	user     *cfg.Profile
	psm      = []int{1, 3, 4, 6, 8, 11, 12}
	altoargs = []string{"--psm", "3", "-c", "tessedit_create_alto=1", "hoot", "quiet"}
	log      *logrus.Logger
)

var (
	send             func(string, string)
	red, green, cyan func(...interface{}) string
	fo               = fmt.Sprintf
	// green = f
	z = func(arr []AltoResult, psm int) string {
		var s string
		s = red(fo("	↓	|> PSM %v <|	↓	\n", psm))
		line := 0

		for i, elem := range arr {

			if elem.LineNo == line {
				s += cyan(fo("{idx:%d}%s ", i, elem))
			} else {
				line = elem.LineNo
				s += cyan(fo("\n{idx:%d}%s ", i, elem))
			}
		}
		s += "\n\n"

		return s
	}
)

type Result struct {
	raw    string
	fields []string
}

func (or Result) String() string {
	return or.raw
}

func (or Result) Fields() []string {
	return or.fields
}

func (or Result) Regex(r string) (res []uint) {
	re := regexp.MustCompile(r)
	for _, v := range re.FindStringSubmatch(or.raw) {
		i, err := strconv.ParseUint(v, 10, 32)
		if err == nil {
			res = append(res, uint(i))
		}
	}
	return
}

func (or Result) Intersect(k []string) (r []string) {
	for _, v := range k {
		if slices.Contains(or.fields, v) {
			r = append(r, v)
		}
	}
	return r
}
func (a AltoResult) String() string {
	return fmt.Sprintf("[%2d|%4dx%4d <| %s|", a.LineNo, a.X, a.Y, a.Linechars)
}
func init() {
	// Fallback to searching on PATH.
	tesser = cfg.LookupPath("tesseract")
	user = cfg.ActiveUser()
	log = cfg.Logger()
	red = color.New(color.FgHiRed).SprintFunc()
	green = color.New(color.FgHiGreen).SprintFunc()
	cyan = color.New(color.FgHiCyan).SprintFunc()
}

func TextExtractAlto(img string) []AltoResult {
	// defer  timeTrack(time.Now(), "AltOcr")
	resu := make([]AltoResult, 0)
	imgPrep := AltOptimize(img)
	for _, v := range psm {
		f, _ := tmpFile()
		tessAlto(imgPrep, f.Name(), customPsm(v)...)
		s := UnmarshalAlto(f.Name()).parse()
		resu = append(resu, s...)
		// resu[v] = r

		log.Debug(fo("Words Onscr: %v\n	Ocred: %v", cyan(len(s)), z(s, v)))

	}

	return unique(resu)
}

func (a Alto) parse() []AltoResult {
	var res []AltoResult
	res = make([]AltoResult, 0)
	tl := a.Layout.Page.PrintSpace.ComposedBlock.TextBlock.TextLine
	for i, line := range tl {
		for _, v := range line.String {
			if len(v.CONTENT) > 3 || slices.Contains(user.Exceptions, v.CONTENT) || strings.ContainsAny(v.CONTENT, "0123456789") {
				res = append(res, AltoResult{Linechars: v.CONTENT, X: cfg.ToInt(v.HPOS), Y: cfg.ToInt(v.VPOS), LineNo: i})
			}
		}
	}
	return res
}

func RegionText(img string, topleft, size image.Point) Result {
	// defer timeTrack(time.Now(), "\nRegionText")
	cropedregion := Concat(img, topleft, size)
	prep := OptimizeForOCR(cropedregion)
	r, e := recognize(prep)
	if e != nil {
		log.Errorf("RegionText fails: %v", e)
	}
	return r
}

func TextExtract(img string) Result {
	// defer timeTrack(time.Now(), "RegularOcr")
	imgPrep := OptimizeForOCR(img)
	t, _ := recognize(imgPrep)
	return t
}

// recognize text on a given img
func recognize(img string) (Result, error) {
	f, _ := tmpFile()
	runOcr(img, f.Name())
	raw, e := readTmp(f.Name() + ".txt")
	r := Result{
		raw: formatStr(strings.TrimSpace(string(raw))),
	}
	log.Tracef("Raw OCR: %s", raw)
	r.fields = cleanText(r.raw)
	return r, e
}

func cleanText(s string) []string {
	res := strings.Fields(s)
	var filtered []string
	for _, v := range res {
		if len(v) > 3 || strings.ContainsAny(v, "01234356789") || slices.Contains(user.Exceptions, v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func formatStr(in string) string {
	res := strings.Split(in, "\n")
	return strings.Join(res, " ")
}

func tmpFile() (*os.File, error) {
	outfile, err := os.CreateTemp("", "ghost-tesseract-out-")
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
	c := color.New(color.BgHiBlue, color.FgCyan, color.Underline, color.Bold).SprintfFunc()
	elapsed := time.Since(start)
	return fmt.Sprintf("%v\n\r", c("\r[%s] %s", name, elapsed.Round(time.Millisecond)))
}

func unique(sample []AltoResult) []AltoResult {
	var unique []AltoResult
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
