package generator

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/0x0f0f0f/tripbot9000/util"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

type MandelgemCommand struct {
	fs        *flag.FlagSet
	filename  string
	img_w     vg.Length
	img_h     vg.Length
	img_w_str string
	img_h_str string
	dpi       int
	canvas    vgimg.PngCanvas
	seed      int64
}

func NewMandelgemCommand() *MandelgemCommand {
	c := &MandelgemCommand{
		fs: flag.NewFlagSet("mandelgem", flag.ContinueOnError),
	}

	now := time.Now()
	defaultfile := fmt.Sprintf("output-%s.png", now.Format("Jan-2-2006-15-04-05"))
	c.fs.StringVar(&c.filename, "o", defaultfile, "output image file name or directory")
	c.fs.StringVar(&c.img_w_str, "w", "80cm", "image width")
	c.fs.StringVar(&c.img_h_str, "h", "80cm", "image height")
	c.fs.IntVar(&c.dpi, "dpi", 120, "resolution in dpi")
	c.fs.Int64Var(&c.seed, "s", now.UnixNano(), "random seed")

	stat, err := os.Stat(c.filename)
	if err == nil {
		if stat.IsDir() {
			c.filename = filepath.Join(c.filename, defaultfile)
		}
	}

	return c
}

func (c *MandelgemCommand) Init(args []string) error {
	err := c.fs.Parse(args)

	c.img_w, err = vg.ParseLength(c.img_w_str)
	if err != nil {
		log.Fatalf("could not parse width: %+v", err)
	}
	c.img_h, err = vg.ParseLength(c.img_h_str)
	if err != nil {
		log.Fatalf("could not parse height: %+v", err)
	}

	c.canvas = vgimg.PngCanvas{
		Canvas: vgimg.NewWith(
			vgimg.UseWH(c.img_w, c.img_h),
			vgimg.UseBackgroundColor(color.Transparent),
			vgimg.UseDPI(c.dpi),
		),
	}
	fmt.Println("seed is", c.seed)
	rand.Seed(c.seed)

	return err
}

func (c *MandelgemCommand) Name() string {
	return c.fs.Name()
}

func (c *MandelgemCommand) Run() error {
	c.canvas.SetLineWidth(1 * vg.Millimeter)
	numfigs := util.RandInt(1, 10)

	for i := 0; i < numfigs; i++ {
		outer_radius := vg.Length(util.RandFloat(float64(c.img_w/10), float64(c.img_w/2)))
		radius_step := outer_radius / vg.Length(util.RandFloat(4, 10))
		center := vg.Point{
			X: vg.Length(util.RandFloat(0, float64(c.img_w))),
			Y: vg.Length(util.RandFloat(0, float64(c.img_h))),
		}
		c.canvas = DrawRandomGemstone(c.canvas, center, outer_radius, radius_step, palette.Plan9).(vgimg.PngCanvas)

	}

	imgGeom := c.canvas.Image()
	imgMandel := RandomMandelbrot(c.filename, c.img_w, c.img_h, c.dpi)

	draw.Draw(imgMandel, imgGeom.Bounds(), imgGeom, image.Point{0, 0}, draw.Over)

	util.SavePNG(imgMandel, c.filename)

	return nil
}
