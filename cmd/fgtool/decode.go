package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quasilyte/fantasy-general-tools/dat"
)

func doDecode(args []string) error {
	var gamePath string
	var outputPath string
	var verbose bool

	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet("fgtool decode", flag.ExitOnError)
	fs.StringVar(&gamePath, "game-root", workdir,
		`Fantasy General game root`)
	fs.StringVar(&outputPath, "o", "_output",
		`fgtool decoding result output directory`)
	fs.BoolVar(&verbose, "v", false,
		`enable verbose output`)
	fs.Parse(args)

	logf := func(format string, args ...any) {}
	if verbose {
		logf = func(format string, args ...any) {
			fmt.Println(fmt.Sprintf(format, args...))
		}
	}

	if gamePath == "" {
		return fmt.Errorf("game-root can't be empty")
	}
	if outputPath == "" {
		return fmt.Errorf("output can't be empty")
	}

	if err := os.RemoveAll(outputPath); err != nil {
		return err
	}
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	magequip, err := dat.ParseMagequipFile(filepath.Join(gamePath, "DAT", "MAGEQUIP.DAT"))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(outputPath, "units"), os.ModePerm); err != nil {
		return err
	}
	for _, u := range magequip.Units {
		filename := filepath.Join(outputPath, "units", fmt.Sprintf("%03d_%s", u.Index, u.Name)+".json")
		if err := writeJSON(filename, u.ToSerdat()); err != nil {
			return err
		}
	}

	paletteData, err := os.ReadFile(filepath.Join(gamePath, "DAT", "FGPAL01.LBM"))
	if err != nil {
		return fmt.Errorf("read palette: %v", err)
	}
	pal, err := dat.ParsePaletteFile(paletteData)
	if err != nil {
		return fmt.Errorf("parse palette: %v", err)
	}

	if err := writeJSON(filepath.Join(outputPath, "palette.json"), pal.ToSerdat()); err != nil {
		return err
	}

	// We can't reconstruct the original LBM file from palette.json as of yet,
	// so copy an original just to make "encode" work.
	if err := os.MkdirAll(filepath.Join(outputPath, "_orig"), os.ModePerm); err != nil {
		return err
	}
	origPalFilename := filepath.Join(outputPath, "_orig", "FGPAL01.LBM")
	if err := os.WriteFile(origPalFilename, paletteData, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(outputPath, "images"), os.ModePerm); err != nil {
		return err
	}
	shpFileList, err := os.ReadDir(filepath.Join(gamePath, "SHP"))
	if err != nil {
		return err
	}
	skippedSHP := 0
	for _, f := range shpFileList {
		if !strings.HasSuffix(f.Name(), ".SHP") {
			continue
		}
		logf("unpacking %q", f.Name())
		shpData, err := os.ReadFile(filepath.Join(gamePath, "SHP", f.Name()))
		if err != nil {
			return err
		}
		shp, err := dat.ParseShpFile(shpData, pal)
		if err != nil {
			return fmt.Errorf("parse shp %q: %v", f.Name(), err)
		}
		if !shp.UniPalette {
			skippedSHP++
			logf("[!] skip %q: multiple palette-spritesheet", f.Name())
			continue
		}
		switch {
		case strings.HasPrefix(f.Name(), "ANIM"):
			// OK: unit death animation
		case strings.HasPrefix(f.Name(), "FANIM"):
			// OK: unit fight animation
		case strings.HasPrefix(f.Name(), "UIF"):
			// OK: unit roster picture
		case strings.HasPrefix(f.Name(), "UNITICON"):
			// OK: unit idle spritesheet (includes all units, varying size frame)
		case strings.HasPrefix(f.Name(), "T_"):
			// OK: tilesets
		case strings.HasPrefix(f.Name(), "ST_"):
			// OK: scaled-down tilesets
		case strings.HasPrefix(f.Name(), "SP"):
			// OK: spells
		case strings.HasPrefix(f.Name(), "VIC"):
			// OK: victory pictures
		default:
			skippedSHP++
			logf("[!] skip %q: this file is not supported yet", f.Name())
			continue
		}
		images := dat.ShpToPNG(shp)
		if len(images) == 1 {
			filename := filepath.Join(outputPath, "images", fmt.Sprintf("%s.%d.png", strings.TrimSuffix(f.Name(), ".SHP"), len(shp.Images)))
			if err := writePNG(filename, images[0]); err != nil {
				return fmt.Errorf("convert shp to png %q: %v", f.Name(), err)
			}
			continue
		}
		dirPath := filepath.Join(outputPath, "images", f.Name())
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
		for i, img := range images {
			filename := filepath.Join(dirPath, fmt.Sprintf("%03d.png", i))
			if err := writePNG(filename, img); err != nil {
				return fmt.Errorf("convert shp[%d] to png %q: %v", i, f.Name(), err)
			}
		}
	}
	logf("skipped %d shp files", skippedSHP)

	return nil
}
