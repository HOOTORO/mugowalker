package navi

type Locator interface {
	Path(Place) []TPoint
}

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
