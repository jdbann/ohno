package textmode

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
)

type imageFileFormat struct {
	Size struct {
		Width  int
		Height int
	}
	Tiles    []byte
	BGColors []byte
	FGColors []byte
	Palette  []string
	Tileset  struct {
		Stride   int
		TileSize int
		Image    []uint8
	}
}

func Encode(w io.Writer, m *Image) error {
	tiles, err := encodeIntSlice(m.tiles)
	if err != nil {
		return err
	}

	bgColors, err := encodeIntSlice(m.bgColors)
	if err != nil {
		return err
	}

	fgColors, err := encodeIntSlice(m.fgColors)
	if err != nil {
		return err
	}

	palette := make([]string, len(m.palette))
	for i, c := range m.palette {
		palette[i] = encodeColor(c)
	}

	f := imageFileFormat{
		Size: struct {
			Width  int
			Height int
		}{
			Width:  m.w,
			Height: m.h,
		},
		Tiles:    tiles,
		BGColors: bgColors,
		FGColors: fgColors,
		Palette:  palette,
		Tileset: struct {
			Stride   int
			TileSize int
			Image    []uint8
		}{
			Stride:   m.tileset.img.Stride,
			TileSize: m.tileset.tileSize,
			Image:    m.tileset.img.Pix,
		},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(f)
}

func Decode(r io.Reader) (*Image, error) {
	var f imageFileFormat

	if err := json.NewDecoder(r).Decode(&f); err != nil {
		return nil, err
	}

	var tiles []int
	if err := gob.NewDecoder(bytes.NewBuffer(f.Tiles)).Decode(&tiles); err != nil {
		return nil, err
	}
	tiles, err := decodeIntSlice(f.Tiles)
	if err != nil {
		return nil, err
	}

	bgColors, err := decodeIntSlice(f.BGColors)
	if err != nil {
		return nil, err
	}

	fgColors, err := decodeIntSlice(f.FGColors)
	if err != nil {
		return nil, err
	}

	palette := make(color.Palette, len(f.Palette))
	for i, c := range f.Palette {
		col := color.RGBA{}
		if _, err := fmt.Sscanf(c, colorFmt, &col.R, &col.G, &col.B, &col.A); err != nil {
			return nil, err
		}
		palette[i] = col
	}

	tilesetHeight := len(f.Tileset.Image) / f.Tileset.Stride
	tilesetImg := image.NewPaletted(image.Rect(0, 0, f.Tileset.Stride, tilesetHeight), tilesetPalette())
	tilesetImg.Pix = f.Tileset.Image

	tileset, err := NewTileset(tilesetImg, f.Tileset.TileSize)
	if err != nil {
		return nil, err
	}

	return &Image{
		w:        f.Size.Width,
		h:        f.Size.Height,
		tiles:    tiles,
		bgColors: bgColors,
		fgColors: fgColors,
		palette:  palette,
		tileset:  tileset,
	}, nil
}

const (
	colorFmt = "#%02x%02x%02x%02x"
)

func encodeColor(in color.Color) string {
	r, g, b, a := in.RGBA()
	return fmt.Sprintf(colorFmt, r>>8, g>>8, b>>8, a>>8)
}

func decodeColor(in string) (color.Color, error) {
	var out color.RGBA
	if _, err := fmt.Sscanf(in, colorFmt, &out.R, &out.G, &out.B, &out.A); err != nil {
		return nil, err
	}
	return out, nil
}

func encodeIntSlice(in []int) ([]byte, error) {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(in); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decodeIntSlice(in []byte) ([]int, error) {
	var out []int
	if err := gob.NewDecoder(bytes.NewBuffer(in)).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
