package ocr

import (
	"errors"
	"fmt"
	c "worker/cfg"
)

var (
	ErrOptimizeImg = errors.New("error during handling optimize image proccess")
)
var (
	psm      = []int{1, 3, 4, 6, 11, 12}
	altoargs = []string{"--psm", "3", "-c", "tessedit_create_alto=1", "hoot", "quiet"}
)

const tessex = "tesseract"

// Tesseract parameters to cfg.RunCmd
type Tesseract struct {
	path string
	args []string
	in   string
	out  string
}

// Path implemetation Runnable
func (t *Tesseract) Path() string {
	return t.path
}

// Args implemetation Runnable
func (t *Tesseract) Args() []string {
	// args := append([]string{t.in, t.out}, t.args...)
	// t.args = args
	return append([]string{t.in, t.out}, t.args...)
}

func (t *Tesseract) SetArgs(in, out string, args ...string) {
	t.in = in
	t.out = out
	t.args = args

	log.Trace("Tessa ARGS --> ", c.Red(tessa.Args()))

}

var tessa = &Tesseract{
	path: tessex,
	args: make([]string, 0),
}

func init() {
	tessa.args = user.TesseractCfg()
}

// PrepareForRecognize alternative optimization
func PrepareForRecognize(f *ImageProfile) error {
	blaine.SetFile(f.original)
	e := c.RunCmd(blaine)
	if e != nil {
		log.Error(ErrOptimizeImg, e)
		return e
	}
	f.prepared = blaine.Prepared()
	return nil
}

func ActivateTesseract(in, out string, args ...string) error {
	tessa.SetArgs(in, out, args...)
	// tessa.out = out
	// tessa.args = args
	return c.RunCmd(tessa)
}

func customPsm(n int) []string {
	return []string{"--psm", fmt.Sprint(n), "-c", "tessedit_create_alto=1", "hoot", "quiet"}
}
