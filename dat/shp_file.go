package dat

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
)

type ShpFile struct {
	UniSize    bool
	UniPalette bool
	Images     []ShpFileImage
}

type ShpFileImage struct {
	DataOffset       int
	PaletteOffset    int
	NumPaletteColors int

	Width  int
	Height int

	StartPixelX int
	StartPixelY int
	EndPixelX   int
	EndPixelY   int

	Pixels [][]color.NRGBA
}

func ShpEncode(f *ShpFile, pal *PaletteFile) ([]byte, error) {
	sizeApprox := 512
	for _, img := range f.Images {
		sizeApprox += 2 * (img.Width * img.Height)
	}
	data := make([]byte, 0, sizeApprox)

	data = append(data, []byte("1.10")...)

	data = binary.LittleEndian.AppendUint32(data, uint32(len(f.Images)))

	// Save this position to handle it later.
	metadataOffset := len(data)
	for range f.Images {
		// Append imgdata and palette offsets.
		// Zero for now, just to allocate space for them (filled later).
		data = binary.LittleEndian.AppendUint32(data, 0)
		data = binary.LittleEndian.AppendUint32(data, 0)
	}

	var encoder shpEncoder
	encoder.pal = pal
	encoder.buf = make([]color.NRGBA, 0, 8)

	for i := range f.Images {
		img := &f.Images[i]

		binary.LittleEndian.PutUint32(data[metadataOffset+(i*8):], uint32(len(data)))

		data = binary.LittleEndian.AppendUint16(data, uint16(img.Height-1))
		data = binary.LittleEndian.AppendUint16(data, uint16(img.Width-1))

		data = append(data, 0, 0, 0, 0) // An unknown chunk

		data = binary.LittleEndian.AppendUint32(data, uint32(img.StartPixelX))
		data = binary.LittleEndian.AppendUint32(data, uint32(img.StartPixelY))
		data = binary.LittleEndian.AppendUint32(data, uint32(img.EndPixelX))
		data = binary.LittleEndian.AppendUint32(data, uint32(img.EndPixelY))

		isEmpty := img.StartPixelX > img.Width && img.StartPixelY > img.Height
		if !isEmpty {
			data = encoder.AppendImgData(data, img)
		}
	}

	// Now that all imgdata is written, we can insert a dummy palette section.
	paletteStart := len(data)
	data = binary.LittleEndian.AppendUint32(data, 0) // 0 as "num colors"

	for i := range f.Images {
		binary.LittleEndian.PutUint32(data[metadataOffset+(i*8)+4:], uint32(paletteStart))
	}

	return data, nil
}

type shpEncoder struct {
	x   int
	y   int
	pal *PaletteFile
	buf []color.NRGBA
}

func (e *shpEncoder) AppendImgData(data []byte, img *ShpFileImage) []byte {
	e.start(img)
	for e.y <= img.EndPixelY {
		data = e.drawRow(data, img)
		data = append(data, 0) // Next row op
		e.x = img.StartPixelX
		e.y++
	}
	return data
}

func (e *shpEncoder) start(img *ShpFileImage) {
	e.x = img.StartPixelX
	e.y = img.StartPixelY
}

func (e *shpEncoder) countRepeats(img *ShpFileImage) int {
	lineLength := 0
	clr := img.Pixels[e.y][e.x]
	for e.x+lineLength+1 <= img.EndPixelX && img.Pixels[e.y][e.x+lineLength+1] == clr {
		lineLength++
	}
	if lineLength == 0 {
		return 0
	}
	return lineLength + 1
}

func (e *shpEncoder) drawLine(data []byte, length int, clr color.NRGBA) []byte {
	paletteIndex := e.pal.GetIndex(clr)

	const maxLineLength = (math.MaxUint8 - 1) / 2
	for length >= 2 {
		segmentLength := length
		if segmentLength > maxLineLength {
			segmentLength = maxLineLength
		}
		data = append(data, byte(segmentLength*2), paletteIndex)
		e.x += segmentLength
		length -= segmentLength
	}

	return data
}

func (e *shpEncoder) flushBuf(data []byte) []byte {
	if len(e.buf) == 0 {
		return data
	}
	buf := e.buf

	const maxSize = (math.MaxUint8 - 1) / 2
	for len(buf) > 0 {
		n := len(buf)
		if n > maxSize {
			n = maxSize
		}
		data = append(data, byte(n*2)+1)
		for i := 0; i < n; i++ {
			paletteIndex := e.pal.GetIndex(buf[i])
			data = append(data, paletteIndex)
		}
		buf = buf[n:]
	}

	e.buf = e.buf[:0]

	return data
}

func (e *shpEncoder) drawRow(data []byte, img *ShpFileImage) []byte {
	for {
		// Skip all transparent pixels using [0 num] op.
		numTransparent := 0
		for e.x <= img.EndPixelX && img.Pixels[e.y][e.x].A == 0 {
			numTransparent++
			e.x++
		}
		// The row could be finished after advancing x.
		if e.x > img.EndPixelX {
			// If it's the end of the row, no need
			// to skip these pixels -- a next row op will do it for us.
			return data
		}
		for numTransparent > 0 {
			skip := numTransparent
			if skip > math.MaxUint8 {
				skip = math.MaxUint8
			}
			data = append(data, 1, uint8(skip))
			numTransparent -= skip
		}

		e.buf = e.buf[:0]
		for e.x <= img.EndPixelX && img.Pixels[e.y][e.x].A != 0 {
			clr := img.Pixels[e.y][e.x]
			lineLength := e.countRepeats(img)
			minLength := 3
			if len(e.buf) == 0 {
				minLength = 2
			}
			if lineLength >= minLength {
				data = e.flushBuf(data)                  // Flush pending buf pixels
				data = e.drawLine(data, lineLength, clr) // Advances x
				continue
			}
			e.buf = append(e.buf, clr)
			e.x++
		}
		data = e.flushBuf(data) // If there are any pending pixels, flush them
		// The row could be finished after advancing x.
		if e.x > img.EndPixelX {
			return data
		}
	}
}

func ShpToPNG(f *ShpFile) []image.Image {
	if f.UniSize {
		width := f.Images[0].Width
		height := f.Images[0].Height
		pngImage := image.NewNRGBA(image.Rectangle{
			image.Pt(0, 0),
			image.Pt(width*len(f.Images), height),
		})
		offsetX := 0
		for i := range f.Images {
			img := &f.Images[i]
			for y := 0; y < img.Height; y++ {
				for x := 0; x < img.Width; x++ {
					pngImage.SetNRGBA(x+offsetX, y, img.Pixels[y][x])
				}
			}
			offsetX += width
		}
		return []image.Image{pngImage}
	}

	list := make([]image.Image, len(f.Images))
	for i := range f.Images {
		img := &f.Images[i]
		pngImage := image.NewNRGBA(image.Rectangle{
			image.Pt(0, 0),
			image.Pt(img.Width, img.Height),
		})
		for y := 0; y < img.Height; y++ {
			for x := 0; x < img.Width; x++ {
				pngImage.SetNRGBA(x, y, img.Pixels[y][x])
			}
		}
		list[i] = pngImage
	}
	return list
}

func ShpFromPNGList(images []image.Image) *ShpFile {
	f := &ShpFile{}
	f.Images = make([]ShpFileImage, len(images))

	clrModel := color.NRGBAModel

	for i, pngImage := range images {
		img := &f.Images[i]
		bounds := pngImage.Bounds()
		frameWidth := bounds.Dx()
		frameHeight := bounds.Dy()
		img.Width = frameWidth
		img.Height = frameHeight

		img.Pixels = make([][]color.NRGBA, frameHeight)
		for i := range img.Pixels {
			img.Pixels[i] = make([]color.NRGBA, frameWidth)
		}
		img.StartPixelX = img.Width + 0xffff
		img.StartPixelY = img.Width + 0xffff
		for y, row := range img.Pixels {
			for x := range row {
				clr := clrModel.Convert(pngImage.At(x, y)).(color.NRGBA)
				if clr.A > 0 {
					if y > img.EndPixelY {
						img.EndPixelY = y
					}
					if y < img.StartPixelY {
						img.StartPixelY = y
					}
					if x > img.EndPixelX {
						img.EndPixelX = x
					}
					if x < img.StartPixelX {
						img.StartPixelX = x
					}
				}
				img.Pixels[y][x] = clr
			}
		}
	}

	return f
}

func ShpFromPNG(filename string, pngImage image.Image) *ShpFile {
	f := &ShpFile{}

	// E.g. foo.3.png => 3 frames,
	// foo.png => 1 frames
	numFrames := 1
	if strings.Count(filename, ".") >= 2 {
		nameParts := strings.Split(filename, ".")
		v, err := strconv.Atoi(nameParts[len(nameParts)-2])
		if err == nil {
			numFrames = v
		}
	}
	f.Images = make([]ShpFileImage, numFrames)

	bounds := pngImage.Bounds()
	frameWidth := bounds.Dx() / numFrames
	frameHeight := bounds.Dy()

	clrModel := color.NRGBAModel

	offsetX := 0
	for i := range f.Images {
		img := &f.Images[i]
		img.Width = frameWidth
		img.Height = frameHeight

		img.Pixels = make([][]color.NRGBA, frameHeight)
		for i := range img.Pixels {
			img.Pixels[i] = make([]color.NRGBA, frameWidth)
		}
		img.StartPixelX = img.Width + 0xffff
		img.StartPixelY = img.Width + 0xffff
		for y, row := range img.Pixels {
			for x := range row {
				clr := clrModel.Convert(pngImage.At(x+offsetX, y)).(color.NRGBA)
				if clr.A > 0 {
					if y > img.EndPixelY {
						img.EndPixelY = y
					}
					if y < img.StartPixelY {
						img.StartPixelY = y
					}
					if x > img.EndPixelX {
						img.EndPixelX = x
					}
					if x < img.StartPixelX {
						img.StartPixelX = x
					}
				}
				img.Pixels[y][x] = clr
			}
		}
		offsetX += frameWidth
	}

	return f
}

func ParseShpFile(data []byte, pal *PaletteFile) (*ShpFile, error) {
	f := &ShpFile{
		UniSize:    true,
		UniPalette: true,
	}

	allData := data

	// The magic number is "1.10", an "eye catcher".
	if magic := string(data[:4]); magic != "1.10" {
		return nil, fmt.Errorf("invalid header, expected 1.10, got %q", magic)
	}
	data = data[len("1.10"):]

	numImages := scanUint32(&data)

	f.Images = make([]ShpFileImage, numImages)
	for i := range f.Images {
		img := &f.Images[i]
		img.DataOffset = int(scanUint32(&data))
		img.PaletteOffset = int(scanUint32(&data))
		img.NumPaletteColors = int(binary.LittleEndian.Uint32(allData[img.PaletteOffset:]))
	}

	for i := range f.Images {
		img := &f.Images[i]
		imgData := allData[img.DataOffset:]

		img.Height = int(scanUint16(&imgData)) + 1
		img.Width = int(scanUint16(&imgData)) + 1

		if img.Width != f.Images[0].Width || img.Height != f.Images[0].Height {
			f.UniSize = false
		}
		if img.PaletteOffset != f.Images[0].PaletteOffset {
			f.UniPalette = false
		}

		img.Pixels = make([][]color.NRGBA, img.Height)
		for i := range img.Pixels {
			img.Pixels[i] = make([]color.NRGBA, img.Width)
		}

		// Skip an unknown chunk.
		imgData = imgData[4:]

		img.StartPixelX = int(scanUint32(&imgData))
		img.StartPixelY = int(scanUint32(&imgData))
		img.EndPixelX = int(scanUint32(&imgData))
		img.EndPixelY = int(scanUint32(&imgData))

		isEmpty := img.StartPixelX > img.Width || img.StartPixelY > img.Height
		if isEmpty {
			// An empty image.
			continue
		}

		y := img.StartPixelY
		x := img.StartPixelX
		for y <= img.EndPixelY {
			op := imgData[0]
			advance := 1
			switch op {
			case 0:
				y++
				x = img.StartPixelX
			case 1:
				advance = 2
				x += int(imgData[1])
			default:
				if op%2 == 0 {
					// An even op byte means "draw a line".
					advance = 2
					lineLength := op / 2
					paletteIndex := int(imgData[1])
					clr := pal.Get(paletteIndex)
					for lineLength > 0 {
						img.Pixels[y][x] = clr
						x++
						lineLength--
					}
					break
				}
				// An odd op means "a variadic sequence of next pixels".
				numBytes := (int(op) - 1) / 2
				advance = numBytes + 1
				for i := 0; i < numBytes; i++ {
					paletteIndex := int(imgData[i+1])
					clr := pal.Get(paletteIndex)
					img.Pixels[y][x] = clr
					x++
				}
			}
			imgData = imgData[advance:]
		}
	}

	// I still have to figure out the local palette meaning.
	// It looks like there are always NumColors identical to
	// the number of main palette color usage.
	// R-component matches the main palette color.
	// Maybe G/B/A components have some meaning too?
	//
	// for index, order := range colorsUsed {
	// 	fmt.Printf("%d (%02x count=%d %x) => %02x %02x %02x\n", index, index, order, order, palette[index].R, palette[index].G, palette[index].B)
	// }
	// fmt.Println("num=", len(colorsUsed))
	// reds := make([]uint8, f.Images[0].NumPaletteColors)
	// for i := 0; i < f.Images[0].NumPaletteColors; i++ {
	// 	offset := f.Images[0].PaletteOffset + (i * 4) + 4
	// 	reds[i] = allData[offset]
	// }
	// for index := range colorsUsed {
	// 	// clr := palette[index]
	// 	i := gslices.Index(reds, uint8(index))
	// 	if i == -1 {
	// 		fmt.Printf("[!] missing %d\n", index)
	// 	} else {
	// 		fmt.Printf("[+] found %d\n", index)
	// 	}
	// }

	return f, nil
}
