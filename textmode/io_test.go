package textmode_test

import (
	"bytes"
	"image/color"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jdbann/ohno/textmode"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestEncodeDecode(t *testing.T) {
	// Build tileset
	tilesetImg := newImage(4, 4, paletteBW, [][]int{
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{1, 0, 1, 0},
	})
	tileset, err := textmode.NewTileset(tilesetImg, 2)
	assert.NilError(t, err)

	// Build palette
	palette := color.Palette{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 127, 255},
	}

	// Build image
	originalImg, err := textmode.NewImage(8, 6, tileset, palette)
	assert.NilError(t, err)

	// Draw in the image
	assert.NilError(t, originalImg.Set(0, 0, 1, 1, 2))
	assert.NilError(t, originalImg.Set(1, 1, 0, 1, 2))

	// Encode
	var b bytes.Buffer
	assert.NilError(t, textmode.Encode(&b, originalImg))

	// Check file format
	golden.Assert(t, b.String(), filepath.Join(t.Name(), "simple-image.json"))

	// Decode
	decodedImg, err := textmode.Decode(&b)
	assert.NilError(t, err)

	// Check images are equivalent
	assert.DeepEqual(t, originalImg, decodedImg, cmp.AllowUnexported(textmode.Image{}, textmode.Tileset{}))
}
