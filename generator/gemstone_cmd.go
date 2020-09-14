package generator

import (
	"flag"
	"fmt"
	"image/color"
	"image/color/palette"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/0x0f0f0f/tripbot9000/util"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

type GemstoneCommand struct {
	fs           *flag.FlagSet
	filename     string
	img_w        vg.Length
	img_h        vg.Length
	img_w_str    string
	img_h_str    string
	dpi          int
	canvas       vgimg.PngCanvas
	outer_radius vg.Length
	radius_step  vg.Length
	center       vg.Point
	seed         int64
}

func NewGemstoneCommand() *GemstoneCommand {
	c := &GemstoneCommand{
		fs: flag.NewFlagSet("gemstone", flag.ContinueOnError),
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

func (c *GemstoneCommand) Init(args []string) error {
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

	c.outer_radius = c.img_w / 2
	c.radius_step = c.img_w / 25
	c.center = vg.Point{X: c.img_w / 2, Y: c.img_h / 2}

	fmt.Println("seed is", c.seed)
	rand.Seed(c.seed)

	return err
}

func (c *GemstoneCommand) Name() string {
	return c.fs.Name()
}

func (c *GemstoneCommand) Run() error {

	c.canvas.SetLineWidth(1 * vg.Millimeter)

	canv := DrawRandomGemstone(c.canvas,
		c.center,
		c.outer_radius,
		c.radius_step,
		palette.Plan9).(vgimg.PngCanvas)
	// c.Fill(bg_rect.Path())
	// for i := 0; i < 10; i++ {
	// }

	util.SaveVGCanvas(canv, c.filename)

	return nil
}
