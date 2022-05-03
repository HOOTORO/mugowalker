package navi

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type UiMap map[string]Place

type Place struct {
	Name   string
	Depth  int
	Entry  TPoint
	Parent *Place
}

type TPoint struct {
	X int
	Y int
}

type Navigator struct {
	curentPLace *Place
}

type Walker interface {
	GoForward(x, y int)
	GoBack()
}

type Locator interface {
	IsPlace(*Place)
}

func (p *Place) Nparent(n int) (nparent *Place) {
	nparent = p
	stepsN := p.Depth - n
	if stepsN > 0 {
		for i := 0; i < stepsN; i++ {
			nparent = nparent.Parent
		}
	}

	return
}

func (n *Navigator) Walk(w Walker, target *Place) {
	if n == nil {
		n.curentPLace = target.Nparent(1)
	}
	log.Debugf("Let's take a walk from >>%v<< to >>%v<<", n.curentPLace.Name, target.Name)
	for n.curentPLace.Depth != 1 {
		w.GoForward(target.Entry.X, target.Entry.Y)
		time.Sleep(3 * time.Second)
		n.curentPLace = n.curentPLace.Parent
	}

	for i := 0; i < target.Depth; i++ {
		nextStep := target.Nparent(i + 1)
		log.Debugf("Make a #%d step to --> %v", i+1, nextStep.Name)
		w.GoForward(nextStep.Entry.X, nextStep.Entry.Y)
		time.Sleep(7 * time.Second)
		n.curentPLace = nextStep
	}
}
