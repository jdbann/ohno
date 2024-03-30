package textmode_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/jdbann/ohno/textmode"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestNewImage(t *testing.T) {
	type testCase struct {
		name    string
		w, h    int
		tileset *textmode.Tileset
		palette color.Palette
		wantErr error
	}

	run := func(t *testing.T, tc testCase) {
		_, err := textmode.NewImage(tc.w, tc.h, tc.tileset, tc.palette)
		assert.ErrorIs(t, err, tc.wantErr)
	}

	validTiles := newImage(4, 4, paletteBW, [][]int{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{1, 0, 1, 0},
	})
	tileset, err := textmode.NewTileset(validTiles, 2)
	assert.NilError(t, err)

	testCases := []testCase{
		{
			name:    "ok",
			w:       4,
			h:       4,
			tileset: tileset,
			palette: paletteBW,
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestImage_At(t *testing.T) {
	type testCase struct {
		name      string
		x, y      int
		wantColor color.Color
	}

	validTiles := newImage(2, 1, paletteBW, [][]int{
		{0, 1},
	})
	tileset, err := textmode.NewTileset(validTiles, 1)
	assert.NilError(t, err)
	img, err := textmode.NewImage(3, 1, tileset, paletteRGB)
	assert.NilError(t, err)
	assert.NilError(t, img.Set(1, 0, 1, 1, 2)) // fg
	assert.NilError(t, img.Set(2, 0, 0, 1, 2)) // bg

	run := func(t *testing.T, tc testCase) {
		assert.Equal(t, img.At(tc.x, tc.y), tc.wantColor)
	}

	testCases := []testCase{
		{
			name:      "zero values",
			x:         0,
			y:         0,
			wantColor: paletteRGB[0],
		},
		{
			name:      "set fg color",
			x:         1,
			y:         0,
			wantColor: paletteRGB[2],
		},
		{
			name:      "set bg color",
			x:         2,
			y:         0,
			wantColor: paletteRGB[1],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestImage_Bounds(t *testing.T) {
	type testCase struct {
		name       string
		w, h       int
		tileset    *textmode.Tileset
		wantBounds image.Rectangle
	}

	run := func(t *testing.T, tc testCase) {
		img, err := textmode.NewImage(tc.w, tc.h, tc.tileset, nil)
		assert.NilError(t, err)
		assert.DeepEqual(t, img.Bounds(), tc.wantBounds)
	}

	validTiles := newImage(4, 4, paletteBW, [][]int{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{1, 0, 1, 0},
	})
	tileset, err := textmode.NewTileset(validTiles, 2)
	assert.NilError(t, err)

	testCases := []testCase{
		{
			name:       "ok",
			w:          4,
			h:          4,
			tileset:    tileset,
			wantBounds: image.Rect(0, 0, 8, 8),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestImage_Set(t *testing.T) {
	t.Run("draw whole tileset", func(t *testing.T) {
		validTiles := newImage(4, 4, paletteBW, [][]int{
			{0, 0, 1, 1},
			{0, 0, 1, 1},
			{0, 1, 1, 0},
			{1, 0, 1, 0},
		})
		tileset, err := textmode.NewTileset(validTiles, 2)
		assert.NilError(t, err)

		image, err := textmode.NewImage(4, 4, tileset, paletteRGB)
		assert.NilError(t, err)

		assert.NilError(t, image.Set(1, 1, 0, 1, 2))
		assert.NilError(t, image.Set(2, 1, 1, 1, 2))
		assert.NilError(t, image.Set(1, 2, 2, 1, 2))
		assert.NilError(t, image.Set(2, 2, 3, 1, 2))

		assertEncodedImage(t, image)
	})
}

func assertEncodedImage(t *testing.T, image image.Image) {
	t.Helper()
	enc := &png.Encoder{}
	var b bytes.Buffer
	assert.NilError(t, enc.Encode(&b, image))
	golden.AssertBytes(t, b.Bytes(), t.Name()+".png")
}
