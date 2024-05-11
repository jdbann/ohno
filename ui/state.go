package ui

import (
	"image"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
)

type State struct {
	image     *textmode.Image
	imageSize image.Point

	tileset       *textmode.Tileset
	tileSize      int
	tileSelection int
	tileTexture   rl.Texture2D

	canvasScroll    rl.Vector2
	canvasSelection image.Point

	tilepickerScroll rl.Vector2

	palette     []color.RGBA
	bgSelection int
	fgSelection int
}

func (s *State) Image() *textmode.Image {
	return s.image
}

func (s *State) SetImage(img *textmode.Image) {
	tilesetImage := rl.NewImageFromImage(img.Tileset().Image())
	rl.ImageColorReplace(tilesetImage, color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 0})

	s.tileset = img.Tileset()
	s.tileSize = img.Tileset().TileSize()
	s.tileTexture = rl.LoadTextureFromImage(tilesetImage)

	s.image = img
	s.imageSize = img.TileBounds().Size()
	s.palette = mapSlice(img.Palette(), toRGBA)
	s.bgSelection = 0
	s.fgSelection = 1
}

func (s State) boundsForTilepickerCell(cell textmode.Cell) rl.Rectangle {
	return rl.NewRectangle(float32(cell.X)*s.cellSize(), float32(cell.Y)*s.cellSize(), s.cellSize(), s.cellSize())
}

func (s State) cellSize() float32 {
	zoom := float32(8)
	return float32(s.tileSize) * zoom
}

func (s State) selectionBounds() rl.Rectangle {
	cell := s.tileset.CellForIndex(s.tileSelection)
	return s.boundsForTilepickerCell(cell)
}

func imageRecToRl(r image.Rectangle) rl.Rectangle {
	return rl.NewRectangle(
		float32(r.Min.X), float32(r.Min.Y),
		float32(r.Dx()), float32(r.Dy()),
	)
}

func recTranslate(r rl.Rectangle, v rl.Vector2) rl.Rectangle {
	return rl.Rectangle{
		X:      r.X + v.X,
		Y:      r.Y + v.Y,
		Width:  r.Width,
		Height: r.Height,
	}
}

func mapSlice[In any, Out any](in []In, fn func(In) Out) []Out {
	out := make([]Out, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}
	return out
}

func toRGBA(in color.Color) color.RGBA {
	r, g, b, a := in.RGBA()
	return color.RGBA{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		uint8(a >> 8),
	}
}
