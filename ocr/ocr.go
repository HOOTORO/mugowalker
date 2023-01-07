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

	"worker/imaginer"

	"golang.org/x/exp/slices"
)

type OcrResult struct {
	raw    string
	fields []string
}

func (or OcrResult) String() string {
	return or.raw//strings.Join(or.fields, " | ")
}

func (or OcrResult) Fields() []string {
	return or.fields
}

func (or OcrResult) Regex(r string) (res []uint) {
	re := regexp.MustCompile(r)
	for _, v := range re.FindStringSubmatch(or.raw) {
		i, err := strconv.ParseUint(v, 10, 32)
		if err == nil {
			res = append(res, uint(i))
		}
	}
	return
}

func (or OcrResult) Intersect(k []string) (r []string) {
	for _, v := range k {
		if slices.Contains(or.fields, v) {
			r = append(r, v)
		}
	}
	return r
}

func RegionText(img string, topleft, size image.Point) OcrResult {
	defer timeTrack(time.Now(), "\nRegionText")
	cropedregion := imaginer.Concat(img, topleft, size)
	prep := OptimizeForOCR(cropedregion)
	r, e := recognize(prep)
	if e != nil {
		log.Errorf("RegionText fails: %v", e)
	}
	return r
}

func TextExtract(img string) OcrResult {
	defer timeTrack(time.Now(), "RegularOcr")
	imgPrep := OptimizeForOCR(img)
	t, _ := recognize(imgPrep)
	return t
}

func TextExtractAlto(img string) Alto {
	defer timeTrack(time.Now(), "RegularOcr")
	imgPrep := OptimizeForOCR(img)
    f, _ := tmpFile()
    runOcr(imgPrep, f.Name())
    return UnmarshalAlto(f.Name())

}

// recognize text on a given img
func recognize(img string) (OcrResult, error) {
	f, _ := tmpFile()
	e := runOcr(img, f.Name())
    raw, e := readTmp(f.Name()+ ".txt")
	r := OcrResult{
		raw: formatStr(strings.TrimSpace(string(raw))),
	}
	log.Tracef("Raw OCR: %s", raw)
//	color.HiCyan("Raw OCR: %s", raw)
	r.fields = cleanText(r.raw)
	return r, e
}

func cleanText(s string) []string {
	res := strings.Fields(s)
	var filtered []string
	for _, v := range res {
		if len(v) > 3 || strings.ContainsAny(v, "01234356789") || slices.Contains(cfg.OcrConf.Exceptions, v) {
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
	defer outfile.Close()
	if err != nil {
		return nil, err
	}
	return outfile, nil
}

func readTmp(fname string) ([]byte, error) {
	bytes, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func timeTrack(start time.Time, name string) {
	c := color.New(color.BgHiBlue, color.FgCyan, color.Underline, color.Bold).SprintfFunc()
	elapsed := time.Since(start)
	fmt.Printf("%v\n", c("[%s] %s", name, elapsed.Round(time.Millisecond)))
}
