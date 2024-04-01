package textmode_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/jdbann/ohno/textmode"
	"gotest.tools/v3/assert"
)

func TestNewTileset(t *testing.T) {
	type testCase struct {
		name    string
		img     image.Image
		size    int
		wantErr error
	}

	run := func(t *testing.T, tc testCase) {
		_, err := textmode.NewTileset(tc.img, tc.size)
		assert.ErrorIs(t, err, tc.wantErr)
	}

	validTiles := newImage(4, 4, color.Palette{color.Black, color.White}, [][]int{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{1, 0, 0, 1},
	})
	colorfulTiles := newImage(4, 4, color.Palette{color.Black, color.White, color.RGBA{255, 0, 0, 0}}, [][]int{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 2, 2, 0},
		{2, 0, 0, 2},
	})

	testCases := []testCase{
		{
			name:    "ok",
			img:     validTiles,
			size:    2,
			wantErr: nil,
		},
		{
			name:    "ErrIncompatibleTileSize",
			img:     validTiles,
			size:    3,
			wantErr: textmode.ErrIncompatibleTileSize,
		},
		{
			name:    "ErrTooManyColors",
			img:     colorfulTiles,
			size:    2,
			wantErr: textmode.ErrTooManyColors,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestTileset_At(t *testing.T) {
	type testCase struct {
		name    string
		tileset *textmode.Tileset
		cell    textmode.Cell
		index   int
		wantImg image.Image
	}

	run := func(t *testing.T, tc testCase) {
		cellActual := tc.tileset.AtCell(tc.cell)
		assert.DeepEqual(t, cellActual, tc.wantImg, cmpImagePixels)

		indexActual := tc.tileset.AtIndex(tc.index)
		assert.DeepEqual(t, indexActual, tc.wantImg, cmpImagePixels)
	}

	validTiles := newImage(4, 4, color.Palette{color.Black, color.White}, [][]int{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{1, 0, 1, 0},
	})
	tile11 := newImage(2, 2, color.Palette{color.Black, color.White}, [][]int{
		{1, 0},
		{1, 0},
	})
	tileset, err := textmode.NewTileset(validTiles, 2)
	assert.NilError(t, err)

	testCases := []testCase{
		{
			name:    "ok",
			tileset: tileset,
			cell:    textmode.Cell{1, 1},
			index:   3,
			wantImg: tile11,
		},
		{
			name:    "ErrInvalidCoords",
			tileset: tileset,
			cell:    textmode.Cell{2, 1},
			index:   4,
			wantImg: image.Rectangle{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestTileset_BoundsAt(t *testing.T) {
	type testCase struct {
		name       string
		tileset    *textmode.Tileset
		cell       textmode.Cell
		index      int
		wantBounds image.Rectangle
	}

	run := func(t *testing.T, tc testCase) {
		cellActual := tc.tileset.BoundsAtCell(tc.cell)
		assert.DeepEqual(t, cellActual, tc.wantBounds)

		indexActual := tc.tileset.BoundsAtIndex(tc.index)
		assert.DeepEqual(t, indexActual, tc.wantBounds)
	}

	validTiles := newImage(4, 4, color.Palette{color.Black, color.White}, [][]int{
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
			tileset:    tileset,
			cell:       textmode.Cell{1, 1},
			index:      3,
			wantBounds: image.Rect(2, 2, 4, 4),
		},
		{
			name:       "out of bounds",
			tileset:    tileset,
			cell:       textmode.Cell{2, 1},
			index:      4,
			wantBounds: image.Rectangle{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestTileset_GridSize(t *testing.T) {
	validTiles := newImage(4, 4, color.Palette{color.Black, color.White}, [][]int{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{1, 0, 1, 0},
	})
	tileset, err := textmode.NewTileset(validTiles, 2)
	assert.NilError(t, err)

	assert.DeepEqual(t, tileset.GridSize(), image.Point{2, 2})
}

func TestTileset_TileSize(t *testing.T) {
	validTiles := newImage(4, 4, color.Palette{color.Black, color.White}, [][]int{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{1, 0, 1, 0},
	})
	tileset, err := textmode.NewTileset(validTiles, 2)
	assert.NilError(t, err)

	assert.Equal(t, tileset.TileSize(), 2)
}
