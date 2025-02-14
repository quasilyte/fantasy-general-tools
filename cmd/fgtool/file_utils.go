package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"os"
)

func readJSON(filename string, dst any) error {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, &dst)
}

func readPNG(filename string) (image.Image, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(data))
}

func writePNG(filename string, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	return os.WriteFile(filename, buf.Bytes(), os.ModePerm)
}

func writeJSON(filename string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
