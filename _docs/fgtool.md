# fgtool

`fgtool` is a command-line tool.

It is intended to be used to `decode` the Fantasy General files into easy-to-edit format and then `encode` them back to patch the game with modified files.

An example of this is `SHP` file. You can't edit sprites in this format easily, but when it's decoded into PNG files, you can use your favorite graphics editor to change them or create new ones. Then you invoke the `encode` command to turn these PNG files into `SHP`. The modified `SHP` files can be placed into the game `SHP/` folder.

Therefore, there are only two commands that are needed for the basic usage: `decode` and `encode`.

## Quick Start

Open the terminal/command line. You may want to change the **working directory** (`cd` into a different folder).

1. Decode the game data:

```
fgtool.exe decode -v --game-root "C:/GOG Games/Fantasy General"
```

This will create `_output` directory in your current directory. Take a look inside, edit some files.

2. Encode the modified game files:

```
fgtool.exe encode -v
```

This will create `_patch` directory. Copy/paste files inside it into the game's folder, replacing the original files.

## decode command

The `decode` command has some extra arguments. Use `--help` to learn about them:

```
fgtool.exe decode --help
```

## encode command

The `encode` command has some extra arguments. Use `--help` to learn about them:

```
fgtool.exe encode --help
```
