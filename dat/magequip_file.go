package dat

import (
	"fmt"
	"os"
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
		u.Name = trimCstring(raw.Name[:])

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
