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
	NameIndex  byte      // An index inside the map's name array
	_          [13]byte  // unknown
	ImageIndex byte      // An image index inside page (e.g. frame)
	_          byte      // usually 0x80 in MAP and 0 in SCN
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
