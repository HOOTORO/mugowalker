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
