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
		x, y    int
		wantImg image.Image
		wantErr error
	}

	run := func(t *testing.T, tc testCase) {
		actual, err := tc.tileset.At(tc.x, tc.y)
		assert.ErrorIs(t, err, tc.wantErr)
		assert.DeepEqual(t, actual, tc.wantImg, cmpImagePixels)
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
			x:       1,
			y:       1,
			wantImg: tile11,
			wantErr: nil,
		},
		{
			name:    "ErrInvalidCoords",
			tileset: tileset,
			x:       2,
			y:       1,
			wantImg: nil,
			wantErr: textmode.ErrInvalidCoords,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestTileset_AtIndex(t *testing.T) {
	type testCase struct {
		name    string
		tileset *textmode.Tileset
		index   int
		wantImg image.Image
		wantErr error
	}

	run := func(t *testing.T, tc testCase) {
		actual, err := tc.tileset.AtIndex(tc.index)
		assert.ErrorIs(t, err, tc.wantErr)
		assert.DeepEqual(t, actual, tc.wantImg, cmpImagePixels)
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
			index:   3,
			wantImg: tile11,
			wantErr: nil,
		},
		{
			name:    "ErrInvalidCoords",
			tileset: tileset,
			index:   4,
			wantImg: nil,
			wantErr: textmode.ErrInvalidCoords,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
