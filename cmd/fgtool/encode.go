package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/quasilyte/fantasy-general-tools/dat"
	"github.com/quasilyte/fantasy-general-tools/serdat"
)

func doEncode(args []string) error {
	var patchPath string
	var sourcePath string
	var verbose bool

	fs := flag.NewFlagSet("fgtool decode", flag.ExitOnError)
	fs.StringVar(&sourcePath, "source", "_output",
		`fgtool-decoded result directory`)
	fs.StringVar(&patchPath, "o", "_patch",
		`generated patch output directory`)
	fs.BoolVar(&verbose, "v", false,
		`enable verbose output`)
	fs.Parse(args)

	logf := func(format string, args ...any) {}
	if verbose {
		logf = func(format string, args ...any) {
			fmt.Println(fmt.Sprintf(format, args...))
		}
	}

	if patchPath == "" {
		return fmt.Errorf("patch path can't be empty")
	}
	if sourcePath == "" {
		return fmt.Errorf("source can't be empty")
	}

	if err := os.RemoveAll(patchPath); err != nil {
		return err
	}
	if err := os.MkdirAll(patchPath, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(patchPath, "DAT"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(patchPath, "SHP"), os.ModePerm); err != nil {
		return err
	}

	paletteData, err := os.ReadFile(filepath.Join(sourcePath, "_orig", "FGPAL01.LBM"))
	if err != nil {
		return fmt.Errorf("read palette: %v", err)
	}
	pal, err := dat.ParsePaletteFile(paletteData)
	if err != nil {
		return fmt.Errorf("parse palette: %v", err)
	}
	var serializedPal serdat.Palette
	if err := readJSON(filepath.Join(sourcePath, "palette.json"), &serializedPal); err != nil {
		return fmt.Errorf("read palette.json: %v", err)
	}
	newPaletteData, err := dat.PaletteEncode(paletteData, serializedPal)
	if err != nil {
		return fmt.Errorf("encode result palette: %v", err)
	}
	if err := os.WriteFile(filepath.Join(patchPath, "DAT", "FGPAL01.LBM"), newPaletteData, os.ModePerm); err != nil {
		return err
	}

	{
		unitsDir := filepath.Join(sourcePath, "units")
		files, err := os.ReadDir(unitsDir)
		if err != nil {
			return err
		}
		var units []serdat.MagequipUnit
		for _, f := range files {
			var u serdat.MagequipUnit
			if err := readJSON(filepath.Join(unitsDir, f.Name()), &u); err != nil {
				return err
			}
			units = append(units, u)
		}
		magequip, err := dat.MagequipFromSerdat(units)
		if err != nil {
			return err
		}
		magequipOutFilename := filepath.Join(patchPath, "DAT", "MAGEQUIP.DAT")
		data, err := dat.MagequipEncode(magequip)
		if err != nil {
			return err
		}
		if err := os.WriteFile(magequipOutFilename, data, os.ModePerm); err != nil {
			return err
		}
	}

	{
		imagesDir := filepath.Join(sourcePath, "images")
		files, err := os.ReadDir(imagesDir)
		if err != nil {
			return err
		}
		for _, f := range files {
			logf("packing %q", f.Name())
			var shp *dat.ShpFile
			if !f.IsDir() {
				pngImage, err := readPNG(filepath.Join(imagesDir, f.Name()))
				if err != nil {
					return fmt.Errorf("decoding png %q: %v", f.Name(), err)
				}
				shp = dat.ShpFromPNG(f.Name(), pngImage)
			} else {
				framesDir := filepath.Join(imagesDir, f.Name())
				frameFiles, err := os.ReadDir(framesDir)
				if err != nil {
					return err
				}
				images := make([]image.Image, 0, len(frameFiles))
				for _, f2 := range frameFiles {
					if !strings.HasSuffix(f2.Name(), ".png") {
						continue
					}
					pngImage, err := readPNG(filepath.Join(framesDir, f2.Name()))
					if err != nil {
						return fmt.Errorf("%q: %v", f.Name(), err)
					}
					images = append(images, pngImage)
				}
				shp = dat.ShpFromPNGList(images)
			}

			shpBytes, err := dat.ShpEncode(shp, pal)
			if err != nil {
				return fmt.Errorf("encode shp %q: %v", f.Name(), err)
			}
			nameParts := strings.Split(f.Name(), ".")
			outFilename := filepath.Join(patchPath, "SHP", nameParts[0]+".SHP")
			if err := os.WriteFile(outFilename, shpBytes, os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}
