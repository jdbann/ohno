package ui

import (
	"image"

	rui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
)

func Canvas(bounds rl.Rectangle, img *textmode.Image, tilesTexture rl.Texture2D, selection *image.Point, scroll *rl.Vector2, zoom float32) bool {
	cellSize := float32(img.Tileset().TileSize()) * zoom
	imgSize := img.TileBounds().Size()

	// Scroll panel
	canvasBounds := rl.NewRectangle(bounds.X, bounds.Y, float32(imgSize.X)*cellSize, float32(imgSize.Y)*cellSize)
	canvasView := rl.Rectangle{}
	rui.ScrollPanel(bounds, "Canvas", canvasBounds, scroll, &canvasView)
	centerScrollPanelContents(bounds, canvasBounds, &canvasView)
	rl.BeginScissorMode(canvasView.ToInt32().X, canvasView.ToInt32().Y, canvasView.ToInt32().Width, canvasView.ToInt32().Height)
	defer rl.EndScissorMode()

	// Canvas grid
	var mouseCell rl.Vector2
	gridBounds := rl.NewRectangle(canvasView.X+scroll.X, canvasView.Y+scroll.Y, canvasBounds.Width, canvasBounds.Height)
	rui.Grid(gridBounds, "Canvas", cellSize, 1, &mouseCell)

	// Change selection if mouse moved inside canvas
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), canvasView) {
		if rl.GetMouseDelta() != rl.NewVector2(0, 0) {
			selection.X = int(mouseCell.X)
			selection.Y = int(mouseCell.Y)
		}
	}

	// Wrap selection within canvas
	x := selection.X
	y := selection.Y

	if x < 0 || x >= imgSize.X {
		x += imgSize.X
		x = x % imgSize.X
	}
	if y < 0 || y >= imgSize.Y {
		y += imgSize.Y
		y = y % imgSize.Y
	}

	*selection = image.Pt(x, y)

	// Draw canvas image
	origin := rl.NewVector2(-canvasView.X-scroll.X, -canvasView.Y-scroll.Y)
	for y := 0; y < imgSize.Y; y++ {
		for x := 0; x < imgSize.X; x++ {
			tile, bg, fg := img.AtCell(textmode.Cell{x, y})
			sourceRec := imageRecToRl(img.Tileset().BoundsAtIndex(tile))
			destRec := rl.NewRectangle(float32(x)*cellSize, float32(y)*cellSize, cellSize, cellSize)
			rl.DrawRectanglePro(destRec, origin, 0, toRGBA(img.Palette()[bg]))
			rl.DrawTexturePro(tilesTexture, sourceRec, destRec, origin, 0, toRGBA(img.Palette()[fg]))
		}
	}

	// Draw selection box
	destRec := rl.NewRectangle(-origin.X+float32(selection.X)*cellSize, -origin.Y+float32(selection.Y)*cellSize, cellSize, cellSize)
	rl.DrawRectangleLinesEx(destRec, 1.5, selectionColor)

	// Return true if clicking in canvas
	return rl.CheckCollisionPointRec(rl.GetMousePosition(), canvasView) && rl.IsMouseButtonPressed(rl.MouseButtonLeft)
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
