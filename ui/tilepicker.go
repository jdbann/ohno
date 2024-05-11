package ui

import (
	"image"
	"image/color"

	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
)

var (
	hoverColor     = rl.Red
	selectionColor = rl.Blue
)

func Tilepicker(bounds rl.Rectangle, img *textmode.Image, tileTexture rl.Texture2D, selection *image.Point, scroll *rl.Vector2, zoom float32, bgColor, fgColor color.Color) {
	cellSize := float32(img.Tileset().TileSize()) * zoom
	tileGridSize := img.Tileset().GridSize()

	// Scroll panel
	tilepickerBounds := rl.NewRectangle(bounds.X, bounds.Y, float32(tileGridSize.X)*cellSize, float32(tileGridSize.Y)*cellSize)
	tilepickerView := rl.Rectangle{}
	rui.ScrollPanel(bounds, "Tilepicker", tilepickerBounds, scroll, &tilepickerView)
	rl.BeginScissorMode(tilepickerView.ToInt32().X, tilepickerView.ToInt32().Y, tilepickerView.ToInt32().Width, tilepickerView.ToInt32().Height)
	defer rl.EndScissorMode()

	// Tilepicker grid
	var mouseCell rl.Vector2
	gridBounds := rl.NewRectangle(tilepickerView.X+scroll.X, tilepickerView.Y+scroll.Y, tilepickerBounds.Width, tilepickerBounds.Height)
	gridOrigin := rl.NewVector2(gridBounds.X, gridBounds.Y)
	rui.Grid(gridBounds, "Tile", cellSize, 1, &mouseCell)

	// Change selection if mouse pressed inside canvas
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), tilepickerView) {
		cellBounds := rl.NewRectangle(float32(mouseCell.X)*cellSize, float32(mouseCell.Y)*cellSize, cellSize, cellSize)
		mouseRec := recTranslate(cellBounds, gridOrigin)
		rl.DrawRectangleLinesEx(mouseRec, 1.5, hoverColor)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			selection.X = int(mouseCell.X)
			selection.Y = int(mouseCell.Y)
		}
	}

	// Draw tiles
	rl.DrawRectanglePro(gridBounds, rl.Vector2{}, 0, toRGBA(bgColor))
	sourceRec := rl.NewRectangle(0, 0, float32(tileTexture.Width), float32(tileTexture.Height))
	rl.DrawTexturePro(tileTexture, sourceRec, gridBounds, rl.Vector2{}, 0, toRGBA(fgColor))

	// Draw selection box
	selectionBounds := rl.NewRectangle(float32(selection.X)*cellSize, float32(selection.Y)*cellSize, cellSize, cellSize)
	selectionRec := recTranslate(selectionBounds, gridOrigin)
	rl.DrawRectangleLinesEx(selectionRec, 1.5, selectionColor)
}
