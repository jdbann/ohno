package ui

import (
	"image"

	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
)

var (
	margins = map[string]float32{"t": 36, "r": 12, "b": 12, "l": 12}

	hoverColor     = rl.Red
	selectionColor = rl.Blue
)

func Tilepicker(bounds rl.Rectangle, state *TilepickerState) {
	// Panel
	rui.Panel(bounds, "Tilepicker")

	if state.texture.ID == 0 {
		return
	}

	var mouseCell rl.Vector2

	tileGridSize := state.tileset.GridSize()
	gridBounds := rl.NewRectangle(bounds.X+margins["l"], bounds.Y+margins["t"], float32(tileGridSize.X)*state.cellSize(), float32(tileGridSize.Y)*state.cellSize())
	gridOrigin := rl.NewVector2(gridBounds.X, gridBounds.Y)

	// Tile grid
	rui.Grid(gridBounds, "Tile", state.cellSize(), 1, &mouseCell)

	var tileIdx int
	for y := 0; y < tileGridSize.Y; y++ {
		for x := 0; x < tileGridSize.X; x++ {
			sourceRec := imageRecToRl(state.tileset.BoundsAtIndex(tileIdx))
			destRec := recTranslate(state.innerBoundsForCoords(textmode.Cell{x, y}), gridOrigin)
			rl.DrawTexturePro(state.texture, sourceRec, destRec, rl.Vector2{}, 0, rl.White)

			tileIdx++
		}
	}

	// Mouse interaction
	if mouseCell.X >= 0 && mouseCell.Y >= 0 {
		mouseRec := recTranslate(state.boundsForCell(textmode.Cell{int(mouseCell.X), int(mouseCell.Y)}), gridOrigin)
		rl.DrawRectangleLinesEx(mouseRec, 1.5, hoverColor)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			state.Selection = state.tileset.IndexForCell(textmode.Cell{int(mouseCell.X), int(mouseCell.Y)})
		}
	}

	// Selected tile
	selectionRec := recTranslate(state.selectionBounds(), gridOrigin)
	rl.DrawRectangleLinesEx(selectionRec, 1.5, selectionColor)
}

type TilepickerState struct {
	Selection int
	Spacing   float32

	tileset  *textmode.Tileset
	tileSize int
	texture  rl.Texture2D
}

func (s *TilepickerState) LoadTileset(img *rl.Image, size int) {
	src := img.ToImage()
	tileset, err := textmode.NewTileset(src, size)
	if err != nil {
		panic(err)
	}

	s.tileset = tileset
	s.tileSize = tileset.TileSize()
	s.texture = rl.LoadTextureFromImage(img)
}

func (s TilepickerState) boundsForCell(cell textmode.Cell) rl.Rectangle {
	return rl.NewRectangle(float32(cell.X)*s.cellSize(), float32(cell.Y)*s.cellSize(), s.cellSize(), s.cellSize())
}

func (s TilepickerState) cellSize() float32 {
	return float32(s.tileSize) + s.Spacing
}

func (s TilepickerState) innerBoundsForCoords(cell textmode.Cell) rl.Rectangle {
	return recTranslate(recAdd(s.boundsForCell(cell), rl.NewVector2(-s.Spacing, -s.Spacing)), rl.NewVector2(s.Spacing, s.Spacing))
}

func (s *TilepickerState) selectionBounds() rl.Rectangle {
	cell := s.tileset.CellForIndex(s.Selection)
	return s.boundsForCell(cell)
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
