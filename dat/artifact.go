package dat

type ArtifactKind byte

const (
	ArtifactNone                 ArtifactKind = 0
	ArtifactMagicalSword         ArtifactKind = 1
	ArtifactMagicalAxe           ArtifactKind = 2
	ArtifactMagicalShield        ArtifactKind = 3
	ArtifactMagicalArmor         ArtifactKind = 4
	ArtifactMagicalBow           ArtifactKind = 5
	ArtifactMagicalCrossbow      ArtifactKind = 6
	ArtifactMagicalSling         ArtifactKind = 7
	ArtifactMagicalStaffSling    ArtifactKind = 8
	ArtifactMagicalJavelin       ArtifactKind = 9
	ArtifactAmuletOfResistance   ArtifactKind = 10
	ArtifactCloakOfHiding        ArtifactKind = 11
	ArtifactStaffOfFireballs     ArtifactKind = 12 // 3-dist fireball
	ArtifactWandOfFireballs      ArtifactKind = 13 // 1-dist fireball
	ArtifactAmuletOfProtection   ArtifactKind = 14
	ArtifactRingOfRegeneration   ArtifactKind = 15
	ArtifactSwordOfDeathWounds   ArtifactKind = 16
	ArtifactStoneOfEagleEye      ArtifactKind = 17
	ArtifactBootsOfRapidMovement ArtifactKind = 18
	ArtifactStormStaff           ArtifactKind = 19
	ArtifactAmuletOfBlessings    ArtifactKind = 20
	ArtifactRingOfFreeMovement   ArtifactKind = 21
	ArtifactStaffOfWinds         ArtifactKind = 22
	ArtifactWandOfSickness       ArtifactKind = 23 // 1-dist weakness
	ArtifactSkullOfCrang         ArtifactKind = 24 // 3-dist weakness
	ArtifactStaffOfConfusion     ArtifactKind = 25
	ArtifactStaffOfFatigue       ArtifactKind = 26
	ArtifactStaffOfHealing       ArtifactKind = 27
)
