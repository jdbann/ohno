package app

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/ui"
)

var (
	windowBounds         = rl.NewRectangle(0, 0, 1280, 720)
	sidebarWidth float32 = 320

	colorpickerHeight float32 = 180
	colorpickerBounds         = rl.NewRectangle(
		windowBounds.Width-sidebarWidth, windowBounds.Height-colorpickerHeight,
		sidebarWidth, colorpickerHeight,
	)

	tilepickerBounds = rl.NewRectangle(
		windowBounds.Width-sidebarWidth, 0,
		sidebarWidth, windowBounds.Height-colorpickerHeight,
	)

	canvasBounds = rl.NewRectangle(
		0, 0,
		windowBounds.Width-sidebarWidth, windowBounds.Height,
	)
)

func Run() {
	rl.InitWindow(windowBounds.ToInt32().Width, windowBounds.ToInt32().Height, "ohno")
	defer rl.CloseWindow()

	state := ui.State{}
	state.LoadTileset(rl.LoadImage("assets/tileset.png"), 8)
	state.NewImage(10, 10, []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	})

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		ui.Tilepicker(tilepickerBounds, &state)
		ui.Canvas(canvasBounds, &state)
		ui.Colorpicker(colorpickerBounds, &state)

		rl.EndDrawing()
	}
}
