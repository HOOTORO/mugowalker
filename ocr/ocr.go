package ocr

import (
	"fmt"
	"os"
	"strings"
	"time"

	"worker/cfg"
	"worker/vendor/github.com/fatih/color"
	"worker/vendor/golang.org/x/exp/slices"
)

func OCRFields(s string) []string {
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
