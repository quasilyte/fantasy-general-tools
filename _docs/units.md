# Unit Data

`MAGEQUIP.DAT` file contains every unit definition.

Upon decoding, these definitions are separated into JSON files.

These JSON files look like this:

```
{
  "Index": 120,
  "ImageID": 110,
  "TechLevel": 0,
  "Race": "mortal",
  "Side": "good",
  "BuyPrice": 205,
  "MeleeDamage": 20,
  "MeleeDamageType": "mech",
  "SkirmishDamage": 12,
  "SkirmishDamageType": "mech",
  "MissileDamage": 10,
  "MissileDamageType": "mech",
  "SiegeDamage": 0,
  "SiegeDamageType": "normal",
  "UnitClass": "heavy_infantry",
  "Name": "Dreggo",
  "Life": 10,
  "NumAttacks": 10,
  "Armor": 20,
  "MagicDefense": 50,
  "Speed": 3,
  "MovementType": "light_infantry",
  "SearchRange": 3,
  "Spell": "",
  "Ability1": "hero",
  "Ability2": "",
  "Ability3": ""
}
```

You can `decode` the DAT-file into a bunch of JSONs, edit those files and then `encode` it back. When the original DAT-file is replaced, the game should reflect your changes.

The auto-generated file names are not important, but they reflect the `Index` property. Dreggo would be named `120_Dreggo.json`. This index is used as a unit ID inside all other files (maps, etc.). You may not want to change the existing unit's `Index` unless you know what you're doing.

Instead, it is advised to change units in-place if you want to replace them with something else (without changing indices). Or you can add a new unit with a new `Index`.

`ImageID` references associated graphical resources. For our Dreggo example, there is `SHP/ANIM120.SHP` which contains death animation. The idle/icon image is located in `SHP/UNITICON.SHP` at `ImageID` offset (it's a combined spritesheet). When decoded, `SHP/UNITICON.SHP` is turned into individual PNG files and naming becomes important (a name is an index/offset of this file when converted back into SHP).

`TechLevel` is required research level to recruit this unit. Heroes are not recruitable normally, so they have a level of 0. Heavy spearmen has a level of 5, for example. Normally, the game limits a max tech level value to 5. Some other examples: Peasants is level 0, Swordsmen is level 1. If there is a level skip in a tech tree for this unit kind, the next available unit will be researched instead. And example of this is Lionmen (lvl3) becoming Elephantmen (lvl5) after a single research round. If there are more than 1 lvlX unit, all of them will be available after 1 unit of this level is researched. This means having more than one unit inside a tech level is valid for a campaign.

`Race` could be `mortal`, `magic`, `beast`, and `mech`. It affects the tech tree this unit belongs to in campaign. Also affects spells like beast healing.

`Side` could be `good`, `evil`, and `neutral`. The player can normally recruit `good` and `neutral` (mechs, beasts, magic) units. The enemy can recruit `evil` and `neutral` units. Changing the unit's side from `evil` to `good` will make them available for recruiting in the campaign.

`BuyPrice` is an amount of gold required to recruit this unit.

`MeleeDamage` is a **melee** damage score of this unit.

`MeleeDamageType` could be `normal`, `magic`, and `mech`.

Skirmish, Missile, and Siege damage use the same rules as **melee** damage properties. A unit can have multiple damage scores as illustrated by Dreggo. If the damage score is **0**, it will not be listed and its damage type doesn't affect anything. So, Dreggo has no **siege** attack showed in the game.

`UnitClass` could be:
* `heavy_infantry`
* `light_infantry`
* `skirmisher`
* `archer`
* `cavalry`
* `light_cavalry`
* `sky_hunter`
* `bomber`
* `siege_engine`
* `spell_caster`

Note that some stuff is hardcoded into the unit class. AFAIK, the support fire is an `archer` thing and there is no `ability` that would cover that.

`Name` is the unit's name and could be anything, but the length is limited to 29 characters (well, it's 30, but could be risky as it would override the terminating null byte which could be important).

`Life` is the number of max **life** score the unit has. Single entities and heroes usually have 10, other units have 15. Any value is valid though.

`Attacks` is the number of max **attack** score the unit has. It drops down when the current **life** is reduced because of wounds and killed units. The drop is proportional to the max life, also affected by whether it is a single entity or not. The game usually uses 15 and 10 values, but anything above or below is valid too.

`Armor` is the unit's **armor** rating.

`MagicDefense` is the unit's **magic defense** rating. The game uses only a few options, like 75, 50, and some other, but any reasonable value can be used.

`Speed` is the number of tiles the unit can cross on a clear terrain. The cost of movement depends on the unit's movement type.

`MovementType` could be `light_infantry`, `heavy_infantry`, `cavalry`, `siege_engine`, `flying`. The open/closed formation of a heavy infantry is usually reflected here: Samurai has `light_infantry` movement type.

`SearchRange` is the unit's **search range** rating.

`Spell` is the the unit's cast ability. There could be only one. Possible options:
* An empty string for "none"
* `death_wounds`
* `summon_banshee`
* `firestorm`
* `animate_trees`
* `raise_dead`
* `heal`
* `earthquake`
* `eagle_eye`
* `invisibility`
* `force_march`
* `armor`
* `dispel_undead`
* `berserk`
* `charm`
* `restore_morale`
* `path_master`
* `panic`
* `whirlwind`
* `bless`
* `storm` - from the item
* `sickness` - from the item
* `drain_life` - an unused (?) spell that deals some 1-target damage
* `confusion` - from the item
* `fatigue` - from the item
* `short_fireball` - from the item
* `fireball` - from the item
* `hero_mass_heal` - a cavalry hero spell
* `hero_fireball` - an archmage hero spell
* `hero_fear` - an archmage hero spell
* `hero_fatigue` - an archmage hero spell
* `hero_whirlwind` - an archmage hero spell
* `hero_heal_beasts` - a beastmaster hero spell
* `hero_plague` - a beastmaster hero spell

Note that units can cast a hero spell. Their usage is not limited. This means a `hero_plague` spell can be cast every turn by a unit. Some of the spells above come from the item.

`Ability1`, `Ability2`, and `Ability3` represent a list of unit's passives. Here is a list of possible values:
* An empty string for "none"
* `negate_cavalry_bonus`
* `vulnerable_to_mech`
* `spell_resistant`
* `cant_be_healed`
* `single_entity`
* `regeneration`
* `heroic`
* `hero`
* `undead`
* `bloodlust`
* `vulnerable_to_magic`
* `mystic`
* `berserk`

If unit has less than 3 abilities, it's recommended to have blanks go after used ability slots.

You may want to refer to a game manual for more details about the specific stats and abilities.
