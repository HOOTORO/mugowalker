package navi

import (
	"image"
	"strconv"
	"worker/imaginer"
	// log "github.com/sirupsen/logrus"
)

func (n *Location) Capture(dev Device) image.Image {
	dev.Screencap(n.Name)
	fpath := dev.PullScreen(n.Name)
	img := imaginer.OpenImg(fpath)

	// log.Debugf("Current Position Captured  --> %v", res)
	return img

}

type Device interface {
	Screencap(string) ([]byte, error)
	PullScreen(string) string
}

func (l *Location) CutSector(img image.Image, x, y, r int, name string) image.Image {

	box := imaginer.Concat(l.Etalons[0], x-r, y-r, x+r, y+r)
	l.Sectors[name] = &Sector{x, y, r, box}
	return box
}

func (l *Location) EtalonSamples(d DSaver) {
	if l.Etalons == nil {
		l.Etalons = make([]image.Image, 0)
		l.Etalons, _ = d.LocEtalons(l.Name)
	}

}

func str(x int) string {
	return strconv.Itoa(x)
}
