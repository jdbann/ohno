package ui

import (
	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
)

var (
	margins = map[string]float32{"t": 36, "r": 12, "b": 12, "l": 12}

	hoverColor     = rl.Red
	selectionColor = rl.Blue
)

func Tilepicker(bounds rl.Rectangle, state *State) {
	// Panel
	rui.Panel(bounds, "Tilepicker")

	if state.tileTexture.ID == 0 {
		return
	}

	var mouseCell rl.Vector2

	tileGridSize := state.tileset.GridSize()
	gridBounds := rl.NewRectangle(bounds.X+margins["l"], bounds.Y+margins["t"], float32(tileGridSize.X)*state.tilepickerCellSize(), float32(tileGridSize.Y)*state.tilepickerCellSize())
	gridOrigin := rl.NewVector2(gridBounds.X, gridBounds.Y)

	// Tile grid
	rui.Grid(gridBounds, "Tile", state.tilepickerCellSize(), 1, &mouseCell)

	// Mouse interaction
	if mouseCell.X >= 0 && mouseCell.Y >= 0 {
		mouseRec := recTranslate(state.boundsForTilepickerCell(textmode.Cell{int(mouseCell.X), int(mouseCell.Y)}), gridOrigin)
		rl.DrawRectangleLinesEx(mouseRec, 1.5, hoverColor)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			state.tileSelection = state.tileset.IndexForCell(textmode.Cell{int(mouseCell.X), int(mouseCell.Y)})
		}
	}

	// Keyboard interaction
	if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
		x := state.tileSelection % tileGridSize.X
		y := state.tileSelection / tileGridSize.X

		if rl.IsKeyPressed(rl.KeyLeft) {
			x--
		} else if rl.IsKeyPressed(rl.KeyRight) {
			x++
		} else if rl.IsKeyPressed(rl.KeyUp) {
			y--
		} else if rl.IsKeyPressed(rl.KeyDown) {
			y++
		}

		if x < 0 || x >= tileGridSize.X {
			x += tileGridSize.X
			x = x % tileGridSize.X
		}
		if y < 0 || y >= tileGridSize.Y {
			y += tileGridSize.Y
			y = y % tileGridSize.Y
		}
		state.tileSelection = y*tileGridSize.X + x
	}

	// Tiles
	var tileIdx int
	for y := 0; y < tileGridSize.Y; y++ {
		for x := 0; x < tileGridSize.X; x++ {
			sourceRec := imageRecToRl(state.tileset.BoundsAtIndex(tileIdx))
			destRec := recTranslate(state.innerBoundsForCoords(textmode.Cell{x, y}), gridOrigin)
			rl.DrawRectanglePro(destRec, rl.Vector2{}, 0, state.palette[state.bgSelection])
			rl.DrawTexturePro(state.tileTexture, sourceRec, destRec, rl.Vector2{}, 0, state.palette[state.fgSelection])

			tileIdx++
		}
	}

	// Selected tile
	selectionRec := recTranslate(state.selectionBounds(), gridOrigin)
	rl.DrawRectangleLinesEx(selectionRec, 1.5, selectionColor)
}
