package navi

import (
	"image"
	"strconv"
	"worker/imaginer"
	// log "github.com/sirupsen/logrus"
)

//Saves the dat

func (l *Location) CutSector(img image.Image, x, y, r int, name string) image.Image {

	box := imaginer.Concat(l.Etalons[0], x-r, y-r, x+r, y+r)
	l.Sectors[name] = &Sector{x, y, r, box}
	return box
}

// func (l *Location) EtalonSamples() {
// 	if l.Etalons == nil {
// 		l.Etalons = make([]image.Image, 0)
// 		l.Etalons, _ =
// 	}

// }

func str(x int) string {
	return strconv.Itoa(x)
}
