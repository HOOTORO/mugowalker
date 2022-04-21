package img

import (
	"image"
	"image/draw"
)

type ImgH interface {
	Compare(image.Image, image.Image) bool
	Text(image.Image) (string, error)
	Concat(image.Image, int, int)
	resize()
}

func Concat(img1 image.Image, x1, y1, x2, y2 int) image.Image {
	sr := image.Rect(x1, y1, x2, y2)
	rect := image.Rectangle{image.Point{}, image.Point{}.Add(sr.Size())}
	dst := image.NewRGBA(rect)
	draw.Draw(dst, rect, img1, sr.Min, draw.Src)
	return dst
}
