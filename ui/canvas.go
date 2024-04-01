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

	// Canvas grid
	rui.Grid(gridBounds, "Canvas", float32(state.tileSize), 1, &mouseCell)

	if mouseCell.X >= 0 && mouseCell.Y >= 0 {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			state.image.Set(int(mouseCell.X), int(mouseCell.Y), state.tileSelection, state.bgSelection, state.fgSelection)

			rl.BeginTextureMode(state.canvasRenderTexture)
			sourceRec := imageRecToRl(state.tileset.BoundsAtIndex(state.tileSelection))
			destRec := rl.NewRectangle(mouseCell.X*float32(state.tileSize), mouseCell.Y*float32(state.tileSize), float32(state.tileSize), float32(state.tileSize))
			rl.DrawRectanglePro(destRec, rl.Vector2{}, 0, state.palette[state.bgSelection])
			rl.DrawTexturePro(state.tileTexture, sourceRec, destRec, rl.Vector2{}, 0, state.palette[state.fgSelection])
			rl.EndTextureMode()
		}
	}

	imageRec := rl.NewRectangle(0, 0, float32(state.imageSize.X*state.tileSize), -float32(state.imageSize.Y*state.tileSize))
	rl.DrawTexturePro(state.canvasRenderTexture.Texture, imageRec, gridBounds, rl.Vector2{}, 0, rl.White)
}
