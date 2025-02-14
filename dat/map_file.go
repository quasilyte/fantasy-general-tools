package dat

type MapFile struct {
	Cols  int
	Rows  int
	Tiles []RawTile

	// [ 0] - volcano/cities
	// [ 1] - plains/river
	// [ 2] - plains/mountains
	// [ 3] - plains/road
	// [ 4] - plains/clear
	// [ 5] - plains/cities
	// [ 6] - desert/clear
	// [ 7] - desert/road
	// [ 8] - desert/river
	// [ 9] - desert/mountains
	// [10] - desert/cities
	// [11] - volcano/river
	// [12] - volcano/road
	// [13] - snow/clear
	// [14] - snow/road
	// [15] - snow/mountains
	// [16] - snow/river
	// [17] - snow/cities
	// [18] - volcano/clear
	// [19] - volcano/mountains
	// [20] - jungle/clear
	// [21] - jungle/road
	// [22] - jungle/mountains
	// [23] - jungle/river
	// [24] - jungle/cities
	Tilebank [25]byte

	Names []string
}

func ParseMapFile(data []byte) (*MapFile, error) {
	return scanMapFile(&data)
}

func scanMapFile(data *[]byte) (*MapFile, error) {
	f := &MapFile{}

	f.Cols = int(scanUint16(data))
	f.Rows = int(scanUint16(data))

	tiles, err := scanTiles(f.Cols, f.Rows, data)
	if err != nil {
		return nil, err
	}
	f.Tiles = tiles

	// The next 25 bytes are related to a tile bank.
	copy(f.Tilebank[:], (*data)[:len(f.Tilebank)])
	*data = (*data)[len(f.Tilebank):]

	// Are next 75 bytes always 0?
	// It looks like they're the same part as [25]byte above,
	// but never used in the game.
	for i := 0; i < 75; i++ {
		if (*data)[i] != 0 {
			panic(i)
		}
	}
	*data = (*data)[75:]

	numNames := scanUint16(data)
	f.Names = make([]string, numNames)
	for i := range f.Names {
		// All names are 30-byte long.
		f.Names[i] = scanString(data, 30)
	}

	return f, nil
}
