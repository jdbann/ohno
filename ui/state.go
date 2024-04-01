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

	canvasRenderTexture rl.RenderTexture2D

	palette     []color.RGBA
	bgSelection int
	fgSelection int

	tilepickerSpacing float32
}

type StateParams struct {
	TilepickerSpacing float32
}

func NewState(params StateParams) State {
	return State{
		tilepickerSpacing: params.TilepickerSpacing,
	}
}

func (s *State) NewImage(w, h int, palette []color.RGBA) {
	imgPalette := make(color.Palette, len(palette))
	for i, c := range palette {
		imgPalette[i] = c
	}
	img, err := textmode.NewImage(w, h, s.tileset, imgPalette)
	if err != nil {
		panic(img)
	}

	s.canvasRenderTexture = rl.LoadRenderTexture(int32(w*s.tileSize), int32(h*s.tileSize))
	rl.BeginTextureMode(s.canvasRenderTexture)
	rl.DrawRectangle(0, 0, int32(w*s.tileSize), int32(h*s.tileSize), palette[0])
	rl.EndTextureMode()

	s.image = img
	s.imageSize = image.Pt(w, h)
	s.palette = palette
	s.bgSelection = 0
	s.fgSelection = 1
}

func (s *State) LoadTileset(img *rl.Image, size int) {
	src := img.ToImage()
	tileset, err := textmode.NewTileset(src, size)
	if err != nil {
		panic(err)
	}

	tilesetImage := rl.NewImageFromImage(tileset.Image())
	rl.ImageColorReplace(tilesetImage, color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 0})

	s.tileset = tileset
	s.tileSize = size
	s.tileTexture = rl.LoadTextureFromImage(tilesetImage)
}

func (s State) boundsForTilepickerCell(cell textmode.Cell) rl.Rectangle {
	return rl.NewRectangle(float32(cell.X)*s.tilepickerCellSize(), float32(cell.Y)*s.tilepickerCellSize(), s.tilepickerCellSize(), s.tilepickerCellSize())
}

func (s State) tilepickerCellSize() float32 {
	return float32(s.tileSize) + s.tilepickerSpacing
}

func (s State) innerBoundsForCoords(cell textmode.Cell) rl.Rectangle {
	return recTranslate(recAdd(s.boundsForTilepickerCell(cell), rl.NewVector2(-s.tilepickerSpacing, -s.tilepickerSpacing)), rl.NewVector2(s.tilepickerSpacing, s.tilepickerSpacing))
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

func recAdd(r rl.Rectangle, v rl.Vector2) rl.Rectangle {
	return rl.Rectangle{
		X:      r.X,
		Y:      r.Y,
		Width:  r.Width + v.X,
		Height: r.Height + v.Y,
	}
}

func recTranslate(r rl.Rectangle, v rl.Vector2) rl.Rectangle {
	return rl.Rectangle{
		X:      r.X + v.X,
		Y:      r.Y + v.Y,
		Width:  r.Width,
		Height: r.Height,
	}
}