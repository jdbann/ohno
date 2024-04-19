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

	cellSize := state.cellSize()

	// Scroll panel
	canvasBounds := rl.NewRectangle(bounds.X, bounds.Y, float32(state.imageSize.X)*cellSize, float32(state.imageSize.Y)*cellSize)
	canvasView := rl.Rectangle{}
	rui.ScrollPanel(bounds, "Canvas", canvasBounds, &state.canvasScroll, &canvasView)
	centerScrollPanelContents(bounds, canvasBounds, &canvasView)
	rl.BeginScissorMode(canvasView.ToInt32().X, canvasView.ToInt32().Y, canvasView.ToInt32().Width, canvasView.ToInt32().Height)
	defer rl.EndScissorMode()

	canvasView.X += state.canvasScroll.X
	canvasView.Y += state.canvasScroll.Y

	// Canvas grid
	var mouseCell rl.Vector2
	gridBounds := rl.NewRectangle(canvasView.X, canvasView.Y, canvasBounds.Width, canvasBounds.Height)
	rui.Grid(gridBounds, "Canvas", cellSize, 1, &mouseCell)

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
			destRec := rl.NewRectangle(float32(x)*cellSize, float32(y)*cellSize, cellSize, cellSize)
			rl.DrawRectanglePro(destRec, origin, 0, state.palette[bg])
			rl.DrawTexturePro(state.tileTexture, sourceRec, destRec, origin, 0, state.palette[fg])
		}
	}

	// Selection
	destRec := rl.NewRectangle(canvasView.X+float32(state.canvasSelection.X)*cellSize, canvasView.Y+float32(state.canvasSelection.Y)*cellSize, cellSize, cellSize)
	rl.DrawRectangleLinesEx(destRec, 1.5, selectionColor)
}

const rayguiWindowboxStatusbarHeight = float32(24)

func centerScrollPanelContents(bounds, content rl.Rectangle, view *rl.Rectangle) {
	borderWidth := float32(rui.GetStyle(rui.DEFAULT, rui.BORDER_WIDTH))

	panelSpace := rl.NewRectangle(
		bounds.X+borderWidth,
		bounds.Y+rayguiWindowboxStatusbarHeight,
		bounds.Width-2*borderWidth,
		bounds.Height-rayguiWindowboxStatusbarHeight-borderWidth,
	)

	if content.Width < bounds.Width-2*borderWidth {
		view.X = bounds.X + (panelSpace.Width-content.Width)/2
	}

	if content.Height < panelSpace.Height {
		view.Y = bounds.Y + rayguiWindowboxStatusbarHeight + (panelSpace.Height-content.Height)/2
	}
}
