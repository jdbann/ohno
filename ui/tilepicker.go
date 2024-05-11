package ui

import (
	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
)

var (
	hoverColor     = rl.Red
	selectionColor = rl.Blue
)

func Tilepicker(bounds rl.Rectangle, state *State) {
	if state.TileTexture.ID == 0 {
		return
	}

	cellSize := state.cellSize()
	tileGridSize := state.tileset.GridSize()

	// Scroll panel
	tilepickerBounds := rl.NewRectangle(bounds.X, bounds.Y, float32(tileGridSize.X)*cellSize, float32(tileGridSize.Y)*cellSize)
	tilepickerView := rl.Rectangle{}
	rui.ScrollPanel(bounds, "Tilepicker", tilepickerBounds, &state.tilepickerScroll, &tilepickerView)
	rl.BeginScissorMode(tilepickerView.ToInt32().X, tilepickerView.ToInt32().Y, tilepickerView.ToInt32().Width, tilepickerView.ToInt32().Height)
	defer rl.EndScissorMode()

	var mouseCell rl.Vector2

	gridBounds := rl.NewRectangle(tilepickerView.X+state.tilepickerScroll.X, tilepickerView.Y+state.tilepickerScroll.Y, tilepickerBounds.Width, tilepickerBounds.Height)
	gridOrigin := rl.NewVector2(gridBounds.X, gridBounds.Y)

	// Tile grid
	rui.Grid(gridBounds, "Tile", cellSize, 1, &mouseCell)

	// Mouse interaction
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), tilepickerView) {
		mouseRec := recTranslate(state.boundsForTilepickerCell(textmode.Cell{int(mouseCell.X), int(mouseCell.Y)}), gridOrigin)
		rl.DrawRectangleLinesEx(mouseRec, 1.5, hoverColor)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			state.TileSelection = state.tileset.IndexForCell(textmode.Cell{int(mouseCell.X), int(mouseCell.Y)})
		}
	}

	// Keyboard interaction
	if !(rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)) {
		x := state.TileSelection % tileGridSize.X
		y := state.TileSelection / tileGridSize.X

		switch {
		case rl.IsKeyPressed(rl.KeyA):
			x--
		case rl.IsKeyPressed(rl.KeyD):
			x++
		case rl.IsKeyPressed(rl.KeyW):
			y--
		case rl.IsKeyPressed(rl.KeyS):
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
		state.TileSelection = y*tileGridSize.X + x
	}

	// Tiles
	rl.DrawRectanglePro(gridBounds, rl.Vector2{}, 0, state.palette[state.BGSelection])
	sourceRec := rl.NewRectangle(0, 0, float32(state.TileTexture.Width), float32(state.TileTexture.Height))
	rl.DrawTexturePro(state.TileTexture, sourceRec, gridBounds, rl.Vector2{}, 0, state.palette[state.FGSelection])

	// Selected tile
	selectionRec := recTranslate(state.selectionBounds(), gridOrigin)
	rl.DrawRectangleLinesEx(selectionRec, 1.5, selectionColor)
}
