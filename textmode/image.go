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

	tile := i.tileset.AtIndex(i.tiles[idx])

	px := tile.(*image.Paletted).ColorIndexAt((x%i.tileset.tileSize)+tile.Bounds().Min.X, (y%i.tileset.tileSize)+tile.Bounds().Min.Y)

	if px == 0 {
		return i.palette[i.bgColors[idx]]
	}
	return i.palette[i.fgColors[idx]]
}

func (i Image) AtCell(cell Cell) (int, int, int) {
	idx := i.tileToIndex(cell)
	return i.tiles[idx], i.bgColors[idx], i.fgColors[idx]
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.w*i.tileset.tileSize, i.h*i.tileset.tileSize)
}

func (i Image) ColorModel() color.Model {
	return i.palette
}

func (i Image) Palette() color.Palette {
	return i.palette
}

func (i *Image) Set(x, y, tile, bg, fg int) error {
	idx := i.tileToIndex(image.Pt(x, y))
	i.tiles[idx] = tile
	i.bgColors[idx] = bg
	i.fgColors[idx] = fg
	return nil
}

func (i Image) TileBounds() image.Rectangle {
	return image.Rect(0, 0, i.w, i.h)
}

func (i Image) Tileset() *Tileset {
	return i.tileset
}

func (i Image) pxToIndex(pt image.Point) int {
	return pt.Y/i.tileset.tileSize*i.w + pt.X/i.tileset.tileSize
}

func (i Image) tileToIndex(pt image.Point) int {
	return pt.Y*i.w + pt.X
}
