package ocr

import (
	"fmt"
	"image"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"worker/imaginer"

	"golang.org/x/exp/slices"
)

var uni string = "“€”\"’^#@™&!~'‘|<$>«»,:¢\\/_;§®‘*~."

type OcrResult struct {
	raw    string
	fields []string
}

func (or OcrResult) String() string {
	return strings.Join(or.fields, " | ")
}

func (or OcrResult) Fields() []string {
	return or.fields
}

func (or OcrResult) Regex(r string) (res []int) {
	re := regexp.MustCompile(r)
	for _, v := range re.FindStringSubmatch(or.raw) {
		i, err := strconv.Atoi(v)
		if err == nil {
			res = append(res, i)
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

// ReadTextFromFile read text from the file. It internally calls ReadText after reading the file.
func ReadTextFromFile(f string) (string, error) {
	r, err := os.Open(f)
	if err != nil {
		return "", err
	}
	defer r.Close()

	return ReadText(r)
}

// ReadText read text from the given io.Reader r. It converts to grayscale first before pass it to tesseract.
// It writes grayscale image and output text file to the os.TempFile.
func ReadText(r io.Reader) (string, error) {
	grayImg, err := covertGrayscale(r)
	outfile, _ := tmpFile()

	if err = runOcr(grayImg.Name(), outfile.Name()); err != nil {
		return "", err
	}
	ocredBytes, _ := readTmp(outfile.Name())
	return formatStr(strings.TrimSpace(string(ocredBytes))), nil
}

func Text(img string) string {
	t, err := ReadTextFromFile(img)
	if err != nil {
		fmt.Printf("OCR Error: %v", err.Error())
	}
	return t
}

func KeywordHits(kw, ocr []string) int {
	res := 0
	for _, word := range kw {
		if slices.Contains(ocr, word) {
			res++
		}
	}
	return res
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

func ImprovedTextExtract(img string) OcrResult {
	defer timeTrack(time.Now(), "ImprovedTextExtract")
	var result OcrResult
	imgPrep := OptimizeForOCR(img)
	images := imaginer.GridCrop(imgPrep)

	for _, v := range images {
		r, _ := recognize(v)
		result.raw += r.raw
		result.fields = append(result.fields, r.Fields()...)
	}
	return result
}

func TextExtract(img string) OcrResult {
	defer timeTrack(time.Now(), "RegularOcr")
	imgPrep := OptimizeForOCR(img)
	t, _ := recognize(imgPrep)
	return t
}

// recognize text on a given img
func recognize(img string) (OcrResult, error) {
	f, _ := tmpFile()
	e := runOcr(img, f.Name())
	raw, e := readTmp(f.Name())
	r := OcrResult{
		raw: formatStr(strings.TrimSpace(string(raw))),
	}
	log.Tracef("Raw OCR: %s", raw)
	r.fields = CleanText(r.raw)
	return r, e
}

func CleanText(s string) []string {
	res := strings.Fields(s)
	var filtered []string
	for _, v := range res {
		trimmed := strings.Trim(v, uni)
		if len(trimmed) > 3 || strings.ContainsAny(trimmed, "01234356789") {
			filtered = append(filtered, trimmed)
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
	bytes, err := os.ReadFile(fname + ".txt")
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func timeTrack(start time.Time, name string) {
	c := color.New(color.BgHiWhite, color.FgHiRed, color.Underline, color.Bold).SprintfFunc()
	elapsed := time.Since(start)
	fmt.Printf("%v\n", c("[%s] %s", name, elapsed.Round(time.Millisecond)))
}
