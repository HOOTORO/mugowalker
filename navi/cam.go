package navi

import (
	"image"
	"strconv"
	"worker/imaginer"

	log "github.com/sirupsen/logrus"
)

type View struct {
	name  string
	img   image.Image
	areas map[string]image.Image
}

func (n *Navigator) Capture(l *Location) *View {
	r, e := n.Screencap(l.Name)
	log.Debugf("Captr adb:  --> %v, %v", r, e)
	newFname := n.PullScreen(l.Name)
	img := imaginer.OpenImg(newFname)
	log.Debugf("Captr adb:  --> %v, %v", r, e)
	res := &View{name: l.Name, img: img, areas: make(map[string]image.Image)}

	log.Debugf("Location Res  --> %v", res)
	return res

}

func (l *View) Area(x, y, r int) {
	areaord := str(x) + str(y) + str(r)
	box := imaginer.Concat(l.img, x-r, y-r, x+r, y+r)
	l.areas[areaord] = box

}

func str(x int) string {
	return strconv.Itoa(x)
}
