package ocr

import (
	"math/rand"
	"time"

	c "worker/cfg"
)

const mag = "magick"

var imagickSets [][]string

// Magick params to run ImageMagick via cfg.Runnable
type Magick struct {
	path string
	f    string
	args []string
	fout string
}

var blaine = &Magick{
	path: mag,
	args: make([]string, 0),
}

// Args implementing Runnable
func (m *Magick) Args() []string {
	newargs := make([]string, 0)
	newargs = append(newargs, m.f)
	newargs = append(newargs, m.args...)
	newargs = append(newargs, m.fout)
	return newargs
}

// Path implementing Runnable
func (m *Magick) Path() string {
	m.path = mag
	return m.path
}
func (m *Magick) SetFile(f string) {
	m.f = f
	m.fout = c.TempFile(magicRandPrefix() + ".png")
}
func (m *Magick) Prepared() string {
	return m.fout
}
func init() {
	rand.Seed(time.Now().Unix())
	blaine.args = append(blaine.args, c.ActiveUser().AltImagick...)
	log = c.Logger()

	// origin
	// imagickSets = append(imagickSets, c.ActiveUser().AltImagick)
	// imagickSets = append(imagickSets, c.ActiveUser().ImagickCfg())

	// great for recognize buttons
	imagickSets = append(imagickSets, []string{"-alpha", "off", "-brightness-contrast", "-70x1", "-negate", "-threshold", "70%"})
	// general recognition
	imagickSets = append(imagickSets, []string{"-colorspace", "Gray", "-alpha", "off", "-threshold", "70%", "-edge", "1", "-negate", "-black-threshold", "70%"}) //80
	imagickSets = append(imagickSets, []string{"-alpha", "off", "-brightness-contrast", "-25x10", "-density", "300", "-negate", "-threshold", "70%"})

	// imagickSets = append(imagickSets, []string{"-alpha", "off", "-fill", "black", "-fuzz", "30%", "+opaque", "#FFFFFF"})
	// imagickSets = append(imagickSets, []string{"-alpha", "off", "-negate", "-threshold", "100", "-negate"})
}

func (m *Magick) NewRandArgs() {
	m.args = imagickSets[rand.Intn(len(imagickSets))]
	log.Debug(c.RFW("NEW IMAGICK ARGS #> ", m.args))
}

func magicRandPrefix() string {
	return c.F("maaagick_%d", rand.Intn(999))
}
