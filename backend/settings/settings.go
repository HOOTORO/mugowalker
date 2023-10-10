package settings

import (
	"fmt"
)

type Settings struct {
	DrawStep   bool      `json:"drawstep"`
	Logfile    string    `json:"logfile"`
	Loglevel   string    `json:"loglevel"`
	Imagick    Imagick   `json:"imagick"`
	Tesseract  Tesseract `json:"tesseract"`
	Bluestacks struct {
		Instance string `json:"instance"`
		Cmd      string `json:"cmd"`
		Package  string `json:"package"`
	} `json:"bluestacks"`
	IgnoredWords []string `json:"ignoredwords"`
}

type Imagick struct {
	Colorspace    string `json:"colorspace"`
	Alpha         string `json:"alpha"`
	Threshold     string `json:"threshold"`
	Edge          string `json:"edge"`
	Negate        bool   `json:"negate"`
	GaussianBlur  string `json:"gaussian-blur"`
	AutoThreshold string `json:"auto-threshold"`
}

type Tesseract struct {
	Psm  int      `json:"psm"`
	Args []string `json:"args"`
}

func Default() *Settings {
	return &Settings{
		DrawStep: false,
		Logfile:  "app.log",
		Loglevel: "FATAL",
		Imagick: Imagick{
			Colorspace:   "Gray",
			Alpha:        "off",
			Threshold:    "90%",
			Edge:         "2",
			Negate:       false,
			GaussianBlur: "0%",
		},
		Tesseract: Tesseract{
			Psm:  6,
			Args: []string{},
		},
		Bluestacks: struct {
			Instance string "json:\"instance\""
			Cmd      string "json:\"cmd\""
			Package  string "json:\"package\""
		}{
			Instance: "rcv64",
			Cmd:      "launch",
			Package:  "com.myapp",
		},
		IgnoredWords: []string{},
	}
}
func (s *Settings) String() string {
	return fmt.Sprintf("\nlogfile\t\tloglevel\n%v\t%v\nIgnored Words:\n\t%v\n%v\n%v", s.Logfile, s.Loglevel, s.IgnoredWords, s.Tesseract, s.Imagick)
}
func (s *Imagick) String() string {
	return fmt.Sprintf("<cfg IMAGICK>: \n%v %v %v %v %v %v", s.Alpha, s.Colorspace, s.Threshold, s.GaussianBlur, s.Edge, s.Negate)
}
func (s *Tesseract) String() string {
	return fmt.Sprintf("<cfg TESSA> \tpsm: %v\nargs:%v\n", s.Psm, s.Args)
}

func (i *Imagick) Args() []string {
	str := make([]string, 0)
	if i.Colorspace != "" {
		str = append(str, "-colorspace", i.Colorspace)
	}
	if i.Alpha != "" {
		str = append(str, "-alpha", i.Alpha)
	}
	if i.Threshold != "" {
		str = append(str, "-threshold", i.Threshold)
	}
	if i.Edge != "" {
		str = append(str, "-edge", i.Edge)
	}
	if i.Negate {
		str = append(str, "-negate")
	}
	if i.GaussianBlur != "" {
		str = append(str, "-gaussian-blur", i.GaussianBlur)
	}
	if i.AutoThreshold != "" {
		str = append(str, "-auto-threshold", i.AutoThreshold)
	}
	return str
}
