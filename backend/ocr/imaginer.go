package ocr

import (
	"fmt"
	"math/rand"
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

// func (mo MagicOptions) Arguments() (args []string) {
// 	if mo.ColorSpace != "" {
// 		args = append(args, "-colorspace", mo.ColorSpace)
// 	}
// 	if mo.Alpha != "" {
// 		args = append(args, "-alpha", mo.Alpha)
// 	}
// 	if mo.Threshold > 0 {
// 		args = append(args, "-threshold", c.F("%v%", mo.Threshold))
// 	}
// 	if mo.Edge > 0 {
// 		args = append(args, "-edge", c.F("%v", mo.Edge))
// 	}
// 	if mo.Negate {
// 		args = append(args, "-negate")
// 	}
// 	if mo.BlackThreshold > 0 {
// 		args = append(args, "-black-threshold", c.F("%v%", mo.BlackThreshold))
// 	}
// 	if mo.BrightnessContrast != "" {
// 		args = append(args, "-brightness-contrast", mo.BrightnessContrast)
// 	}
// 	if mo.Density > 0 {
// 		args = append(args, c.F("%v", mo.Density))
// 	}

// 	return
// }

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
	m.fout = "wd/temp/" + magicRandPrefix() + ".png"
}
func (m *Magick) Prepared() string {
	return m.fout
}
func init() {
	// rand.Seed(time.Now().Unix())
	// blaine.args = append(blaine.args, c.ActiveUser().AltImagick...)
	// log = c.Logger()

	// origin
	// imagickSets = append(imagickSets, c.ActiveUser().AltImagick)
	// imagickSets = append(imagickSets, c.ActiveUser().ImagickCfg())

	// great for recognize buttons
	// buttonsOpts := MagicOptions{
	// 	ColorSpace:         "",
	// 	Alpha:              "off",
	// 	Threshold:          70,
	// 	Edge:               0,
	// 	Negate:             true,
	// 	BlackThreshold:     0,
	// 	BrightnessContrast: "-70x1",
	// 	Density:            0,
	// }

	// genOpts := MagicOptions{
	// 	ColorSpace:         "Gray",
	// 	Alpha:              "off",
	// 	Threshold:          70,
	// 	Edge:               1,
	// 	Negate:             true,
	// 	BlackThreshold:     70,
	// 	BrightnessContrast: "",
	// 	Density:            0,
	// }

	// altGenOpts := MagicOptions{
	// 	ColorSpace:         "",
	// 	Alpha:              "off",
	// 	Threshold:          0,
	// 	Edge:               0,
	// 	Negate:             true,
	// 	BlackThreshold:     0,
	// 	BrightnessContrast: "-25x10",
	// 	Density:            300,
	// }
	// imagickSets = append(imagickSets, buttonsOpts.Arguments(), genOpts.Arguments(), altGenOpts.Arguments())

	// imagickSets = append(imagickSets, []string
	// btn recog
	//  {"-alpha", "off", "-brightness-contrast", "-70x1", "-negate", "-threshold", "70%"})
	// general recognition
	// {"-colorspace", "Gray", "-alpha", "off", "-threshold", "70%", "-edge", "1", "-negate", "-black-threshold", "70%"}
	// {"-alpha", "off", "-brightness-contrast", "-25x10", "-density", "300", "-negate", "-threshold", "70%"}

}

func magicRandPrefix() string {
	return fmt.Sprintf("maaagick_%d", rand.Intn(999))
}
