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
	Name    string
	Depth   int
	Entry   TPoint
	etalons []image.Image
	areas   map[string]image.Image
	Parent  *Location
	Peers   []*Location
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

//Saves the data
type DSaver interface {
	SaveLoc(string, image.Image) error
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
	//give oit time to load
	time.Sleep(3 * time.Second)
	screen := n.Capture(target)
	if !n.ExpectedLocation(screen, target) && trylim > 0 {
		n.Step(w, target)
		trylim--
	} else {
		trylim = 5
		panic("WE FAILED MASTQA")
	}

}

func (n *Navigator) ExpectedLocation(v *View, loc *Location) bool {
	return imaginer.Similarity(v.img, loc.etalons[0])

}

func (n *Navigator) EtalonStep(w Walker, dm DSaver, t *Location) {
	log.Debugf("ETALON STEP for bot (>>%v>> to <<%v>>", n.Liveloc.Name, t.Name)
	captscr := n.Capture(n.Liveloc)
	n.Liveloc.etalons = append(n.Liveloc.etalons, captscr.img)
	dm.SaveLoc(captscr.name, captscr.img)
	w.GoForward(t.Entry.X, t.Entry.Y)
	time.Sleep(10 * time.Second)
	n.Liveloc = t
}

// func
