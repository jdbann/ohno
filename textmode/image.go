package textmode

import (
	"image"
	"image/color"
)

type Image struct {
	w, h     int
	tiles    []int
	bgColors []int
	fgColors []int

	palette color.Palette
	tileset *Tileset
}

func NewImage(w, h int, tileset *Tileset, palette color.Palette) (*Image, error) {
	return &Image{
		w:        w,
		h:        h,
		tiles:    make([]int, w*h),
		bgColors: make([]int, w*h),
		fgColors: make([]int, w*h),
		palette:  palette,
		tileset:  tileset,
	}, nil
}

func (i Image) At(x, y int) color.Color {
	idx := i.pxToIndex(image.Pt(x, y))

	tile, err := i.tileset.AtIndex(i.tiles[idx])
	if err != nil {
		panic(err)
	}

	px := tile.(*image.Paletted).ColorIndexAt((x%i.tileset.size)+tile.Bounds().Min.X, (y%i.tileset.size)+tile.Bounds().Min.Y)

	if px == 0 {
		return i.palette[i.bgColors[idx]]
	}
	return i.palette[i.fgColors[idx]]
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.w*i.tileset.size, i.h*i.tileset.size)
}

func (i Image) ColorModel() color.Model {
	return i.palette
}

func (i *Image) Set(x, y, tile, bg, fg int) error {
	idx := i.tileToIndex(image.Pt(x, y))
	i.tiles[idx] = tile
	i.bgColors[idx] = bg
	i.fgColors[idx] = fg
	return nil
}

func (i Image) pxToIndex(pt image.Point) int {
	return pt.Y/i.tileset.size*i.w + pt.X/i.tileset.size
}

func (i Image) tileToIndex(pt image.Point) int {
	return pt.Y*i.w + pt.X
}
