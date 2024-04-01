package ui

import (
	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Canvas(bounds rl.Rectangle, state *State) {
	rui.DummyRec(bounds, "Canvas")

	if state.image == nil {
		return
	}

	var mouseCell rl.Vector2

	gridBounds := rl.NewRectangle(bounds.X, bounds.Y, float32(state.imageSize.X*state.tileSize), float32(state.imageSize.Y*state.tileSize))

	// Tile grid
	rui.Grid(gridBounds, "Canvas", float32(state.tileSize), 1, &mouseCell)
}
