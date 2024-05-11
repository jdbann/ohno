package app

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/ohno/textmode"
	"github.com/jdbann/ohno/ui"
	gap "github.com/muesli/go-app-paths"
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
	tilepickerSelection = image.Point{}
	tilepickerScroll    = rl.Vector2{}

	canvasBounds = rl.NewRectangle(
		0, 0,
		windowBounds.Width-sidebarWidth, windowBounds.Height,
	)
	canvasSelection = image.Point{}
	canvasScroll    = rl.Vector2{}

	zoom float32 = 8

	defaultFilename = "canvas.json"
)

func Run() {
	scope := gap.NewScope(gap.User, "ohno")
	img := loadFile(scope, defaultFilename)
	if img == nil {
		img = defaultImg()
	}

	rl.InitWindow(windowBounds.ToInt32().Width, windowBounds.ToInt32().Height, "ohno")
	defer rl.CloseWindow()

	state := ui.State{}
	state.SetImage(img)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		bgColor, fgColor := img.Palette()[state.BGSelection], img.Palette()[state.FGSelection]
		ui.Tilepicker(tilepickerBounds, img, state.TileTexture, &tilepickerSelection, &tilepickerScroll, zoom, bgColor, fgColor)
		if ui.Canvas(canvasBounds, img, state.TileTexture, &canvasSelection, &canvasScroll, zoom) {
			tileIdx := img.Tileset().IndexForCell(textmode.Cell(tilepickerSelection))
			img.Set(canvasSelection.X, canvasSelection.Y, tileIdx, state.BGSelection, state.FGSelection)
		}
		ui.Colorpicker(colorpickerBounds, &state)

		if rl.IsKeyDown(rl.KeyLeftSuper) || rl.IsKeyDown(rl.KeyRightSuper) {
			switch {
			case rl.IsKeyPressed(rl.KeyS):
				saveFile(scope, defaultFilename, img)
			}
		} else {
			switch {
			case rl.IsKeyPressed(rl.KeySpace):
				tileIdx := img.Tileset().IndexForCell(textmode.Cell(tilepickerSelection))
				img.Set(canvasSelection.X, canvasSelection.Y, tileIdx, state.BGSelection, state.FGSelection)

			case rl.IsKeyPressed(rl.KeyLeft):
				canvasSelection.X--
			case rl.IsKeyPressed(rl.KeyRight):
				canvasSelection.X++
			case rl.IsKeyPressed(rl.KeyUp):
				canvasSelection.Y--
			case rl.IsKeyPressed(rl.KeyDown):
				canvasSelection.Y++

			case rl.IsKeyPressed(rl.KeyA):
				tilepickerSelection.X--
			case rl.IsKeyPressed(rl.KeyD):
				tilepickerSelection.X++
			case rl.IsKeyPressed(rl.KeyW):
				tilepickerSelection.Y--
			case rl.IsKeyPressed(rl.KeyS):
				tilepickerSelection.Y++
			}
		}

		canvasSelection = wrapSelection(img.TileBounds().Size(), canvasSelection)
		tilepickerSelection = wrapSelection(img.Tileset().GridSize(), tilepickerSelection)

		rl.EndDrawing()
	}
}

func defaultImg() *textmode.Image {
	tilesetFile, err := os.Open("assets/tileset.png")
	if err != nil {
		log.Fatal(err)
	}
	defer tilesetFile.Close()

	tilesetImg, err := png.Decode(tilesetFile)
	if err != nil {
		log.Fatal(err)
	}

	tileset, err := textmode.NewTileset(tilesetImg, 8)
	if err != nil {
		log.Fatal(err)
	}

	palette := color.Palette{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
	}

	img, err := textmode.NewImage(10, 10, tileset, palette)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func loadFile(scope *gap.Scope, name string) *textmode.Image {
	path, err := scope.DataPath(name)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		log.Fatal(err)
	}

	img, err := textmode.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func saveFile(scope *gap.Scope, name string, img *textmode.Image) {
	path, err := scope.DataPath(name)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o770); err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := textmode.Encode(file, img); err != nil {
		log.Fatal(err)
	}
}

func wrapSelection(bounds, selection image.Point) image.Point {
	return image.Point{
		X: (selection.X + bounds.X) % bounds.X,
		Y: (selection.Y + bounds.Y) % bounds.Y,
	}
}
