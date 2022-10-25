package ocr

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func OCRFields(s string) []string {
	res := strings.Fields(s)
	var filtered []string
	for _, v := range res {
		if len(v) > 2 {
			filtered = append(filtered, v)
		}
	}
	return filtered
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

	outfile, err := ioutil.TempFile("", "ghost-tesseract-out-")
	defer outfile.Close()
	if err != nil {
		return "", err
	}

	if err = runOcr(grayImg.Name(), outfile.Name()); err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(outfile.Name() + ".txt")
	if err != nil {
		return "", err
	}
	return formatStr(strings.TrimSpace(string(bytes))), nil
}

func Text(img string) string {
	t, err := ReadTextFromFile(img)
	if err != nil {
		fmt.Printf("OCR Error: %v", err.Error())
	}
	return t
}

func formatStr(in string) string {
	res := strings.Split(in, "\n")
	return strings.Join(res, " ")
}
