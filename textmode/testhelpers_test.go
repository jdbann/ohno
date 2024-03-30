package textmode_test

import (
	"image"
	"image/color"

	"github.com/google/go-cmp/cmp"
)

func newImage(w, h int, palette color.Palette, pixels [][]int) image.Image {
	img := image.NewPaletted(image.Rect(0, 0, w, h), palette)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, palette[pixels[y][x]])
		}
	}

	return img
}

var cmpImagePixels = cmp.Transformer("image.Image", func(img image.Image) []color.Color {
	s := img.Bounds()
	pixels := make([]color.Color, 0, s.Dx()*s.Dy())
	for y := s.Min.Y; y < s.Max.Y; y++ {
		for x := s.Min.X; x < s.Max.X; x++ {
			pixels = append(pixels, img.At(x, y))
		}
	}
	return pixels
})

var (
	paletteBW  = color.Palette{color.Black, color.White}
	paletteRGB = color.Palette{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}}
)
