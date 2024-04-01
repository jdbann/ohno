package ui

import (
	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var colorCellSize = 16

func Colorpicker(bounds rl.Rectangle, state *State) {
	// Panel
	rui.Panel(bounds, "Colorpicker")

	var mouseCell rl.Vector2
	gridBounds := rl.NewRectangle(bounds.X+margins["l"], bounds.Y+margins["t"], float32(len(state.palette)*colorCellSize), float32(colorCellSize))

	// Color grid
	rui.Grid(gridBounds, "Color", float32(colorCellSize), 1, &mouseCell)

	if mouseCell.X >= 0 && mouseCell.Y >= 0 {
		newIdx := int(mouseCell.X)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			state.fgSelection = newIdx
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
			state.bgSelection = newIdx
		}
	}

	// Colors
	for i, c := range state.palette {
		destRec := rl.NewRectangle(gridBounds.X+float32(i*colorCellSize), gridBounds.Y, float32(colorCellSize), float32(colorCellSize))
		rl.DrawRectanglePro(destRec, rl.Vector2{}, 0, c)
	}
}
