package ui

import (
	"image"

	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
)

func Canvas(bounds rl.Rectangle, state *State) {
	if state.image == nil {
		return
	}

	// Scroll panel
	canvasBounds := rl.NewRectangle(bounds.X, bounds.Y, float32(state.imageSize.X*state.tileSize), float32(state.imageSize.Y*state.tileSize))
	canvasView := rl.Rectangle{}
	rui.ScrollPanel(bounds, "Canvas", canvasBounds, &state.canvasScroll, &canvasView)
	rl.BeginScissorMode(canvasView.ToInt32().X, canvasView.ToInt32().Y, canvasView.ToInt32().Width, canvasView.ToInt32().Height)
	defer rl.EndScissorMode()

	// Canvas grid
	var mouseCell rl.Vector2
	gridBounds := rl.NewRectangle(canvasView.X, canvasView.Y, canvasBounds.Width, canvasBounds.Height)
	rui.Grid(gridBounds, "Canvas", float32(state.tileSize), 1, &mouseCell)

	// Mouse interaction
	if mouseCell.X >= 0 && mouseCell.Y >= 0 {
		if rl.GetMouseDelta() != rl.NewVector2(0, 0) {
			state.canvasSelection = image.Pt(int(mouseCell.X), int(mouseCell.Y))
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			state.image.Set(int(mouseCell.X), int(mouseCell.Y), state.tileSelection, state.bgSelection, state.fgSelection)
		}
	}

	// Keyboard interaction
	if !(rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) || rl.IsKeyDown(rl.KeyLeftAlt) || rl.IsKeyDown(rl.KeyRightAlt)) {
		x := state.canvasSelection.X
		y := state.canvasSelection.Y

		if rl.IsKeyPressed(rl.KeyLeft) {
			x--
		} else if rl.IsKeyPressed(rl.KeyRight) {
			x++
		} else if rl.IsKeyPressed(rl.KeyUp) {
			y--
		} else if rl.IsKeyPressed(rl.KeyDown) {
			y++
		}

		if x < 0 || x >= state.imageSize.X {
			x += state.imageSize.X
			x = x % state.imageSize.X
		}
		if y < 0 || y >= state.imageSize.Y {
			y += state.imageSize.Y
			y = y % state.imageSize.Y
		}

		state.canvasSelection = image.Pt(x, y)
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		state.image.Set(state.canvasSelection.X, state.canvasSelection.Y, state.tileSelection, state.bgSelection, state.fgSelection)
	}

	// Canvas image
	origin := rl.NewVector2(-canvasView.X, -canvasView.Y)
	for y := 0; y < state.imageSize.Y; y++ {
		for x := 0; x < state.imageSize.X; x++ {
			tile, bg, fg := state.image.AtCell(textmode.Cell{x, y})
			sourceRec := imageRecToRl(state.tileset.BoundsAtIndex(tile))
			destRec := rl.NewRectangle(float32(x*state.tileSize), float32(y*state.tileSize), float32(state.tileSize), float32(state.tileSize))
			rl.DrawRectanglePro(destRec, origin, 0, state.palette[bg])
			rl.DrawTexturePro(state.tileTexture, sourceRec, destRec, origin, 0, state.palette[fg])
		}
	}

	// Selection
	destRec := rl.NewRectangle(canvasView.X+float32(state.canvasSelection.X*state.tileSize), canvasView.Y+float32(state.canvasSelection.Y*state.tileSize), float32(state.tileSize), float32(state.tileSize))
	rl.DrawRectangleLinesEx(destRec, 1.5, selectionColor)
}
