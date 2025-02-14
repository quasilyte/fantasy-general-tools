# Data Files

* DAT - contains maps, campaign scenarios, palettes
* EXE - executable file
* MUSIC - OGG music files, no decoding/encoding required
* SAVES - stores SAV and ARM files, **neither are researched yet**
* SHP - image files (sprites, UI files, animations)
* SOUND - MEL-encoded sound files, **not researched yet**

An interesting fact: removing a game file, like `SCEN2.SCN` seem to make the game load that file from the mounted disk. If you have `SCEN2.SCN` present, it will be loaded instead. Even the GoG version of the game mounts the disk implicitly (`game.gog` looks like a disk image).

## DAT

The DAT folder contains these files:

```
      1 LBM
      1 SMK
      4 DAT
     62 MAP
    121 SCN
    372 GEO
```

`LBM` stores the game palette (256 colors).
The format is described here: https://wiki.amigaos.net/wiki/ILBM_IFF_Interleaved_Bitmap

`SMK` is most likely a file containing the in-game videos.
Could be the SMK2 type code.
https://en.wikipedia.org/wiki/Smacker_video

`DAT` files have different contents and encodings:
* `MAGEQUIP.DAT` - contains stats of all units
* `PREFS.DAT` - could be the settings file, not sure
* `TFONT2.DAT`, `TFONTBIG.DAT` - font-related files

> See [units.md](./units.md) to learn more about MAGEQUIP

`MAP` files are a shorter form of `SCN`. They contain the tilemap for arena battles. Some of the maps (most of them?) have identical landscape as campaign maps (SCN).

`SCN` contains a `MAP` data inside as well as extra campaign-related info. For example, it describes deployment area, enemy starting units and their stats/locations.

`GEO` files are still a mystery and TODO. I haven't figured them out yet.

## SHP

FG1 SHP-files have a format very close to PG2 SHP files:
https://www.luis-guzman.com/zips/PG2_FilesSpec.html

An SHP file can encode one or more images. They may have different sizes.

## SAVES

TODO: ARM and SAV files.

## SOUND

TODO: MEL files.
