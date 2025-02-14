package serdat

type MagequipUnit struct {
	Index     int
	ImageID   int
	TechLevel int
	Race      string
	Side      string
	BuyPrice  int

	MeleeDamage     int
	MeleeDamageType string

	SkirmishDamage     int
	SkirmishDamageType string

	MissileDamage     int
	MissileDamageType string

	SiegeDamage     int
	SiegeDamageType string

	UnitClass    string
	Name         string
	Life         int
	NumAttacks   int
	Armor        int
	MagicDefense int
	Speed        int
	MovementType string
	SearchRange  int

	Spell    string
	Ability1 string
	Ability2 string
	Ability3 string
}
