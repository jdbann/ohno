package ui

import (
	"image"

	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var colorCellSize = 16

func Colorpicker(bounds rl.Rectangle, state *State) {
	// Panel
	rui.Panel(bounds, "Colorpicker")

	var mouseCell rl.Vector2
	colorGridSize := image.Pt(len(state.palette), 1)
	gridBounds := rl.NewRectangle(bounds.X+margins["l"], bounds.Y+margins["t"], float32(colorGridSize.X*colorCellSize), float32(colorGridSize.Y*colorCellSize))

	// Color grid
	rui.Grid(gridBounds, "Color", float32(colorCellSize), 1, &mouseCell)

	// Mouse interaction
	if mouseCell.X >= 0 && mouseCell.Y >= 0 {
		newIdx := int(mouseCell.X)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			state.fgSelection = newIdx
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
			state.bgSelection = newIdx
		}
	}

	// Keyboard interaction
	if rl.IsKeyDown(rl.KeyLeftAlt) || rl.IsKeyDown(rl.KeyRightAlt) {
		x := state.fgSelection % colorGridSize.X
		y := state.fgSelection / colorGridSize.X

		if rl.IsKeyPressed(rl.KeyLeft) {
			x--
		} else if rl.IsKeyPressed(rl.KeyRight) {
			x++
		} else if rl.IsKeyPressed(rl.KeyUp) {
			y--
		} else if rl.IsKeyPressed(rl.KeyDown) {
			y++
		}

		if x < 0 || x >= colorGridSize.X {
			x += colorGridSize.X
			x = x % colorGridSize.X
		}
		if y < 0 || y >= colorGridSize.Y {
			y += colorGridSize.Y
			y = y % colorGridSize.Y
		}

		state.fgSelection = y*colorGridSize.X + x
	}

	if rl.IsKeyPressed(rl.KeyC) {
		state.fgSelection, state.bgSelection = state.bgSelection, state.fgSelection
	}

	// Colors
	for i, c := range state.palette {
		destRec := rl.NewRectangle(gridBounds.X+float32(i*colorCellSize), gridBounds.Y, float32(colorCellSize), float32(colorCellSize))
		rl.DrawRectanglePro(destRec, rl.Vector2{}, 0, c)
	}
}
