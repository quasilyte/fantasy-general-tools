package dat

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/quasilyte/fantasy-general-tools/serdat"
	"github.com/quasilyte/gslices"
)

const unitSize = 74

type MagequipFile struct {
	Units []*MagequipUnit
}

type MagequipUnit struct {
	Name string

	BuyPrice int

	Index int

	Raw RawMagequipUnit
}

func (u *MagequipUnit) ToSerdat() serdat.MagequipUnit {
	return serdat.MagequipUnit{
		Index:              u.Index,
		ImageID:            int(u.Raw.ImageID),
		BuyPrice:           int(u.BuyPrice),
		MeleeDamage:        int(u.Raw.MeleeDamage),
		MeleeDamageType:    u.Raw.MeleeDamageType.String(),
		SkirmishDamage:     int(u.Raw.SkirmishDamage),
		SkirmishDamageType: u.Raw.SkirmishDamageType.String(),
		MissileDamage:      int(u.Raw.MissileDamage),
		MissileDamageType:  u.Raw.MissileDamageType.String(),
		SiegeDamage:        int(u.Raw.SiegeDamage),
		SiegeDamageType:    u.Raw.SiegeDamageType.String(),
		Side:               u.Raw.Side.String(),
		UnitClass:          u.Raw.UnitClass.String(),
		Name:               u.Name,
		Life:               int(u.Raw.Life),
		NumAttacks:         int(u.Raw.NumAttacks),
		Armor:              int(u.Raw.Armor),
		MagicDefense:       int(u.Raw.MagicDefense),
		Speed:              int(u.Raw.Speed),
		SearchRange:        int(u.Raw.SearchRange),
		Spell:              u.Raw.Spell.String(),
		MovementType:       u.Raw.MovementType.String(),
		Ability1:           u.Raw.Ability1.String(),
		Ability2:           u.Raw.Ability2.String(),
		Ability3:           u.Raw.Ability3.String(),
		Race:               u.Raw.Race.String(),
		TechLevel:          int(u.Raw.TechLevel),
	}
}

type RawMagequipUnit struct {
	ImageID  byte
	Unknown2 byte // Could be the second byte of image ID

	BuyPrice uint16le

	MeleeDamage     byte
	Unknown5        byte // Always seem to be zero
	MeleeDamageType DamageType
	Unknown7        byte // Always seem to be zero

	SkirmishDamage     byte
	Unknown9           byte // Always seem to be zero
	SkirmishDamageType DamageType
	Unknown11          byte // Always seem to be zero

	MissileDamage     byte
	Unknown13         byte // Always seem to be zero
	MissileDamageType DamageType
	Unknown15         byte // Always seem to be zero

	SiegeDamage     byte
	Unknown17       byte // Always seem to be zero
	SiegeDamageType DamageType
	Unknown19       byte // Always seem to be zero

	Side         UnitSide
	Unknown21    byte
	Unknown22    byte
	Unknown23    byte
	Unknown24    byte
	Unknown25    byte
	Unknown26    byte
	Unknown27    byte
	Unknown28    byte
	Unknown29    byte
	UnitClass    UnitClass
	Name         [30]byte
	Life         byte
	NumAttacks   byte
	Armor        byte
	MagicDefense byte
	Speed        byte
	SearchRange  byte
	Spell        UnitSpell
	MovementType UnitMovementType
	Ability1     UnitAbility
	Ability2     UnitAbility
	Ability3     UnitAbility
	Race         UnitRace
	TechLevel    byte
}

func (u *MagequipUnit) CommentString() string {
	raw := &u.Raw
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("    0) [%02x] -- image ID\n", raw.ImageID))
	buf.WriteString(fmt.Sprintf("    1) [%02x] -- ?\n", raw.Unknown2))
	buf.WriteString(fmt.Sprintf("  2-3) [%02x][%02x] -- buying price (%d$)\n", raw.BuyPrice.Low, raw.BuyPrice.High, u.BuyPrice))
	buf.WriteString(fmt.Sprintf("    4) [%02x] -- melee damage (%v)\n", raw.MeleeDamage, raw.MeleeDamage))
	buf.WriteString(fmt.Sprintf("    5) [%02x] -- ?\n", raw.Unknown5))
	buf.WriteString(fmt.Sprintf("    6) [%02x] -- melee damage type (%v)\n", byte(raw.MeleeDamageType), raw.MeleeDamageType))
	buf.WriteString(fmt.Sprintf("    7) [%02x] -- ?\n", raw.Unknown7))
	buf.WriteString(fmt.Sprintf("    8) [%02x] -- skirmish damage (%v)\n", raw.SkirmishDamage, raw.SkirmishDamage))
	buf.WriteString(fmt.Sprintf("    9) [%02x] -- ?\n", raw.Unknown9))
	buf.WriteString(fmt.Sprintf("   10) [%02x] -- skirmish damage type (%v)\n", byte(raw.SkirmishDamageType), raw.SkirmishDamageType))
	buf.WriteString(fmt.Sprintf("   11) [%02x] -- ?\n", raw.Unknown11))
	buf.WriteString(fmt.Sprintf("   12) [%02x] -- missile damage (%v)\n", byte(raw.MissileDamage), raw.MissileDamage))
	buf.WriteString(fmt.Sprintf("   13) [%02x] -- ?\n", raw.Unknown13))
	buf.WriteString(fmt.Sprintf("   14) [%02x] -- missile damage type (%v)\n", byte(raw.MissileDamageType), raw.MissileDamageType))
	buf.WriteString(fmt.Sprintf("   15) [%02x] -- ?\n", raw.Unknown15))
	buf.WriteString(fmt.Sprintf("   16) [%02x] -- siege damage (%v)\n", raw.SiegeDamage, raw.SiegeDamage))
	buf.WriteString(fmt.Sprintf("   17) [%02x] -- ?\n", raw.Unknown17))
	buf.WriteString(fmt.Sprintf("   18) [%02x] -- siege damage type (%v)\n", byte(raw.SiegeDamageType), raw.SiegeDamageType))
	buf.WriteString(fmt.Sprintf("   19) [%02x] -- ?\n", raw.Unknown19))
	buf.WriteString(fmt.Sprintf("   20) [%02x] -- side (%v)\n", byte(raw.Side), raw.Side))
	buf.WriteString(fmt.Sprintf("   21) [%02x] -- ?\n", raw.Unknown21))
	buf.WriteString(fmt.Sprintf("   22) [%02x] -- ?\n", raw.Unknown22))
	buf.WriteString(fmt.Sprintf("   23) [%02x] -- ?\n", raw.Unknown23))
	buf.WriteString(fmt.Sprintf("   24) [%02x] -- ?\n", raw.Unknown24))
	buf.WriteString(fmt.Sprintf("   25) [%02x] -- ?\n", raw.Unknown25))
	buf.WriteString(fmt.Sprintf("   26) [%02x] -- ?\n", raw.Unknown26))
	buf.WriteString(fmt.Sprintf("   27) [%02x] -- ?\n", raw.Unknown27))
	buf.WriteString(fmt.Sprintf("   28) [%02x] -- ?\n", raw.Unknown28))
	buf.WriteString(fmt.Sprintf("   29) [%02x] -- ?\n", raw.Unknown29))
	buf.WriteString(fmt.Sprintf("   30) [%02x] -- unit type (%v)\n", byte(raw.UnitClass), raw.UnitClass))
	buf.WriteString(fmt.Sprintf("31-60) [%02x]...[%02x] -- unit name (%q)\n",
		u.Name[0], u.Name[len(u.Name)-1], u.Name))
	buf.WriteString(fmt.Sprintf("   61) [%02x] -- life (%v)\n", raw.Life, raw.Life))
	buf.WriteString(fmt.Sprintf("   62) [%02x] -- attacks (%v)\n", raw.NumAttacks, raw.NumAttacks))
	buf.WriteString(fmt.Sprintf("   63) [%02x] -- armor (%v)\n", raw.Armor, raw.Armor))
	buf.WriteString(fmt.Sprintf("   64) [%02x] -- magic defense (%v)\n", raw.MagicDefense, raw.MagicDefense))
	buf.WriteString(fmt.Sprintf("   65) [%02x] -- speed (%v)\n", raw.Speed, raw.Speed))
	buf.WriteString(fmt.Sprintf("   66) [%02x] -- search range (%v)\n", raw.SearchRange, raw.SearchRange))
	if raw.Spell != 0 {
		buf.WriteString(fmt.Sprintf("   67) [%02x] -- spell (%v)\n", byte(raw.Spell), raw.Spell))
	} else {
		buf.WriteString(fmt.Sprintf("   67) [%02x] -- spell\n", byte(raw.Spell)))
	}
	buf.WriteString(fmt.Sprintf("   68) [%02x] -- movement type (%v)\n", byte(raw.MovementType), raw.MovementType))
	if raw.Ability1 != 0 {
		buf.WriteString(fmt.Sprintf("   69) [%02x] -- ability A (%v)\n", byte(raw.Ability1), raw.Ability1))
	} else {
		buf.WriteString(fmt.Sprintf("   69) [%02x] -- ability A\n", byte(raw.Ability1)))
	}
	if raw.Ability2 != 0 {
		buf.WriteString(fmt.Sprintf("   70) [%02x] -- ability B (%v)\n", byte(raw.Ability2), raw.Ability2))
	} else {
		buf.WriteString(fmt.Sprintf("   70) [%02x] -- ability B\n", byte(raw.Ability2)))
	}
	if raw.Ability3 != 0 {
		buf.WriteString(fmt.Sprintf("   71) [%02x] -- ability C (%v)\n", byte(raw.Ability3), raw.Ability3))
	} else {
		buf.WriteString(fmt.Sprintf("   71) [%02x] -- ability C\n", byte(raw.Ability3)))
	}
	buf.WriteString(fmt.Sprintf("   72) [%02x] -- race (%v)\n", byte(raw.Race), raw.Race))

	return buf.String()
}

type UnitRace byte

const (
	UnitRaceMortal UnitRace = iota
	UnitRaceMagic
	UnitRaceBeast
	UnitRaceMech
)

var unitRaceStrings = []string{
	UnitRaceMortal: "mortal",
	UnitRaceMagic:  "magic",
	UnitRaceBeast:  "beast",
	UnitRaceMech:   "mech",
}

func (race UnitRace) String() string {
	return unitRaceStrings[race]
}

type UnitAbility byte

const (
	UnitAbilityNone               UnitAbility = 0
	UnitAbilityNegateCavalryBonus UnitAbility = 1
	UnitAbilityVulnerableToMech   UnitAbility = 2
	UnitAbilitySpellResistant     UnitAbility = 3
	UnitAbilityCantBeHealed       UnitAbility = 4
	UnitAbilitySingleEntity       UnitAbility = 5
	UnitAbilityRegeneration       UnitAbility = 6
	UnitAbilityHeroic             UnitAbility = 7
	UnitAbilityHero               UnitAbility = 8
	UnitAbilityUndead             UnitAbility = 9
	UnitAbilityBloodlust          UnitAbility = 10
	UnitAbilityVulnerableToMagic  UnitAbility = 11
	UnitAbilityMystic             UnitAbility = 12
	UnitAbilityBerserk            UnitAbility = 13
)

var unitAbilityStrings = []string{
	UnitAbilityNone:               "",
	UnitAbilityNegateCavalryBonus: "negate_cavalry_bonus",
	UnitAbilityVulnerableToMech:   "vulnerable_to_mech",
	UnitAbilitySpellResistant:     "spell_resistant",
	UnitAbilityCantBeHealed:       "cant_be_healed",
	UnitAbilitySingleEntity:       "single_entity",
	UnitAbilityRegeneration:       "regeneration",
	UnitAbilityHeroic:             "heroic",
	UnitAbilityHero:               "hero",
	UnitAbilityUndead:             "undead",
	UnitAbilityBloodlust:          "bloodlust",
	UnitAbilityVulnerableToMagic:  "vulnerable_to_magic",
	UnitAbilityMystic:             "mystic",
	UnitAbilityBerserk:            "berserk",
}

func (ability UnitAbility) String() string {
	return unitAbilityStrings[ability]
}

type UnitMovementType byte

const (
	UnitMovementLightInfantry UnitMovementType = iota
	UnitMovementHeavyInfantry
	UnitMovementCavalry
	UnitMovementSiegeEngine
	UnitMovementFlying
)

var unitMovementStrings = []string{
	UnitMovementLightInfantry: "light_infantry",
	UnitMovementHeavyInfantry: "heavy_infantry",
	UnitMovementCavalry:       "cavalry",
	UnitMovementSiegeEngine:   "siege_engine",
	UnitMovementFlying:        "flying",
}

func (movement UnitMovementType) String() string {
	return unitMovementStrings[movement]
}

type UnitSpell byte

const (
	UnitSpellNone UnitSpell = iota
	UnitSpellDeathWounds
	UnitSpellSummonBanshee
	UnitSpellFirestorm
	UnitSpellAnimateTrees
	UnitSpellRaiseDead
	UnitSpellHeal
	UnitSpellEarthquake
	UnitSpellEagleEye
	UnitSpellInvisibility
	UnitSpellForceMarch
	UnitSpellArmor
	UnitSpellDispelUndead
	UnitSpellBerserk
	UnitSpellCharm
	UnitSpellRestoreMorale
	UnitSpellPathMaster
	UnitSpellPanic
	UnitSpellWhirlwind
	UnitSpellBless
	UnitSpellStorm          // From the item
	UnitSpellSickness       // From the item
	UnitSpellDrainLife      // ?
	UnitSpellConfusion      // From the item
	UnitSpellFatigue        // From the item
	UnitSpellShortFireball  // From the item
	UnitSpellFireball       // From the item
	UnitSpellHeroMassHeal   // A hero spell
	UnitSpellHeroFireball   // A hero spell
	UnitSpellHeroFear       // A hero spell
	UnitSpellHeroFatigue    // A hero spell
	UnitSpellHeroWhirlwind  // A hero spell
	UnitSpellHeroHealBeasts // A hero spell
	UnitSpellHeroPlague     // A hero spell
)

var unitSpellStrings = []string{
	UnitSpellNone:           "",
	UnitSpellDeathWounds:    "death_wounds",
	UnitSpellSummonBanshee:  "summon_banshee",
	UnitSpellFirestorm:      "firestorm",
	UnitSpellAnimateTrees:   "animate_trees",
	UnitSpellRaiseDead:      "raise_dead",
	UnitSpellHeal:           "heal",
	UnitSpellEarthquake:     "earthquake",
	UnitSpellEagleEye:       "eagle_eye",
	UnitSpellInvisibility:   "invisibility",
	UnitSpellForceMarch:     "force_march",
	UnitSpellArmor:          "armor",
	UnitSpellDispelUndead:   "dispel_undead",
	UnitSpellBerserk:        "berserk",
	UnitSpellCharm:          "charm",
	UnitSpellRestoreMorale:  "restore_morale",
	UnitSpellPathMaster:     "path_master",
	UnitSpellPanic:          "panic",
	UnitSpellWhirlwind:      "whirlwind",
	UnitSpellBless:          "bless",
	UnitSpellStorm:          "storm",
	UnitSpellSickness:       "sickness",
	UnitSpellDrainLife:      "drain_life",
	UnitSpellConfusion:      "confusion",
	UnitSpellFatigue:        "fatigue",
	UnitSpellShortFireball:  "short_fireball",
	UnitSpellFireball:       "fireball",
	UnitSpellHeroMassHeal:   "hero_mass_heal",
	UnitSpellHeroFireball:   "hero_fireball",
	UnitSpellHeroFear:       "hero_fear",
	UnitSpellHeroFatigue:    "hero_fatigue",
	UnitSpellHeroWhirlwind:  "hero_whirlwind",
	UnitSpellHeroHealBeasts: "hero_heal_beasts",
	UnitSpellHeroPlague:     "hero_plague",
}

func (spell UnitSpell) String() string {
	return unitSpellStrings[spell]
}

type UnitClass byte

const (
	UnitClassHeavyInfantry UnitClass = 0
	UnitClassLightInfantry UnitClass = 1
	UnitClassSkirmisher    UnitClass = 2
	UnitClassArcher        UnitClass = 3
	UnitClassCavalry       UnitClass = 4
	UnitClassLightCavalry  UnitClass = 5
	UnitClassSkyHunter     UnitClass = 6
	UnitClassBomber        UnitClass = 7
	UnitClassSiegeEngine   UnitClass = 8
	UnitClassSpellCaster   UnitClass = 9
)

var unitClassStrings = []string{
	UnitClassHeavyInfantry: "heavy_infantry",
	UnitClassLightInfantry: "light_infantry",
	UnitClassSkirmisher:    "skirmisher",
	UnitClassArcher:        "archer",
	UnitClassCavalry:       "cavalry",
	UnitClassLightCavalry:  "light_cavalry",
	UnitClassSkyHunter:     "sky_hunter",
	UnitClassBomber:        "bomber",
	UnitClassSiegeEngine:   "siege_engine",
	UnitClassSpellCaster:   "spell_caster",
}

func (c UnitClass) String() string {
	return unitClassStrings[c]
}

type DamageType byte

const (
	DamageNormal DamageType = iota
	DamageMagic
	DamageMech
)

var damageTypeStrings = []string{
	DamageNormal: "normal",
	DamageMagic:  "magic",
	DamageMech:   "mech",
}

func (t DamageType) String() string {
	return damageTypeStrings[t]
}

type UnitSide byte

const (
	UnitSideGood UnitSide = iota
	UnitSideEvil
	UnitSideNeutral
)

var unitSideStrings = []string{
	UnitSideGood:    "good",
	UnitSideEvil:    "evil",
	UnitSideNeutral: "neutral",
}

func (side UnitSide) String() string {
	return unitSideStrings[side]
}

func MagequipEncode(f *MagequipFile) ([]byte, error) {
	data := make([]byte, 0, 2+(len(f.Units)*unitSize))

	numUnits := len(f.Units)
	data = append(data,
		byte(numUnits),
		byte(numUnits>>8))

	for _, u := range f.Units {
		data = append(data, asBytes(&u.Raw)...)
	}

	if len(data) != cap(data) {
		return nil, fmt.Errorf("encoding error: cap and len mismatch (%d vs %d)", len(data), cap(data))
	}

	return data, nil
}

func MagequipFromSerdat(units []serdat.MagequipUnit) (*MagequipFile, error) {
	f := &MagequipFile{
		Units: make([]*MagequipUnit, 0, len(units)),
	}

	for _, u := range units {
		name := u.Name
		u, err := magequipUnitFromSerdat(u)
		if err != nil {
			return nil, fmt.Errorf("unit %q: %v", name, err)
		}
		f.Units = append(f.Units, &u)
	}

	gslices.SortFunc(f.Units, func(x, y *MagequipUnit) bool {
		return x.Index < y.Index
	})

	return f, nil
}

func ParseMagequipFile(filename string) (*MagequipFile, error) {
	f := &MagequipFile{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	numExpected := int(data[0]) | (int(data[1]) << 8)

	// Skip the header.
	data = data[2:]

	for i := 0; i < numExpected; i++ {
		chunk := data[:unitSize]

		unitData := [unitSize]byte(chunk)
		raw := *(*RawMagequipUnit)(unsafe.Pointer(&unitData))
		u := &MagequipUnit{
			Raw:   raw,
			Index: i,
		}
		nameLength := 0
		for _, ch := range raw.Name {
			if ch == 0 {
				break
			}
			nameLength++
		}
		u.Name = string(raw.Name[:nameLength])

		u.BuyPrice = raw.BuyPrice.ToInt()
		f.Units = append(f.Units, u)

		data = data[unitSize:]
	}

	if len(f.Units) != numExpected {
		return nil, fmt.Errorf("bad file: expected %d units, got %d", numExpected, len(f.Units))
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("bad fine: unconsumed %d bytes", len(data))
	}

	return f, nil
}

func magequipUnitFromSerdat(data serdat.MagequipUnit) (MagequipUnit, error) {
	u := MagequipUnit{
		Name:     data.Name,
		BuyPrice: data.BuyPrice,
		Index:    data.Index,
		Raw: RawMagequipUnit{
			ImageID:        byte(data.ImageID),
			MeleeDamage:    byte(data.MeleeDamage),
			SkirmishDamage: byte(data.SkirmishDamage),
			MissileDamage:  byte(data.MissileDamage),
			SiegeDamage:    byte(data.SiegeDamage),
			Life:           byte(data.Life),
			NumAttacks:     byte(data.NumAttacks),
			Armor:          byte(data.Armor),
			MagicDefense:   byte(data.MagicDefense),
			Speed:          byte(data.Speed),
			SearchRange:    byte(data.SearchRange),
			TechLevel:      byte(data.TechLevel),
		},
	}

	u.Raw.BuyPrice.SetInt(data.BuyPrice)

	if v := gslices.Index(damageTypeStrings, data.MeleeDamageType); v != -1 {
		u.Raw.MeleeDamageType = DamageType(v)
	} else {
		return u, fmt.Errorf("invalid melee damage type: %q", data.MeleeDamageType)
	}

	if v := gslices.Index(damageTypeStrings, data.SkirmishDamageType); v != -1 {
		u.Raw.SkirmishDamageType = DamageType(v)
	} else {
		return u, fmt.Errorf("invalid skirmish damage type: %q", data.SkirmishDamageType)
	}

	if v := gslices.Index(damageTypeStrings, data.MissileDamageType); v != -1 {
		u.Raw.MissileDamageType = DamageType(v)
	} else {
		return u, fmt.Errorf("invalid missile damage type: %q", data.MissileDamageType)
	}

	if v := gslices.Index(damageTypeStrings, data.SiegeDamageType); v != -1 {
		u.Raw.SiegeDamageType = DamageType(v)
	} else {
		return u, fmt.Errorf("invalid siege damage type: %q", data.SiegeDamageType)
	}

	if v := gslices.Index(unitSideStrings, data.Side); v != -1 {
		u.Raw.Side = UnitSide(v)
	} else {
		return u, fmt.Errorf("invalid side: %q", data.Side)
	}

	if v := gslices.Index(unitClassStrings, data.UnitClass); v != -1 {
		u.Raw.UnitClass = UnitClass(v)
	} else {
		return u, fmt.Errorf("invalid class: %q", data.UnitClass)
	}

	if v := gslices.Index(unitSpellStrings, data.Spell); v != -1 {
		u.Raw.Spell = UnitSpell(v)
	} else {
		return u, fmt.Errorf("invalid spell: %q", data.Spell)
	}

	if v := gslices.Index(unitMovementStrings, data.MovementType); v != -1 {
		u.Raw.MovementType = UnitMovementType(v)
	} else {
		return u, fmt.Errorf("invalid movement type: %q", data.MovementType)
	}

	if v := gslices.Index(unitAbilityStrings, data.Ability1); v != -1 {
		u.Raw.Ability1 = UnitAbility(v)
	} else {
		return u, fmt.Errorf("invalid ability 1 type: %q", data.Ability1)
	}
	if v := gslices.Index(unitAbilityStrings, data.Ability2); v != -1 {
		u.Raw.Ability2 = UnitAbility(v)
	} else {
		return u, fmt.Errorf("invalid ability 2 type: %q", data.Ability2)
	}
	if v := gslices.Index(unitAbilityStrings, data.Ability3); v != -1 {
		u.Raw.Ability3 = UnitAbility(v)
	} else {
		return u, fmt.Errorf("invalid ability 3 type: %q", data.Ability3)
	}

	if v := gslices.Index(unitRaceStrings, data.Race); v != -1 {
		u.Raw.Race = UnitRace(v)
	} else {
		return u, fmt.Errorf("invalid race: %q", data.Race)
	}

	if len(data.Name) > 29 {
		return u, fmt.Errorf("unit name is too long")
	}

	copy(u.Raw.Name[:], []byte(u.Name))

	return u, nil
}
