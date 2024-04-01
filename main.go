package main

import (
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
)

func main() {
	rl.InitWindow(windowBounds.ToInt32().Width, windowBounds.ToInt32().Height, "ohno")
	defer rl.CloseWindow()

	tilepickerState := &ui.TilepickerState{
		Spacing: 1,
	}
	tilepickerState.LoadTileset(rl.LoadImage("assets/tileset.png"), 8)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		ui.Tilepicker(tilepickerBounds, tilepickerState)

		rl.EndDrawing()
	}
}
