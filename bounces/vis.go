package bounces

import (
	img "image"
	"image/color"
	"image/png"
	"io"
	"math"
)

func (r Results) Draw(w io.Writer, f func(a float64) float64) error {
	res := int(math.Sqrt(float64(len(r.Image))))
	img := img.NewGray(img.Rect(0, 0, res, res))

	max := int32(0)
	for _, v := range r.Image {
		if v > max {
			max = v
		}
	}
	for yi := 0; yi < res; yi++ {
		for xi := 0; xi < res; xi++ {
			i := xi + yi*res
			val := f(float64(r.Image[i])) / f(float64(max)) * float64(math.MaxUint8-1)
			img.SetGray(xi, yi, color.Gray{math.MaxUint8 - uint8(val)})
		}
	}
	return png.Encode(w, img)
}
