package dat

import (
	"bytes"
	"fmt"
)

type ScnFile struct {
	DeployTiles    []Coord
	VictoryFlags   []Coord
	Map            *MapFile
	TurnLimit      int
	SurrenderRate  int
	EnemyUnitLimit int
	EnemyBehavior  int

	EnemyUnits        []SceneUnit
	EnemyStartingGold int

	MortalTechLevelCap [10]uint8
	MagicTechLevelCap  [10]uint8
	BeastTechLevelCap  [10]uint8
	MechTechLevelCap   [10]uint8

	EnemyLeaderClass LeaderKind

	EnemyLeaderAdvantages [11]bool
}

const rawSceneUnitSize = 77

type SceneUnit struct {
	Name string

	Raw RawSceneUnit
}

type UnitState uint8

const (
	UnitStateDead       UnitState = 0
	UnitStateNormal     UnitState = 2
	UnitStateDisordered UnitState = 3
	UnitStateBroken     UnitState = 4
)

type RawSceneUnit struct {
	Pos               Coord
	Index             uint16le
	Experience        uint16le
	NumKills          uint16le
	_                 [8]byte // unknown, always 0 in all scenes
	Artifact          ArtifactKind
	Name              [30]byte
	IsLegal           byte // If 0, the game will crash with "trying to move illegal unit"
	Level             byte
	Morale            byte
	NumDead           byte // Initially dead units, 0 means "everyone alive"
	NumWounded        byte // Initially wounded units, 0 means "no wounded"
	State             UnitState
	Subclass          byte // Affects the unit type; TODO: figure out the effects of this
	StatusFatigue     byte
	StatusDeathWounds byte
	StatusInvisible   byte
	StatusForceMarch  byte
	StatusArmor       byte
	StatusBerserk     byte
	StatusPathFinder  byte
	StatusBless       byte
	_                 byte    // unknown
	_                 byte    // unknown
	BoostSearchRange  byte    // Adds to search score (described as "with magic")
	BoostMelee        byte    // Adds to melee score
	BoostArmor        byte    // Adds to armor rating (described as "with magic")
	BoostAttacks      byte    // Adds to num of attacks
	BoostLife         byte    // Adds to current and max life
	BoostMorale       byte    // Adds to current and max morale
	_                 [5]byte // unknown
}

type LeaderKind uint8

const (
	LeaderHealer LeaderKind = iota
	LeaderWarlord
	LeaderArchmage
	LeaderBeastmaster
)

type LeaderAdvantageKind uint8

const (
	LeaderAdvantageFame LeaderAdvantageKind = iota
	LeaderAdvantageWarlord
	LeaderAdvantageGeneral
	LeaderAdvantageCharisma
	LeaderAdvantageHealer
	LeaderAdvantageEngineer
	LeaderAdvantageArchMage
	LeaderAdvantageMechMaster
	LeaderAdvantageBeastMaster
	LeaderAdvantageRanger
	LeaderAdvantageDwarven
)

type Coord struct {
	Col uint16le
	Row uint16le
}

func ParseScnFile(data []byte) (*ScnFile, error) {
	f := &ScnFile{}

	// The last 12 bytes encode enemy leader.
	enemyLeaderDataStart := data[len(data)-12:]

	mapDataStart := data[1+9218:]

	// The first byte is weird: everything breaks if
	// it's zero. Maybe it's a flag for "fixed map"
	// (as opposed to generated map)?
	// TODO: SCEN198.SCN has it set to 1!
	isFixedMap := scanUint8(&data)
	if isFixedMap != 1 {
		return nil, fmt.Errorf("expected the first byte to be 0x01, found %02x", isFixedMap)
	}

	numDeployTiles := scanUint16(&data)
	f.DeployTiles = make([]Coord, numDeployTiles)
	for i := range f.DeployTiles {
		col := scanUint16(&data)
		row := scanUint16(&data)
		f.DeployTiles[i].Col.SetInt(int(col))
		f.DeployTiles[i].Row.SetInt(int(row))
	}

	// Fast-forward to relevant data.
	data = mapDataStart

	mapfile, err := scanMapFile(&data, true)
	if err != nil {
		return nil, fmt.Errorf("map data: %v", err)
	}
	f.Map = mapfile

	// From here, it's 409 bytes forward until the next data block.
	// Everything between this and the next block is filled with 0xff.
	nextDataChunk := data[409:]

	// The current section has variable-length data (but it never exceeds
	// the limits described above.

	// TODO: analyze this section.
	// It has something to do with shrine rewards.
	// data[4] in SCEN2 can change the loot table

	// After we handled the section, continue
	// from the next section to skip all 0xff in-between.
	data = nextDataChunk

	numFlags := scanUint8(&data)
	if data[0] != 0 {
		panic("sanity check failed: second numFlags byte is not zero")
	}
	// Could be the second byte for numFlags, unsure.
	// I'm using it above as a sanity check as the game
	// can't have more than 10 flags anyway.
	_ = data[1:]

	nextDataChunk = data[2*2*10:] // After capture flag coord pairs

	f.VictoryFlags = make([]Coord, numFlags)
	for i := range f.VictoryFlags {
		col := scanUint16(&data)
		row := scanUint16(&data)
		f.VictoryFlags[i].Col.SetInt(int(col))
		f.VictoryFlags[i].Row.SetInt(int(row))
	}

	data = nextDataChunk

	// The following byte is always 0.
	// TODO: can we still figure it out?
	data = data[1:]

	// This one has 2 possible values: 1 and 5.
	// Not even 0 can be there.
	// Only one map sets it to 5: SCEN43.
	// Not sure about its usage yet. TODO
	data = data[1:]

	// The next two bytes only take two values: 0 and 1.
	// Most likely some kind of bool flags.
	// TODO: figure them out.
	// The second byte seem to control the enemy surrender and win/lose states.
	data = data[2:]

	f.TurnLimit = int(scanUint8(&data))

	// These 8 bytes are always 0.
	// Can't study them.
	if !bytes.Equal([]byte{0, 0, 0, 0, 0, 0, 0, 0}, data[:8]) {
		panic("sanity check failed: unknown 8 bytes are not zeroed")
	}
	data = data[8:]

	// A % of troops a computer needs to lose before surrendering.
	// A campaign never uses a value lower than 70.
	// Setting it too low may cause the player to lose after accepting surrender (!)
	f.SurrenderRate = int(scanUint8(&data))
	if f.SurrenderRate > 100 {
		panic("sanity check failed: surrender rate exceeds 100%")
	}
	// The next byte is always 0.
	// Could be the second byte of SurrenderRate, but it's not needed.
	if data[0] != 0 {
		panic("sanity check failed: second surrenderRate byte is not zero")
	}
	data = data[1:]

	f.EnemyStartingGold = int(scanUint16(&data))

	// These 2 bytes are always 0.
	// Setting them grants enemy infinite (?) gold.
	// I would assume them to be a gold income, but
	// setting anything here allows the enemy to buy tons of units.
	if !bytes.Equal([]byte{0, 0}, data[:2]) {
		panic("sanity check failed: unknown 2 bytes are not zeroed")
	}
	data = data[2:]

	// Then there is an enemy's name.
	// Probably 30 bytes-long, as usual.
	// It's not important as the game doesn't load this string
	// even if it doesn't match the enemy name.
	data = data[30:]

	// Unsure about this one, but it seem to be increasing
	// further into campaign. Starts small, gets bigger.
	// It it AI level?
	f.EnemyBehavior = int(scanUint8(&data))

	// These 3 bytes are always 0.
	// Can't study them.
	if !bytes.Equal([]byte{0, 0, 0}, data[:3]) {
		panic("sanity check failed: unknown 3 bytes are not zeroed")
	}
	data = data[3:]

	// Only a few weird scenes (not part of campaign I think)
	// use these. Could be hard to figure out.
	// Namely, SCEN250.SCN SCEN251.SCN SCEN252.SCN SCEN253.SCN SCEN254.SCN SCEN255.SCN
	data = data[4:]

	// The 40 next byte has these values: 0, 1, 2, 3, 4, 5.
	// Probably the max tech level for every unit and every race.
	techCaps := []*[10]uint8{
		&f.MortalTechLevelCap,
		&f.MagicTechLevelCap,
		&f.BeastTechLevelCap,
		&f.MechTechLevelCap,
	}
	for _, array := range techCaps {
		for j := range array {
			(*array)[j] = scanUint8(&data)
		}
	}

	f.EnemyUnitLimit = int(scanUint8(&data))
	f.EnemyUnits = make([]SceneUnit, 0, 8)
	for {
		var u RawSceneUnit
		copy(asBytes(&u), data)
		name := trimCstring(u.Name[:])
		if name == "" {
			break
		}
		f.EnemyUnits = append(f.EnemyUnits, SceneUnit{
			Name: name,
			Raw:  u,
		})
		data = data[rawSceneUnitSize:]
	}

	data = enemyLeaderDataStart

	f.EnemyLeaderClass = LeaderKind(scanUint8(&data))
	// 11 leader advantages as byte-wide flags.
	for i := range f.EnemyLeaderAdvantages {
		f.EnemyLeaderAdvantages[i] = scanBool(&data)
	}

	return f, nil
}
