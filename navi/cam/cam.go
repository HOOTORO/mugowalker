package cam

import (
	"image"
	"strconv"
	"worker/imaginer"

	log "github.com/sirupsen/logrus"
)

type Cam interface {
	Screencap(string) ([]byte, error)
	PullScreen(string) string
}

type Gallery interface {
	OpenPng(string) image.Image
	SaveImg(string, image.Image) error
}

type View struct {
	name    string
	imgPath string
	areas   map[string]image.Image
}

func Capture(device Cam, locname string) *View {
	r, e := device.Screencap(locname)
	log.Debugf("Captr adb:  --> %v, %v", r, e)
	newFname := device.PullScreen(locname)
	log.Debugf("Captr adb:  --> %v, %v", r, e)
	res := &View{name: locname, imgPath: newFname, areas: make(map[string]image.Image)}

	log.Debugf("Location Res  --> %v", res)
	return res

}

func (l *View) Area(gallery Gallery, x, y, r int) {

	areaord := str(x) + str(y) + str(r)
	img := gallery.OpenPng(l.imgPath)
	box := imaginer.Concat(img, x-r, y-r, x+r, y+r)
	l.areas[areaord] = box

}

func str(x int) string {
	return strconv.Itoa(x)
}
