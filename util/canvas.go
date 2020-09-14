package util

import (
	"image"
	"image/png"
	"log"
	"os"

	"gonum.org/v1/plot/vg"
)

func SaveVGCanvas(c vg.CanvasWriterTo, file string) {
	// Save image
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("could not create file %s: %+v", file, err)
	}
	defer f.Close()
	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatalf("could not encode image: %+v", err)
	}
	return
}
func SavePNG(i *image.RGBA, file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("could not create file %s: %+v", file, err)
	}
	defer f.Close()
	err = png.Encode(f, i)
	if err != nil {
		log.Fatalf("could not encode image: %+v", err)
	}
}
