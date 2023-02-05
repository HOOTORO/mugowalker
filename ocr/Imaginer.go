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
}

var blaine = Magick{
	path: mag,
	args: make([]string, 0),
}

// Args implementing Runnable
func (m Magick) Args() []string {
	out := c.TempFile(m.f)
	m.args = append([]string{m.f}, user.AltImagick...)
	m.args = append(m.args, out)
	return m.args
}

// Path implementing Runnable
func (m Magick) Path() string {
	m.path = mag
	return m.path
}

func (m *Magick) out() string {
	return c.TempFile(m.f)
}

func init() {
	blaine.args = append(blaine.args, c.ActiveUser().AltImagick...)
	log = c.Logger()

	// origin
	imagickSets = append(imagickSets, []string{"-colorspace", "Gray", "-alpha", "off", "threshold", "75%", "-edge", "2", "-negate", "-black-threshold", "90%"})
	imagickSets = append(imagickSets, []string{"-alpha", "off", "-fill", "black", "-fuzz", "30%", "+opaque", "#FFFFFF"})
	imagickSets = append(imagickSets, []string{"-alpha", "off", "-brightness-contrast", "-40x10", "-units pixelsperinch", "-density", "300", "-negate", "-noise", "10", "-threshold", "70%"})
	imagickSets = append(imagickSets, []string{"-alpha", "off", "-negate", "-threshold", "100", "-negate"})
}

func MagickArgs() []string {
	rand.Seed(time.Now().Unix())
	return imagickSets[rand.Intn(len(imagickSets))]
}
