package textmode

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"slices"
)

var (
	ErrIncompatibleTileSize = errors.New("incompatible tile size")
	ErrInvalidCoords        = errors.New("invalid tile coordinates")
	ErrTooManyColors        = errors.New("too many colors")
)

type Tileset struct {
	img  *image.Paletted
	size int
}

func NewTileset(src image.Image, size int) (*Tileset, error) {
	srcSize := src.Bounds()
	if srcSize.Dx()%size != 0 || srcSize.Dy()%size != 0 {
		return nil, ErrIncompatibleTileSize
	}

	var colors []color.Color
	img := image.NewPaletted(src.Bounds(), color.Palette{color.Black, color.White})

	for y := srcSize.Min.Y; y <= srcSize.Max.Y; y++ {
		for x := srcSize.Min.X; x <= srcSize.Max.X; x++ {
			c := src.At(x, y)
			if !slices.Contains(colors, c) {
				colors = append(colors, c)
			}
			if len(colors) > 2 {
				return nil, ErrTooManyColors
			}

			img.Set(x, y, src.At(x, y))
		}
	}

	return &Tileset{
		img:  img,
		size: size,
	}, nil
}

func (t Tileset) At(x, y int) (image.Image, error) {
	r := image.Rect(0, 0, t.size, t.size).Add(image.Point{x * t.size, y * t.size})
	if !r.In(t.img.Bounds()) {
		return nil, fmt.Errorf("%w: (%d, %d)", ErrInvalidCoords, x, y)
	}

	return t.img.SubImage(r), nil
}

func (t Tileset) AtIndex(idx int) (image.Image, error) {
	y := idx / (t.img.Bounds().Dx() / t.size)
	x := idx % (t.img.Bounds().Dx() / t.size)
	return t.At(x, y)
}
