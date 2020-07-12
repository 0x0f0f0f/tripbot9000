package main

import (
	"fmt"
	"gonum.org/v1/plot/vg"
	// "gonum.org/v1/plot/vg/draw"
	// "gonum.org/v1/plot/vg/vgsvg"
	// "gonum.org/v1/plot/vg/vgeps"
	"flag"
	t "github.com/0x0f0f0f/trippygen/trip"
	"gonum.org/v1/plot/vg/vgimg"
	"image/color"
	"image/color/palette"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
	// "strings"
)

func SaveImage(c vg.CanvasWriterTo, file string) {
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

func main() {
	now := time.Now()
	defaultfile := fmt.Sprintf("output-%s.png", now.Format("Jan-2-2006-15-04-05"))
	// TODO parse flags
	filename := flag.String("o", defaultfile, "output image file name or directory")
	seed := flag.Int64("s", now.UnixNano(), "random seed")
	img_w_str := flag.String("w", "80cm", "image width")
	img_h_str := flag.String("h", "80cm", "image height")

	flag.Parse()

	//////////////////////////////////////////////////////////////////

	stat, err := os.Stat(*filename)
	if err == nil {
		if stat.IsDir() {
			*filename = filepath.Join(*filename, defaultfile)
		}
	}

	fmt.Println("seed is", seed)

	img_w, err := vg.ParseLength(*img_w_str)
	if err != nil {
		log.Fatalf("could not parse width: %+v", err)
	}
	img_h, err := vg.ParseLength(*img_h_str)
	if err != nil {
		log.Fatalf("could not parse height: %+v", err)
	}

	c := vgimg.PngCanvas{Canvas: vgimg.NewWith(
		vgimg.UseWH(img_w, img_h), vgimg.UseBackgroundColor(color.Transparent))}
	// c := vgsvg.New(img_w, img_h)
	// d := draw.New(c)

	center := vg.Point{X: img_w / 2, Y: img_h / 2}

	outer_radius := img_w / 2
	radius_step := 5 * vg.Centimeter

	rand.Seed(*seed)
	// bg_rect := vg.Rectangle{Min: vg.Point{0, 0}, Max: vg.Point{img_w, img_h}}
	c.SetLineWidth(2 * vg.Millimeter)

	c = t.DrawRandomGeom(c, center, outer_radius, radius_step, palette.Plan9).(vgimg.PngCanvas)
	// c.Fill(bg_rect.Path())
	// for i := 0; i < 10; i++ {
	// }
	SaveImage(c, *filename)
}
