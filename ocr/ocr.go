package ocr

import (
	"fmt"
	"image"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"worker/cfg"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

type AltoResult struct {
	Linechars string
	X, Y      int
}

var send func(string, string)

func (a AltoResult) String() string {
	return fmt.Sprintf("%vx%v | |> %v <|", a.X, a.Y, a.Linechars)
}

func TextExtractAlto(img string) []AltoResult {
	// defer  timeTrack(time.Now(), "AltOcr")
	imgPrep := AltOptimize(img)
	f, _ := tmpFile()
	tessAlto(imgPrep, f.Name())
	s := UnmarshalAlto(f.Name())
	return s.parse()
}

func (a Alto) parse() []AltoResult {
	var res []AltoResult
	res = make([]AltoResult, 0)
	//    fmt.Printf("%v", pass("%v",))
	tl := a.Layout.Page.PrintSpace.ComposedBlock.TextBlock.TextLine
	for _, line := range tl {
		for _, v := range line.String {
			if len(v.CONTENT) > 3 || slices.Contains(user.Exceptions, v.CONTENT) {
				res = append(res, AltoResult{Linechars: v.CONTENT, X: cfg.ToInt(v.HPOS), Y: cfg.ToInt(v.VPOS)})
			}
		}
	}
	return res
}

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
