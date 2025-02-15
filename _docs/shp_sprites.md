# SHP and Sprites

## SHP Names

* `ANIM` prefix - unit death animations
* `FANIM` prefix - unit fight animation
* `UIF` prefix - unit sketch art in the detailed view
* `UNITICON` - contains all "idle" sprites of units
* `T_` prefix - tiles
* `ST_` prefix - "small" tiles, scaled down for the map view
* `SP` prefix - spell effects
* `VIC` prefix - victory screen images (a hero posing)

## Tiles

Tiles:

* `T_${LOC}CL.SHP` (e.g. `T_DESCL.SHP`) - clear and rough tiles
* `T_${LOC}CT.SHP` (e.g. `T_DESCT.SHP`) - cities, shrines, towers
* `T_${LOC}MN.SHP` (e.g. `T_DESMN.SHP`) - mountains
* `T_${LOC}RD.SHP` (e.g. `T_DESRD.SHP`) - roads
* `T_${LOC}RV.SHP` (e.g. `T_DESRV.SHP`) - river

Location codes:

* `DES` - desert
* `JUN` - jungle
* `SNW` - snowy
* `VOL` - volcano

There is an exception to this: plains have no "code" and their files are named like `T_CLEAR.SHP`, `T_RIVER.SHP`.

## Palette

Fantasy General 1 uses a limited 256-color palette. It's described inside `DAT/FGPAL01.LBM`.

We're mostly interested in the CMAP section. The `decode` command would extract the relevant parts into `_output/palette.json`. You can edit the colors, but keep in mind that it should stay the 256-color table - no extra colors. You can change the existing colors though.

The colors in that palette only contain 3 components: R, G, and B, the alpha channel is always opaque (`0xff`). The string `01ffab` means `R=0x01`, `G=0xff`, and `B=0xab` (in hex notation). The index of the color is hex-based too, so `ff` is 255.

Keep in mind that generated PNG use the palette that was inside the `game-root` folder, so if you want to change the palette and re-draw the sprites, do an extra round first:

1. `decode` to get the current `palette.json`
2. modify the palette colors
3. `encode` to put an updated palette inside the game
4. `decode` again to extract the files with new colors applied

The SHP files do not describe the exact colors on their own but refer to the global palette in LBM file.
