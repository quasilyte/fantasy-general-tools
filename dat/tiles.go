package dat

type ImagePage byte

const (
	ImagePageGround    ImagePage = 0x06
	ImagePageRoad      ImagePage = 0x07
	ImagePageLake      ImagePage = 0x08
	ImagePageMountains ImagePage = 0x09
	ImagePageBuilding  ImagePage = 0x0a
)

const rawTileSize = 19

type RawTile struct {
	_          [14]byte  // unknown
	ImageIndex byte      // An image index inside page (e.g. frame)
	_          byte      // unknown, usually 0x80, could be related to ImageID (uint16?)
	ImagePage  ImagePage // An image page
	Kind       TileKind
	_          byte
}

type TileKind byte

const (
	TileClear TileKind = 0x03
	TileRoad  TileKind = 0x0d
)

func scanTiles(numCols, numRows int, data *[]byte) ([]RawTile, error) {
	numTiles := numCols * numRows

	tiles := make([]RawTile, numTiles)

	for i := range tiles {
		dst := &tiles[i]
		s := asBytes(dst)
		copy(s, (*data)[:rawTileSize])
		*data = (*data)[rawTileSize:]
	}

	return tiles, nil
}
