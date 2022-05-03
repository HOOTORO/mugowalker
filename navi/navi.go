package navi

import (
	"image"
	"time"
	"worker/adb"
	"worker/imaginer"

	log "github.com/sirupsen/logrus"
)

var trylim int = 5

type UiMap map[string]Location

type Location struct {
	Name   string
	Depth  int
	Entry  TPoint
	etalon image.Image
	areas  map[string]image.Image
	Parent *Location
	Peers  []*Location
}

type TPoint struct {
	X int
	Y int
}

type Navigator struct {
	*adb.Device
	Liveloc *Location
}

// Travels the world
type Walker interface {
	GoForward(x, y int)
	GoBack()
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

func (n *Navigator) Step(w Walker, target *Location) {
	log.Debugf("Little step for bot (>>%v>> to <<%v>>", n.Liveloc.Name, target.Name)
	w.GoForward(target.Entry.X, target.Entry.Y)
	capture := n.Capture(target)
	if !n.ExpectedLocation(capture, target) && trylim > 0 {
		time.Sleep(3 * time.Second)
		n.Step(w, target)
	} else {
		trylim = 5
		panic("WE FAILED MASTQA")
	}

}

func (n *Navigator) ExpectedLocation(v *View, loc *Location) bool {
	return imaginer.Similarity(v.img, loc.etalon)

}
