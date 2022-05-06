package navi

import (
	"image"
	"worker/adb"
	"worker/imaginer"
	// log "github.com/sirupsen/logrus"
)

var trylim int = 5

type UiMap map[string]Location

type Location struct {
	Name    string
	Depth   int
	Entry   TPoint
	Etalons []image.Image
	Parent  *Location
	Sectors map[string]*Sector
}

type Sector struct {
	X, Y, R int
	area    image.Image
}

type TPoint struct {
	X int
	Y int
}

type Navigator struct {
	*adb.Device
	Liveloc *Location
}

func (p *Location) Nparent(n int) (nparent *Location) {
	nparent = p
	stepsN := p.Depth - n
	if stepsN > 0 {
		for i := 0; i < stepsN; i++ {
			nparent = nparent.Parent
		}
	}

	return
}

func (n *Location) IsLocation(img image.Image) bool {
	return imaginer.Similarity(img, n.Etalons[0])

}

// func
