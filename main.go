package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/ui"
)

var (
	windowBounds = rl.NewRectangle(0, 0, 1280, 720)

	tilepickerWidth  float32 = 320
	tilepickerBounds         = rl.NewRectangle(
		windowBounds.Width-tilepickerWidth, 0,
		tilepickerWidth, windowBounds.Height,
	)

	canvasBounds = rl.NewRectangle(
		0, 0,
		windowBounds.Width-tilepickerWidth, windowBounds.Height,
	)
)

func main() {
	rl.InitWindow(windowBounds.ToInt32().Width, windowBounds.ToInt32().Height, "ohno")
	defer rl.CloseWindow()

	stateParams := ui.StateParams{
		TilepickerSpacing: 1,
	}
	state := ui.NewState(stateParams)
	state.LoadTileset(rl.LoadImage("assets/tileset.png"), 8)
	state.NewImage(10, 10, color.Palette{})

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		ui.Tilepicker(tilepickerBounds, &state)
		ui.Canvas(canvasBounds, &state)

		rl.EndDrawing()
	}
}
