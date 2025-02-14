package dat

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/quasilyte/fantasy-general-tools/serdat"
)

// https://wiki.amigaos.net/wiki/ILBM_IFF_Interleaved_Bitmap

type PaletteFile struct {
	TransparentIndex uint8

	Colors [256]color.NRGBA

	ReverseIndex map[color.NRGBA]uint8
}

func ParsePaletteFile(data []byte) (*PaletteFile, error) {
	f := &PaletteFile{}

	bmhdIndex := bytes.Index(data, []byte("BMHD"))
	if bmhdIndex == -1 {
		return nil, fmt.Errorf("BMHD chunk not found")
	}

	data = data[bmhdIndex+len("BMHD"):]
	data = data[2:] // skip width UWORD
	data = data[2:] // skip height UWORD
	data = data[2:] // skip x WORD
	data = data[2:] // skip y WORD

	nPlanes, data := scanUint8(data)
	if nPlanes != 0 {
		return nil, fmt.Errorf("BMHD: expected zero nPlanes, found %02x", nPlanes)
	}

	data = data[1:] // skip masking UBYTE
	data = data[1:] // skip compression UBYTE
	data = data[1:] // skip pad1 UBYTE

	transparentIndex, data := scanUint16(data)
	f.TransparentIndex = uint8(transparentIndex)

	cmapIndex := bytes.Index(data, []byte("CMAP"))
	if cmapIndex == -1 {
		return nil, fmt.Errorf("CMAP chunk not found")
	}

	data = data[cmapIndex+len("CMAP"):]

	paletteSize, data := scanUint32BE(data)
	if paletteSize%3 != 0 {
		return nil, fmt.Errorf("invalid palette size: %d", paletteSize)
	}
	colorIndex := 0
	for i := 0; i < int(paletteSize); i += 3 {
		f.Colors[colorIndex] = color.NRGBA{
			R: data[i+0],
			G: data[i+1],
			B: data[i+2],
			A: 255,
		}
		colorIndex++
	}

	f.ReverseIndex = make(map[color.NRGBA]uint8, len(f.Colors))
	for i, clr := range &f.Colors {
		f.ReverseIndex[clr] = uint8(i)
	}

	return f, nil
}

func (p *PaletteFile) Get(index int) color.NRGBA {
	return p.Colors[uint8(index)]
}

func (p *PaletteFile) GetIndex(clr color.NRGBA) uint8 {
	i, ok := p.ReverseIndex[clr]
	if !ok {
		return p.TransparentIndex
	}
	return i
}

func (p *PaletteFile) ToSerdat() serdat.Palette {
	serialized := serdat.Palette{
		Colors: make(map[string]string, len(p.Colors)),
	}
	for i, clr := range p.Colors {
		colorKey := fmt.Sprintf("%02x", i)
		colorHex := fmt.Sprintf("%02x%02x%02x", clr.R, clr.G, clr.B)
		serialized.Colors[colorKey] = colorHex
	}
	return serialized
}
