package textmode

import (
	"errors"
	"image"
	"image/color"
	"slices"
)

var (
	ErrIncompatibleTileSize = errors.New("incompatible tile size")
	ErrInvalidCoords        = errors.New("invalid tile coordinates")
	ErrTooManyColors        = errors.New("too many colors")
)

type Cell = image.Point

type Tileset struct {
	gridSize image.Point
	img      *image.Paletted
	tileSize int
}

func NewTileset(src image.Image, tileSize int) (*Tileset, error) {
	srcSize := src.Bounds()
	gridSize := srcSize.Size().Div(tileSize)
	if gridSize.Mul(tileSize) != srcSize.Size() {
		return nil, ErrIncompatibleTileSize
	}

	var colors []color.Color
	img := image.NewPaletted(src.Bounds(), color.Palette{color.Black, color.White})

	for y := srcSize.Min.Y; y < srcSize.Max.Y; y++ {
		for x := srcSize.Min.X; x < srcSize.Max.X; x++ {
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
		gridSize: gridSize,
		img:      img,
		tileSize: tileSize,
	}, nil
}

func (t Tileset) AtCell(cell Cell) image.Image {
	if !t.cellInGrid(cell) {
		return image.Rectangle{}
	}

	r := t.BoundsAtCell(cell)
	return t.img.SubImage(r)
}

func (t Tileset) AtIndex(idx int) image.Image {
	cell := t.CellForIndex(idx)
	return t.AtCell(cell)
}

func (t Tileset) BoundsAtCell(cell Cell) image.Rectangle {
	if !t.cellInGrid(cell) {
		return image.Rectangle{}
	}

	return image.Rect(0, 0, t.tileSize, t.tileSize).Add(cell.Mul(t.tileSize))
}

func (t Tileset) BoundsAtIndex(idx int) image.Rectangle {
	cell := t.CellForIndex(idx)
	return t.BoundsAtCell(cell)
}

func (t Tileset) CellForIndex(idx int) Cell {
	y := idx / (t.img.Bounds().Dx() / t.tileSize)
	x := idx % (t.img.Bounds().Dx() / t.tileSize)
	return Cell{x, y}
}

func (t Tileset) GridSize() image.Point {
	return t.gridSize
}

func (t Tileset) Image() image.Image {
	return t.img
}

func (t Tileset) TileSize() int {
	return t.tileSize
}

func (t Tileset) IndexForCell(cell Cell) int {
	return cell.Y*t.img.Bounds().Dx()/t.tileSize + cell.X
}

func (t Tileset) cellInGrid(cell Cell) bool {
	return cell.In(image.Rectangle{Max: t.gridSize})
}
