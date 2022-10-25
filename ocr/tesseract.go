package ocr

import (
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/harrydb/go/img/grayscale"
	// "github.com/otiai10/gosseract/v2"
)

var tesser string

func init() {
	// Fallback to searching on PATH.
	if p, err := exec.LookPath("tesseract"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			tesser = p
		}
	}
}

func covertGrayscale(r io.Reader) (*os.File, error) {
	src, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	gray := grayscale.Convert(src, grayscale.ToGrayLuminance)
	grayImg, err := ioutil.TempFile("", "tesseract-gray-")
	defer grayImg.Close()
	if err != nil {
		return nil, err
	}

	err = jpeg.Encode(grayImg, gray, &jpeg.Options{Quality: 100})
	// orig quality 80
	if err != nil {
		return nil, err
	}

	return grayImg, nil
}

func runOcr(in string, out string) error {
	ocr, err := filepath.Abs(tesser)
	if err != nil {
		return err
	}

	cmd := exec.Command(ocr, in, out, "--psm", "12")
	// uncomment for ocr log
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}

// func TextGOS(img string) string {
// 	client := gosseract.NewClient()

// 	defer client.Close()
// 	client.SetImage(img)
// 	text, _ := client.Text()
// 	// fmt.Println(text)
// 	// Hello, World!
// 	return text
// }
