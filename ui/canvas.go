package ui

import (
	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
)

func Canvas(bounds rl.Rectangle, state *State) {
	rui.DummyRec(bounds, "Canvas")

	if state.image == nil {
		return
	}

	var mouseCell rl.Vector2

	gridBounds := rl.NewRectangle(bounds.X, bounds.Y, float32(state.imageSize.X*state.tileSize), float32(state.imageSize.Y*state.tileSize))

	// Canvas grid
	rui.Grid(gridBounds, "Canvas", float32(state.tileSize), 1, &mouseCell)

	if mouseCell.X >= 0 && mouseCell.Y >= 0 {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			state.image.Set(int(mouseCell.X), int(mouseCell.Y), state.tileSelection, state.bgSelection, state.fgSelection)
		}
	}

	// Canvas image
	for y := 0; y < state.imageSize.Y; y++ {
		for x := 0; x < state.imageSize.X; x++ {
			tile, bg, fg := state.image.AtCell(textmode.Cell{x, y})
			sourceRec := imageRecToRl(state.tileset.BoundsAtIndex(tile))
			destRec := rl.NewRectangle(float32(x*state.tileSize), float32(y*state.tileSize), float32(state.tileSize), float32(state.tileSize))
			rl.DrawRectanglePro(destRec, rl.Vector2{}, 0, state.palette[bg])
			rl.DrawTexturePro(state.tileTexture, sourceRec, destRec, rl.Vector2{}, 0, state.palette[fg])
		}
	}
}
