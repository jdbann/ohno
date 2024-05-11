package app

import (
	"errors"
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

	canvasBounds = rl.NewRectangle(
		0, 0,
		windowBounds.Width-sidebarWidth, windowBounds.Height,
	)

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

		ui.Tilepicker(tilepickerBounds, &state)
		ui.Canvas(canvasBounds, &state)
		ui.Colorpicker(colorpickerBounds, &state)

		if rl.IsKeyDown(rl.KeyLeftSuper) || rl.IsKeyDown(rl.KeyRightSuper) {
			switch {
			case rl.IsKeyPressed(rl.KeyS):
				saveFile(scope, defaultFilename, state.Image())
			}
		}

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
