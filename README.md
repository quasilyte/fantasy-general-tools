# Fantasy General 1 Tools

## Overview

`fgtool` is a Fantasy General modding toolkit.

* Decode game files (SHP->PNG, DAT->JSON, ...)
* Encode game files (PNG->SHP, JSON->DAT, ...)
* Apply binary patches to change some in-game constants

> You would need a legally-obtained copy of the game to use this tool.

## Installation

It's recommended to download the [latest release binary](https://github.com/quasilyte/fantasy-general-tools/releases) for your system.

You can also [build it from sources](_docs/from_sources.md).

## Documentation

* [fgtool.md](_docs/fgtool.md) - how to use this tool
* [data_files.md](_docs/data_files.md) - describes the known parts of Fantasy General file structure
* [units.md](_docs/units.md) - describes unit stats and how to edit them

## Support

Want to support the author? Add [my game](https://store.steampowered.com/app/3024370/NebuLeet) to your Steam wishlist.

## Useful Resources

Information used to create this tool:

* https://wiki.amigaos.net/wiki/ILBM_IFF_Interleaved_Bitmap
* https://moddingwiki.shikadi.net/wiki/LBM_Format#CMAP:_Palette
* http://blog.ssokolow.com/archives/2018/12/02/resources-for-reverse-engineering-16-bit-applications/
* https://www.luis-guzman.com/zips/PG2_FilesSpec.html
* https://groups.google.com/g/comp.sys.ibm.pc.games.strategic/c/Yt_sxvy67CM
* https://forum.shrapnelgames.com/showthread.php?p=856948

Tools used:

* `ndisasm` to dump DOS binary assembly
* `hexl` mode in Emacs to inspect/edit binary files
* basic tools like `grep`, `awk`
* custom Go scripts to patch the binary

